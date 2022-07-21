package chainreg

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/brsuite/brond/chaincfg/chainhash"
	"github.com/brsuite/brond/rpcclient"
	"github.com/brsuite/bronutil"
	"github.com/brsuite/bronwallet/chain"
	"github.com/brsuite/neutrino"
	"github.com/brsuite/broln/blockcache"
	"github.com/brsuite/broln/chainntnfs"
	"github.com/brsuite/broln/chainntnfs/bitcoindnotify"
	"github.com/brsuite/broln/chainntnfs/btcdnotify"
	"github.com/brsuite/broln/chainntnfs/neutrinonotify"
	"github.com/brsuite/broln/channeldb"
	"github.com/brsuite/broln/htlcswitch"
	"github.com/brsuite/broln/input"
	"github.com/brsuite/broln/keychain"
	"github.com/brsuite/broln/kvdb"
	"github.com/brsuite/broln/lncfg"
	"github.com/brsuite/broln/lnwallet"
	"github.com/brsuite/broln/lnwallet/chainfee"
	"github.com/brsuite/broln/lnwire"
	"github.com/brsuite/broln/routing/chainview"
	"github.com/brsuite/broln/walletunlocker"
)

// Config houses necessary fields that a chainControl instance needs to
// function.
type Config struct {
	// Brocoin defines settings for the Brocoin chain.
	Brocoin *lncfg.Chain

	// Litecoin defines settings for the Litecoin chain.
	Litecoin *lncfg.Chain

	// PrimaryChain is a function that returns our primary chain via its
	// ChainCode.
	PrimaryChain func() ChainCode

	// HeightHintCacheQueryDisable is a boolean that disables height hint
	// queries if true.
	HeightHintCacheQueryDisable bool

	// NeutrinoMode defines settings for connecting to a neutrino
	// light-client.
	NeutrinoMode *lncfg.Neutrino

	// BrocoindMode defines settings for connecting to a brocoind node.
	BrocoindMode *lncfg.Brocoind

	// LitecoindMode defines settings for connecting to a litecoind node.
	LitecoindMode *lncfg.Brocoind

	// BrondMode defines settings for connecting to a brond node.
	BrondMode *lncfg.Brond

	// LtcdMode defines settings for connecting to an ltcd node.
	LtcdMode *lncfg.Brond

	// HeightHintDB is a pointer to the database that stores the height
	// hints.
	HeightHintDB kvdb.Backend

	// ChanStateDB is a pointer to the database that stores the channel
	// state.
	ChanStateDB *channeldb.ChannelStateDB

	// BlockCache is the main cache for storing block information.
	BlockCache *blockcache.BlockCache

	// WalletUnlockParams are the parameters that were used for unlocking
	// the main wallet.
	WalletUnlockParams *walletunlocker.WalletUnlockParams

	// NeutrinoCS is a pointer to a neutrino ChainService. Must be non-nil
	// if using neutrino.
	NeutrinoCS *neutrino.ChainService

	// ActiveNetParams details the current chain we are on.
	ActiveNetParams BrocoinNetParams

	// FeeURL defines the URL for fee estimation we will use. This field is
	// optional.
	FeeURL string

	// Dialer is a function closure that will be used to establish outbound
	// TCP connections to Brocoin peers in the event of a pruned block being
	// requested.
	Dialer chain.Dialer
}

const (
	// DefaultBrocoinMinHTLCInMSat is the default smallest value htlc this
	// node will accept. This value is proposed in the channel open sequence
	// and cannot be changed during the life of the channel. It is 1 msat by
	// default to allow maximum flexibility in deciding what size payments
	// to forward.
	//
	// All forwarded payments are subjected to the min htlc constraint of
	// the routing policy of the outgoing channel. This implicitly controls
	// the minimum htlc value on the incoming channel too.
	DefaultBrocoinMinHTLCInMSat = lnwire.MilliSatoshi(1)

	// DefaultBrocoinMinHTLCOutMSat is the default minimum htlc value that
	// we require for sending out htlcs. Our channel peer may have a lower
	// min htlc channel parameter, but we - by default - don't forward
	// anything under the value defined here.
	DefaultBrocoinMinHTLCOutMSat = lnwire.MilliSatoshi(1000)

	// DefaultBrocoinBaseFeeMSat is the default forwarding base fee.
	DefaultBrocoinBaseFeeMSat = lnwire.MilliSatoshi(1000)

	// DefaultBrocoinFeeRate is the default forwarding fee rate.
	DefaultBrocoinFeeRate = lnwire.MilliSatoshi(1)

	// DefaultBrocoinTimeLockDelta is the default forwarding time lock
	// delta.
	DefaultBrocoinTimeLockDelta = 40

	DefaultLitecoinMinHTLCInMSat  = lnwire.MilliSatoshi(1)
	DefaultLitecoinMinHTLCOutMSat = lnwire.MilliSatoshi(1000)
	DefaultLitecoinBaseFeeMSat    = lnwire.MilliSatoshi(1000)
	DefaultLitecoinFeeRate        = lnwire.MilliSatoshi(1)
	DefaultLitecoinTimeLockDelta  = 576
	DefaultLitecoinDustLimit      = bronutil.Amount(54600)

	// DefaultBrocoinStaticFeePerKW is the fee rate of 50 sat/vbyte
	// expressed in sat/kw.
	DefaultBrocoinStaticFeePerKW = chainfee.SatPerKWeight(12500)

	// DefaultBrocoinStaticMinRelayFeeRate is the min relay fee used for
	// static estimators.
	DefaultBrocoinStaticMinRelayFeeRate = chainfee.FeePerKwFloor

	// DefaultLitecoinStaticFeePerKW is the fee rate of 200 sat/vbyte
	// expressed in sat/kw.
	DefaultLitecoinStaticFeePerKW = chainfee.SatPerKWeight(50000)

	// BtcToLtcConversionRate is a fixed ratio used in order to scale up
	// payments when running on the Litecoin chain.
	BtcToLtcConversionRate = 60
)

// DefaultLtcChannelConstraints is the default set of channel constraints that
// are meant to be used when initially funding a Litecoin channel.
var DefaultLtcChannelConstraints = channeldb.ChannelConstraints{
	DustLimit:        DefaultLitecoinDustLimit,
	MaxAcceptedHtlcs: input.MaxHTLCNumber / 2,
}

// PartialChainControl contains all the primary interfaces of the chain control
// that can be purely constructed from the global configuration. No wallet
// instance is required for constructing this partial state.
type PartialChainControl struct {
	// Cfg is the configuration that was used to create the partial chain
	// control.
	Cfg *Config

	// HealthCheck is a function which can be used to send a low-cost, fast
	// query to the chain backend to ensure we still have access to our
	// node.
	HealthCheck func() error

	// FeeEstimator is used to estimate an optimal fee for transactions
	// important to us.
	FeeEstimator chainfee.Estimator

	// ChainNotifier is used to receive blockchain events that we are
	// interested in.
	ChainNotifier chainntnfs.ChainNotifier

	// ChainView is used in the router for maintaining an up-to-date graph.
	ChainView chainview.FilteredChainView

	// ChainSource is the primary chain interface. This is used to operate
	// the wallet and do things such as rescanning, sending transactions,
	// notifications for received funds, etc.
	ChainSource chain.Interface

	// RoutingPolicy is the routing policy we have decided to use.
	RoutingPolicy htlcswitch.ForwardingPolicy

	// MinHtlcIn is the minimum HTLC we will accept.
	MinHtlcIn lnwire.MilliSatoshi

	// ChannelConstraints is the set of default constraints that will be
	// used for any incoming or outgoing channel reservation requests.
	ChannelConstraints channeldb.ChannelConstraints
}

// ChainControl couples the three primary interfaces broln utilizes for a
// particular chain together. A single ChainControl instance will exist for all
// the chains broln is currently active on.
type ChainControl struct {
	// PartialChainControl is the part of the chain control that was
	// initialized purely from the configuration and doesn't contain any
	// wallet related elements.
	*PartialChainControl

	// ChainIO represents an abstraction over a source that can query the
	// blockchain.
	ChainIO lnwallet.BlockChainIO

	// Signer is used to provide signatures over things like transactions.
	Signer input.Signer

	// KeyRing represents a set of keys that we have the private keys to.
	KeyRing keychain.SecretKeyRing

	// Wc is an abstraction over some basic wallet commands. This base set
	// of commands will be provided to the Wallet *LightningWallet raw
	// pointer below.
	Wc lnwallet.WalletController

	// MsgSigner is used to sign arbitrary messages.
	MsgSigner lnwallet.MessageSigner

	// Wallet is our LightningWallet that also contains the abstract Wc
	// above. This wallet handles all of the lightning operations.
	Wallet *lnwallet.LightningWallet
}

// GenDefaultBtcConstraints generates the default set of channel constraints
// that are to be used when funding a Brocoin channel.
func GenDefaultBtcConstraints() channeldb.ChannelConstraints {
	// We use the dust limit for the maximally sized witness program with
	// a 40-byte data push.
	dustLimit := lnwallet.DustLimitForSize(input.UnknownWitnessSize)

	return channeldb.ChannelConstraints{
		DustLimit:        dustLimit,
		MaxAcceptedHtlcs: input.MaxHTLCNumber / 2,
	}
}

// NewPartialChainControl creates a new partial chain control that contains all
// the parts that can be purely constructed from the passed in global
// configuration and doesn't need any wallet instance yet.
func NewPartialChainControl(cfg *Config) (*PartialChainControl, func(), error) {
	// Set the RPC config from the "home" chain. Multi-chain isn't yet
	// active, so we'll restrict usage to a particular chain for now.
	homeChainConfig := cfg.Brocoin
	if cfg.PrimaryChain() == LitecoinChain {
		homeChainConfig = cfg.Litecoin
	}
	log.Infof("Primary chain is set to: %v", cfg.PrimaryChain())

	cc := &PartialChainControl{
		Cfg: cfg,
	}

	switch cfg.PrimaryChain() {
	case BrocoinChain:
		cc.RoutingPolicy = htlcswitch.ForwardingPolicy{
			MinHTLCOut:    cfg.Brocoin.MinHTLCOut,
			BaseFee:       cfg.Brocoin.BaseFee,
			FeeRate:       cfg.Brocoin.FeeRate,
			TimeLockDelta: cfg.Brocoin.TimeLockDelta,
		}
		cc.MinHtlcIn = cfg.Brocoin.MinHTLCIn
		cc.FeeEstimator = chainfee.NewStaticEstimator(
			DefaultBrocoinStaticFeePerKW,
			DefaultBrocoinStaticMinRelayFeeRate,
		)
	case LitecoinChain:
		cc.RoutingPolicy = htlcswitch.ForwardingPolicy{
			MinHTLCOut:    cfg.Litecoin.MinHTLCOut,
			BaseFee:       cfg.Litecoin.BaseFee,
			FeeRate:       cfg.Litecoin.FeeRate,
			TimeLockDelta: cfg.Litecoin.TimeLockDelta,
		}
		cc.MinHtlcIn = cfg.Litecoin.MinHTLCIn
		cc.FeeEstimator = chainfee.NewStaticEstimator(
			DefaultLitecoinStaticFeePerKW, 0,
		)
	default:
		return nil, nil, fmt.Errorf("default routing policy for chain "+
			"%v is unknown", cfg.PrimaryChain())
	}

	var err error
	heightHintCacheConfig := chainntnfs.CacheConfig{
		QueryDisable: cfg.HeightHintCacheQueryDisable,
	}
	if cfg.HeightHintCacheQueryDisable {
		log.Infof("Height Hint Cache Queries disabled")
	}

	// Initialize the height hint cache within the chain directory.
	hintCache, err := chainntnfs.NewHeightHintCache(
		heightHintCacheConfig, cfg.HeightHintDB,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to initialize height hint "+
			"cache: %v", err)
	}

	// If spv mode is active, then we'll be using a distinct set of
	// chainControl interfaces that interface directly with the p2p network
	// of the selected chain.
	switch homeChainConfig.Node {
	case "neutrino":
		// We'll create ChainNotifier and FilteredChainView instances,
		// along with the wallet's ChainSource, which are all backed by
		// the neutrino light client.
		cc.ChainNotifier = neutrinonotify.New(
			cfg.NeutrinoCS, hintCache, hintCache, cfg.BlockCache,
		)
		cc.ChainView, err = chainview.NewCfFilteredChainView(
			cfg.NeutrinoCS, cfg.BlockCache,
		)
		if err != nil {
			return nil, nil, err
		}

		// Map the deprecated neutrino feeurl flag to the general fee
		// url.
		if cfg.NeutrinoMode.FeeURL != "" {
			if cfg.FeeURL != "" {
				return nil, nil, errors.New("feeurl and " +
					"neutrino.feeurl are mutually " +
					"exclusive")
			}

			cfg.FeeURL = cfg.NeutrinoMode.FeeURL
		}

		cc.ChainSource = chain.NewNeutrinoClient(
			cfg.ActiveNetParams.Params, cfg.NeutrinoCS,
		)

		// Get our best block as a health check.
		cc.HealthCheck = func() error {
			_, _, err := cc.ChainSource.GetBestBlock()
			return err
		}

	case "brocoind", "litecoind":
		var brocoindMode *lncfg.Brocoind
		switch {
		case cfg.Brocoin.Active:
			brocoindMode = cfg.BrocoindMode
		case cfg.Litecoin.Active:
			brocoindMode = cfg.LitecoindMode
		}
		// Otherwise, we'll be speaking directly via RPC and ZMQ to a
		// brocoind node. If the specified host for the brond/ltcd RPC
		// server already has a port specified, then we use that
		// directly. Otherwise, we assume the default port according to
		// the selected chain parameters.
		var brocoindHost string
		if strings.Contains(brocoindMode.RPCHost, ":") {
			brocoindHost = brocoindMode.RPCHost
		} else {
			// The RPC ports specified in chainparams.go assume
			// brond, which picks a different port so that btcwallet
			// can use the same RPC port as brocoind. We convert
			// this back to the btcwallet/brocoind port.
			rpcPort, err := strconv.Atoi(cfg.ActiveNetParams.RPCPort)
			if err != nil {
				return nil, nil, err
			}
			rpcPort -= 2
			brocoindHost = fmt.Sprintf("%v:%d",
				brocoindMode.RPCHost, rpcPort)
			if (cfg.Brocoin.Active &&
				(cfg.Brocoin.RegTest || cfg.Brocoin.SigNet)) ||
				(cfg.Litecoin.Active && cfg.Litecoin.RegTest) {

				conn, err := net.Dial("tcp", brocoindHost)
				if err != nil || conn == nil {
					switch {
					case cfg.Brocoin.Active && cfg.Brocoin.RegTest:
						rpcPort = 18871
					case cfg.Litecoin.Active && cfg.Litecoin.RegTest:
						rpcPort = 19443
					case cfg.Brocoin.Active && cfg.Brocoin.SigNet:
						rpcPort = 38332
					}
					brocoindHost = fmt.Sprintf("%v:%d",
						brocoindMode.RPCHost,
						rpcPort)
				} else {
					conn.Close()
				}
			}
		}

		// Establish the connection to brocoind and create the clients
		// required for our relevant subsystems.
		brocoindConn, err := chain.NewBrocoindConn(&chain.BrocoindConfig{
			ChainParams:        cfg.ActiveNetParams.Params,
			Host:               brocoindHost,
			User:               brocoindMode.RPCUser,
			Pass:               brocoindMode.RPCPass,
			ZMQBlockHost:       brocoindMode.ZMQPubRawBlock,
			ZMQTxHost:          brocoindMode.ZMQPubRawTx,
			ZMQReadDeadline:    5 * time.Second,
			Dialer:             cfg.Dialer,
			PrunedModeMaxPeers: brocoindMode.PrunedNodeMaxPeers,
		})
		if err != nil {
			return nil, nil, err
		}

		if err := brocoindConn.Start(); err != nil {
			return nil, nil, fmt.Errorf("unable to connect to "+
				"brocoind: %v", err)
		}

		cc.ChainNotifier = bitcoindnotify.New(
			brocoindConn, cfg.ActiveNetParams.Params, hintCache,
			hintCache, cfg.BlockCache,
		)
		cc.ChainView = chainview.NewBrocoindFilteredChainView(
			brocoindConn, cfg.BlockCache,
		)
		cc.ChainSource = brocoindConn.NewBrocoindClient()

		// If we're not in regtest mode, then we'll attempt to use a
		// proper fee estimator for testnet.
		rpcConfig := &rpcclient.ConnConfig{
			Host:                 brocoindHost,
			User:                 brocoindMode.RPCUser,
			Pass:                 brocoindMode.RPCPass,
			DisableConnectOnNew:  true,
			DisableAutoReconnect: false,
			DisableTLS:           true,
			HTTPPostMode:         true,
		}
		if cfg.Brocoin.Active && !cfg.Brocoin.RegTest {
			log.Infof("Initializing brocoind backed fee estimator "+
				"in %s mode", brocoindMode.EstimateMode)

			// Finally, we'll re-initialize the fee estimator, as
			// if we're using brocoind as a backend, then we can
			// use live fee estimates, rather than a statically
			// coded value.
			fallBackFeeRate := chainfee.SatPerKVByte(25 * 1000)
			cc.FeeEstimator, err = chainfee.NewBrocoindEstimator(
				*rpcConfig, brocoindMode.EstimateMode,
				fallBackFeeRate.FeePerKWeight(),
			)
			if err != nil {
				return nil, nil, err
			}
		} else if cfg.Litecoin.Active && !cfg.Litecoin.RegTest {
			log.Infof("Initializing litecoind backed fee "+
				"estimator in %s mode",
				brocoindMode.EstimateMode)

			// Finally, we'll re-initialize the fee estimator, as
			// if we're using litecoind as a backend, then we can
			// use live fee estimates, rather than a statically
			// coded value.
			fallBackFeeRate := chainfee.SatPerKVByte(25 * 1000)
			cc.FeeEstimator, err = chainfee.NewBrocoindEstimator(
				*rpcConfig, brocoindMode.EstimateMode,
				fallBackFeeRate.FeePerKWeight(),
			)
			if err != nil {
				return nil, nil, err
			}
		}

		// We need to use some apis that are not exposed by btcwallet,
		// for a health check function so we create an ad-hoc brocoind
		// connection.
		chainConn, err := rpcclient.New(rpcConfig, nil)
		if err != nil {
			return nil, nil, err
		}

		// The api we will use for our health check depends on the
		// brocoind version.
		cmd, err := getBrocoindHealthCheckCmd(chainConn)
		if err != nil {
			return nil, nil, err
		}

		cc.HealthCheck = func() error {
			_, err := chainConn.RawRequest(cmd, nil)
			return err
		}

	case "brond", "ltcd":
		// Otherwise, we'll be speaking directly via RPC to a node.
		//
		// So first we'll load brond/ltcd's TLS cert for the RPC
		// connection. If a raw cert was specified in the config, then
		// we'll set that directly. Otherwise, we attempt to read the
		// cert from the path specified in the config.
		var brondMode *lncfg.Brond
		switch {
		case cfg.Brocoin.Active:
			brondMode = cfg.BrondMode
		case cfg.Litecoin.Active:
			brondMode = cfg.LtcdMode
		}
		var rpcCert []byte
		if brondMode.RawRPCCert != "" {
			rpcCert, err = hex.DecodeString(brondMode.RawRPCCert)
			if err != nil {
				return nil, nil, err
			}
		} else {
			certFile, err := os.Open(brondMode.RPCCert)
			if err != nil {
				return nil, nil, err
			}
			rpcCert, err = ioutil.ReadAll(certFile)
			if err != nil {
				return nil, nil, err
			}
			if err := certFile.Close(); err != nil {
				return nil, nil, err
			}
		}

		// If the specified host for the brond/ltcd RPC server already
		// has a port specified, then we use that directly. Otherwise,
		// we assume the default port according to the selected chain
		// parameters.
		var brondHost string
		if strings.Contains(brondMode.RPCHost, ":") {
			brondHost = brondMode.RPCHost
		} else {
			brondHost = fmt.Sprintf("%v:%v", brondMode.RPCHost,
				cfg.ActiveNetParams.RPCPort)
		}

		brondUser := brondMode.RPCUser
		brondPass := brondMode.RPCPass
		rpcConfig := &rpcclient.ConnConfig{
			Host:                 brondHost,
			Endpoint:             "ws",
			User:                 brondUser,
			Pass:                 brondPass,
			Certificates:         rpcCert,
			DisableTLS:           false,
			DisableConnectOnNew:  true,
			DisableAutoReconnect: false,
		}
		cc.ChainNotifier, err = btcdnotify.New(
			rpcConfig, cfg.ActiveNetParams.Params, hintCache,
			hintCache, cfg.BlockCache,
		)
		if err != nil {
			return nil, nil, err
		}

		// Finally, we'll create an instance of the default chain view
		// to be used within the routing layer.
		cc.ChainView, err = chainview.NewBrondFilteredChainView(
			*rpcConfig, cfg.BlockCache,
		)
		if err != nil {
			log.Errorf("unable to create chain view: %v", err)
			return nil, nil, err
		}

		// Create a special websockets rpc client for brond which will be
		// used by the wallet for notifications, calls, etc.
		chainRPC, err := chain.NewRPCClient(
			cfg.ActiveNetParams.Params, brondHost, brondUser,
			brondPass, rpcCert, false, 20,
		)
		if err != nil {
			return nil, nil, err
		}

		cc.ChainSource = chainRPC

		// Use a query for our best block as a health check.
		cc.HealthCheck = func() error {
			_, _, err := cc.ChainSource.GetBestBlock()
			return err
		}

		// If we're not in simnet or regtest mode, then we'll attempt
		// to use a proper fee estimator for testnet.
		if !cfg.Brocoin.SimNet && !cfg.Litecoin.SimNet &&
			!cfg.Brocoin.RegTest && !cfg.Litecoin.RegTest {

			log.Info("Initializing brond backed fee estimator")

			// Finally, we'll re-initialize the fee estimator, as
			// if we're using brond as a backend, then we can use
			// live fee estimates, rather than a statically coded
			// value.
			fallBackFeeRate := chainfee.SatPerKVByte(25 * 1000)
			cc.FeeEstimator, err = chainfee.NewBrondEstimator(
				*rpcConfig, fallBackFeeRate.FeePerKWeight(),
			)
			if err != nil {
				return nil, nil, err
			}
		}

	case "nochainbackend":
		backend := &NoChainBackend{}
		source := &NoChainSource{
			BestBlockTime: time.Now(),
		}

		cc.ChainNotifier = backend
		cc.ChainView = backend
		cc.FeeEstimator = backend

		cc.ChainSource = source
		cc.HealthCheck = func() error {
			return nil
		}

	default:
		return nil, nil, fmt.Errorf("unknown node type: %s",
			homeChainConfig.Node)
	}

	switch {
	// If the fee URL isn't set, and the user is running mainnet, then
	// we'll return an error to instruct them to set a proper fee
	// estimator.
	case cfg.FeeURL == "" && cfg.Brocoin.MainNet &&
		homeChainConfig.Node == "neutrino":

		return nil, nil, fmt.Errorf("--feeurl parameter required " +
			"when running neutrino on mainnet")

	// Override default fee estimator if an external service is specified.
	case cfg.FeeURL != "":
		// Do not cache fees on regtest to make it easier to execute
		// manual or automated test cases.
		cacheFees := !cfg.Brocoin.RegTest

		log.Infof("Using external fee estimator %v: cached=%v",
			cfg.FeeURL, cacheFees)

		cc.FeeEstimator = chainfee.NewWebAPIEstimator(
			chainfee.SparseConfFeeSource{
				URL: cfg.FeeURL,
			},
			!cacheFees,
		)
	}

	ccCleanup := func() {
		if cc.FeeEstimator != nil {
			if err := cc.FeeEstimator.Stop(); err != nil {
				log.Errorf("Failed to stop feeEstimator: %v",
					err)
			}
		}
	}

	// Start fee estimator.
	if err := cc.FeeEstimator.Start(); err != nil {
		return nil, nil, err
	}

	// Select the default channel constraints for the primary chain.
	cc.ChannelConstraints = GenDefaultBtcConstraints()
	if cfg.PrimaryChain() == LitecoinChain {
		cc.ChannelConstraints = DefaultLtcChannelConstraints
	}

	return cc, ccCleanup, nil
}

// NewChainControl attempts to create a ChainControl instance according
// to the parameters in the passed configuration. Currently three
// branches of ChainControl instances exist: one backed by a running brond
// full-node, another backed by a running brocoind full-node, and the other
// backed by a running neutrino light client instance. When running with a
// neutrino light client instance, `neutrinoCS` must be non-nil.
func NewChainControl(walletConfig lnwallet.Config,
	msgSigner lnwallet.MessageSigner,
	pcc *PartialChainControl) (*ChainControl, func(), error) {

	cc := &ChainControl{
		PartialChainControl: pcc,
		MsgSigner:           msgSigner,
		Signer:              walletConfig.Signer,
		ChainIO:             walletConfig.ChainIO,
		Wc:                  walletConfig.WalletController,
		KeyRing:             walletConfig.SecretKeyRing,
	}

	ccCleanup := func() {
		if cc.Wallet != nil {
			if err := cc.Wallet.Shutdown(); err != nil {
				log.Errorf("Failed to shutdown wallet: %v", err)
			}
		}
	}

	lnWallet, err := lnwallet.NewLightningWallet(walletConfig)
	if err != nil {
		return nil, ccCleanup, fmt.Errorf("unable to create wallet: %v",
			err)
	}
	if err := lnWallet.Startup(); err != nil {
		return nil, ccCleanup, fmt.Errorf("unable to create wallet: %v",
			err)
	}

	log.Info("LightningWallet opened")
	cc.Wallet = lnWallet

	return cc, ccCleanup, nil
}

// getBrocoindHealthCheckCmd queries brocoind for its version to decide which
// api we should use for our health check. We prefer to use the uptime
// command, because it has no locking and is an inexpensive call, which was
// added in version 0.15. If we are on an earlier version, we fallback to using
// getblockchaininfo.
func getBrocoindHealthCheckCmd(client *rpcclient.Client) (string, error) {
	// Query brocoind to get our current version.
	resp, err := client.RawRequest("getnetworkinfo", nil)
	if err != nil {
		return "", err
	}

	// Parse the response to retrieve brocoind's version.
	info := struct {
		Version int64 `json:"version"`
	}{}
	if err := json.Unmarshal(resp, &info); err != nil {
		return "", err
	}

	// Brocoind returns a single value representing the semantic version:
	// 1000000 * CLIENT_VERSION_MAJOR + 10000 * CLIENT_VERSION_MINOR
	// + 100 * CLIENT_VERSION_REVISION + 1 * CLIENT_VERSION_BUILD
	//
	// The uptime call was added in version 0.15.0, so we return it for
	// any version value >= 150000, as per the above calculation.
	if info.Version >= 150000 {
		return "uptime", nil
	}

	return "getblockchaininfo", nil
}

var (
	// BrocoinTestnetGenesis is the genesis hash of Brocoin's testnet
	// chain.
	BrocoinTestnetGenesis = chainhash.Hash([chainhash.HashSize]byte{
		0xce, 0xeb, 0x1a, 0x38, 0x06, 0xf0, 0x71, 0x0c, 
0xda, 0x62, 0xbd, 0x49, 0x32, 0x8d, 0x70, 0x62,
0xb1, 0x0c, 0xca, 0x75, 0xd1, 0xd8, 0x16, 0x9f,
0xa9, 0x9e, 0xed, 0x16, 0xfa, 0x0c, 0x00, 0x00, 
	})

	// BrocoinSignetGenesis is the genesis hash of Brocoin's signet chain.
	BrocoinSignetGenesis = chainhash.Hash([chainhash.HashSize]byte{
		0xf6, 0x1e, 0xee, 0x3b, 0x63, 0xa3, 0x80, 0xa4,
		0x77, 0xa0, 0x63, 0xaf, 0x32, 0xb2, 0xbb, 0xc9,
		0x7c, 0x9f, 0xf9, 0xf0, 0x1f, 0x2c, 0x42, 0x25,
		0xe9, 0x73, 0x98, 0x81, 0x08, 0x00, 0x00, 0x00,
	})

	// BrocoinMainnetGenesis is the genesis hash of Brocoin's main chain.
	BrocoinMainnetGenesis = chainhash.Hash([chainhash.HashSize]byte{
		0xd2, 0x28, 0x0d, 0x8c, 0xf4, 0x3a, 0x3e, 0xfd,
0x9a, 0x51, 0x41, 0x4b, 0x15, 0xbf, 0x6a, 0xb0,
0x2b, 0x1c, 0x12, 0xfb, 0x78, 0xd6, 0xb6, 0x9e,
0x63, 0xf8, 0x88, 0xc5, 0xe3, 0x18, 0xbf, 0x05, 
	})

	// LitecoinTestnetGenesis is the genesis hash of Litecoin's testnet4
	// chain.
	LitecoinTestnetGenesis = chainhash.Hash([chainhash.HashSize]byte{
		0xa0, 0x29, 0x3e, 0x4e, 0xeb, 0x3d, 0xa6, 0xe6,
		0xf5, 0x6f, 0x81, 0xed, 0x59, 0x5f, 0x57, 0x88,
		0x0d, 0x1a, 0x21, 0x56, 0x9e, 0x13, 0xee, 0xfd,
		0xd9, 0x51, 0x28, 0x4b, 0x5a, 0x62, 0x66, 0x49,
	})

	// LitecoinMainnetGenesis is the genesis hash of Litecoin's main chain.
	LitecoinMainnetGenesis = chainhash.Hash([chainhash.HashSize]byte{
		0xe2, 0xbf, 0x04, 0x7e, 0x7e, 0x5a, 0x19, 0x1a,
		0xa4, 0xef, 0x34, 0xd3, 0x14, 0x97, 0x9d, 0xc9,
		0x98, 0x6e, 0x0f, 0x19, 0x25, 0x1e, 0xda, 0xba,
		0x59, 0x40, 0xfd, 0x1f, 0xe3, 0x65, 0xa7, 0x12,
	})

	// chainMap is a simple index that maps a chain's genesis hash to the
	// ChainCode enum for that chain.
	chainMap = map[chainhash.Hash]ChainCode{
		BrocoinTestnetGenesis:  BrocoinChain,
		LitecoinTestnetGenesis: LitecoinChain,

		BrocoinMainnetGenesis:  BrocoinChain,
		LitecoinMainnetGenesis: LitecoinChain,
	}

	// ChainDNSSeeds is a map of a chain's hash to the set of DNS seeds
	// that will be use to bootstrap peers upon first startup.
	//
	// The first item in the array is the primary host we'll use to attempt
	// the SRV lookup we require. If we're unable to receive a response
	// over UDP, then we'll fall back to manual TCP resolution. The second
	// item in the array is a special A record that we'll query in order to
	// receive the IP address of the current authoritative DNS server for
	// the network seed.
	//
	// TODO(roasbeef): extend and collapse these and chainparams.go into
	// struct like chaincfg.Params
	ChainDNSSeeds = map[chainhash.Hash][][2]string{
		BrocoinMainnetGenesis: {
			{
				"207.180.196.129",
			},
			{
				//"lseed.brocoinstats.com",
			},
		},

		BrocoinTestnetGenesis: {
			{
				"207.180.196.129",
			},
		},

		BrocoinSignetGenesis: {
			{
				"ln.signet.secp.tech",
			},
		},

		LitecoinMainnetGenesis: {
			{
				"ltc.nodes.lightning.directory",
				"soa.nodes.lightning.directory",
			},
		},
	}
)

// ChainRegistry keeps track of the current chains
type ChainRegistry struct {
	sync.RWMutex

	activeChains map[ChainCode]*ChainControl
	netParams    map[ChainCode]*BrocoinNetParams

	primaryChain ChainCode
}

// NewChainRegistry creates a new ChainRegistry.
func NewChainRegistry() *ChainRegistry {
	return &ChainRegistry{
		activeChains: make(map[ChainCode]*ChainControl),
		netParams:    make(map[ChainCode]*BrocoinNetParams),
	}
}

// RegisterChain assigns an active ChainControl instance to a target chain
// identified by its ChainCode.
func (c *ChainRegistry) RegisterChain(newChain ChainCode,
	cc *ChainControl) {

	c.Lock()
	c.activeChains[newChain] = cc
	c.Unlock()
}

// LookupChain attempts to lookup an active ChainControl instance for the
// target chain.
func (c *ChainRegistry) LookupChain(targetChain ChainCode) (
	*ChainControl, bool) {

	c.RLock()
	cc, ok := c.activeChains[targetChain]
	c.RUnlock()
	return cc, ok
}

// LookupChainByHash attempts to look up an active ChainControl which
// corresponds to the passed genesis hash.
func (c *ChainRegistry) LookupChainByHash(
	chainHash chainhash.Hash) (*ChainControl, bool) {

	c.RLock()
	defer c.RUnlock()

	targetChain, ok := chainMap[chainHash]
	if !ok {
		return nil, ok
	}

	cc, ok := c.activeChains[targetChain]
	return cc, ok
}

// RegisterPrimaryChain sets a target chain as the "home chain" for broln.
func (c *ChainRegistry) RegisterPrimaryChain(cc ChainCode) {
	c.Lock()
	defer c.Unlock()

	c.primaryChain = cc
}

// PrimaryChain returns the primary chain for this running broln instance. The
// primary chain is considered the "home base" while the other registered
// chains are treated as secondary chains.
func (c *ChainRegistry) PrimaryChain() ChainCode {
	c.RLock()
	defer c.RUnlock()

	return c.primaryChain
}

// ActiveChains returns a slice containing the active chains.
func (c *ChainRegistry) ActiveChains() []ChainCode {
	c.RLock()
	defer c.RUnlock()

	chains := make([]ChainCode, 0, len(c.activeChains))
	for activeChain := range c.activeChains {
		chains = append(chains, activeChain)
	}

	return chains
}

// NumActiveChains returns the total number of active chains.
func (c *ChainRegistry) NumActiveChains() uint32 {
	c.RLock()
	defer c.RUnlock()

	return uint32(len(c.activeChains))
}
