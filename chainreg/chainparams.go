package chainreg

import (
	"github.com/brsuite/brond/chaincfg"
	brocoinCfg "github.com/brsuite/brond/chaincfg"
	"github.com/brsuite/brond/chaincfg/chainhash"
	brocoinWire "github.com/brsuite/brond/wire"
	"github.com/brsuite/broln/keychain"
	litecoinCfg "github.com/ltcsuite/ltcd/chaincfg"
	litecoinWire "github.com/ltcsuite/ltcd/wire"
)

// BrocoinNetParams couples the p2p parameters of a network with the
// corresponding RPC port of a daemon running on the particular network.
type BrocoinNetParams struct {
	*brocoinCfg.Params
	RPCPort  string
	CoinType uint32
}

// LitecoinNetParams couples the p2p parameters of a network with the
// corresponding RPC port of a daemon running on the particular network.
type LitecoinNetParams struct {
	*litecoinCfg.Params
	RPCPort  string
	CoinType uint32
}

// BrocoinTestNetParams contains parameters specific to the 3rd version of the
// test network.
var BrocoinTestNetParams = BrocoinNetParams{
	Params:   &brocoinCfg.TestNet3Params,
	RPCPort:  "18740",
	CoinType: keychain.CoinTypeTestnet,
}

// BrocoinMainNetParams contains parameters specific to the current Brocoin
// mainnet.
var BrocoinMainNetParams = BrocoinNetParams{
	Params:   &brocoinCfg.MainNetParams,
	RPCPort:  "8360",
	CoinType: keychain.CoinTypeBrocoin,
}

// BrocoinSimNetParams contains parameters specific to the simulation test
// network.
var BrocoinSimNetParams = BrocoinNetParams{
	Params:   &brocoinCfg.SimNetParams,
	RPCPort:  "18556",
	CoinType: keychain.CoinTypeTestnet,
}

// BrocoinSigNetParams contains parameters specific to the signet test network.
var BrocoinSigNetParams = BrocoinNetParams{
	Params:   &brocoinCfg.SigNetParams,
	RPCPort:  "38332",
	CoinType: keychain.CoinTypeTestnet,
}

// LitecoinSimNetParams contains parameters specific to the simulation test
// network.
var LitecoinSimNetParams = LitecoinNetParams{
	Params:   &litecoinCfg.TestNet4Params,
	RPCPort:  "18556",
	CoinType: keychain.CoinTypeTestnet,
}

// LitecoinTestNetParams contains parameters specific to the 4th version of the
// test network.
var LitecoinTestNetParams = LitecoinNetParams{
	Params:   &litecoinCfg.TestNet4Params,
	RPCPort:  "19334",
	CoinType: keychain.CoinTypeTestnet,
}

// LitecoinMainNetParams contains the parameters specific to the current
// Litecoin mainnet.
var LitecoinMainNetParams = LitecoinNetParams{
	Params:   &litecoinCfg.MainNetParams,
	RPCPort:  "9334",
	CoinType: keychain.CoinTypeLitecoin,
}

// LitecoinRegTestNetParams contains parameters specific to a local litecoin
// regtest network.
var LitecoinRegTestNetParams = LitecoinNetParams{
	Params:   &litecoinCfg.RegressionNetParams,
	RPCPort:  "18334",
	CoinType: keychain.CoinTypeTestnet,
}

// BrocoinRegTestNetParams contains parameters specific to a local brocoin
// regtest network.
var BrocoinRegTestNetParams = BrocoinNetParams{
	Params:   &brocoinCfg.RegressionNetParams,
	RPCPort:  "18871",
	CoinType: keychain.CoinTypeTestnet,
}

// ApplyLitecoinParams applies the relevant chain configuration parameters that
// differ for litecoin to the chain parameters typed for btcsuite derivation.
// This function is used in place of using something like interface{} to
// abstract over _which_ chain (or fork) the parameters are for.
func ApplyLitecoinParams(params *BrocoinNetParams,
	litecoinParams *LitecoinNetParams) {

	params.Name = litecoinParams.Name
	params.Net = brocoinWire.BrocoinNet(litecoinParams.Net)
	params.DefaultPort = litecoinParams.DefaultPort
	params.CoinbaseMaturity = litecoinParams.CoinbaseMaturity

	copy(params.GenesisHash[:], litecoinParams.GenesisHash[:])

	// Address encoding magics
	params.PubKeyHashAddrID = litecoinParams.PubKeyHashAddrID
	params.ScriptHashAddrID = litecoinParams.ScriptHashAddrID
	params.PrivateKeyID = litecoinParams.PrivateKeyID
	params.WitnessPubKeyHashAddrID = litecoinParams.WitnessPubKeyHashAddrID
	params.WitnessScriptHashAddrID = litecoinParams.WitnessScriptHashAddrID
	params.Bech32HRPSegwit = litecoinParams.Bech32HRPSegwit

	copy(params.HDPrivateKeyID[:], litecoinParams.HDPrivateKeyID[:])
	copy(params.HDPublicKeyID[:], litecoinParams.HDPublicKeyID[:])

	params.HDCoinType = litecoinParams.HDCoinType

	checkPoints := make([]chaincfg.Checkpoint, len(litecoinParams.Checkpoints))
	for i := 0; i < len(litecoinParams.Checkpoints); i++ {
		var chainHash chainhash.Hash
		copy(chainHash[:], litecoinParams.Checkpoints[i].Hash[:])

		checkPoints[i] = chaincfg.Checkpoint{
			Height: litecoinParams.Checkpoints[i].Height,
			Hash:   &chainHash,
		}
	}
	params.Checkpoints = checkPoints

	params.RPCPort = litecoinParams.RPCPort
	params.CoinType = litecoinParams.CoinType
}

// IsTestnet tests if the givern params correspond to a testnet
// parameter configuration.
func IsTestnet(params *BrocoinNetParams) bool {
	switch params.Params.Net {
	case brocoinWire.TestNet3, brocoinWire.BrocoinNet(litecoinWire.TestNet4):
		return true
	default:
		return false
	}
}
