package routing

import (
	"bytes"
	"fmt"

	"github.com/brsuite/brond/btcec"
	"github.com/brsuite/brond/chaincfg/chainhash"
	"github.com/brsuite/bronutil"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-errors/errors"
	"github.com/brsuite/broln/lnwire"
)

// ValidateChannelAnn validates the channel announcement message and checks
// that node signatures covers the announcement message, and that the brocoin
// signatures covers the node keys.
func ValidateChannelAnn(a *lnwire.ChannelAnnouncement) error {
	// First, we'll compute the digest (h) which is to be signed by each of
	// the keys included within the node announcement message. This hash
	// digest includes all the keys, so the (up to 4 signatures) will
	// attest to the validity of each of the keys.
	data, err := a.DataToSign()
	if err != nil {
		return err
	}
	dataHash := chainhash.DoubleHashB(data)

	// First we'll verify that the passed brocoin key signature is indeed a
	// signature over the computed hash digest.
	brocoinSig1, err := a.BrocoinSig1.ToSignature()
	if err != nil {
		return err
	}
	brocoinKey1, err := btcec.ParsePubKey(a.BrocoinKey1[:], btcec.S256())
	if err != nil {
		return err
	}
	if !brocoinSig1.Verify(dataHash, brocoinKey1) {
		return errors.New("can't verify first brocoin signature")
	}

	// If that checks out, then we'll verify that the second brocoin
	// signature is a valid signature of the brocoin public key over hash
	// digest as well.
	brocoinSig2, err := a.BrocoinSig2.ToSignature()
	if err != nil {
		return err
	}
	brocoinKey2, err := btcec.ParsePubKey(a.BrocoinKey2[:], btcec.S256())
	if err != nil {
		return err
	}
	if !brocoinSig2.Verify(dataHash, brocoinKey2) {
		return errors.New("can't verify second brocoin signature")
	}

	// Both node signatures attached should indeed be a valid signature
	// over the selected digest of the channel announcement signature.
	nodeSig1, err := a.NodeSig1.ToSignature()
	if err != nil {
		return err
	}
	nodeKey1, err := btcec.ParsePubKey(a.NodeID1[:], btcec.S256())
	if err != nil {
		return err
	}
	if !nodeSig1.Verify(dataHash, nodeKey1) {
		return errors.New("can't verify data in first node signature")
	}

	nodeSig2, err := a.NodeSig2.ToSignature()
	if err != nil {
		return err
	}
	nodeKey2, err := btcec.ParsePubKey(a.NodeID2[:], btcec.S256())
	if err != nil {
		return err
	}
	if !nodeSig2.Verify(dataHash, nodeKey2) {
		return errors.New("can't verify data in second node signature")
	}

	return nil

}

// ValidateNodeAnn validates the node announcement by ensuring that the
// attached signature is needed a signature of the node announcement under the
// specified node public key.
func ValidateNodeAnn(a *lnwire.NodeAnnouncement) error {
	// Reconstruct the data of announcement which should be covered by the
	// signature so we can verify the signature shortly below
	data, err := a.DataToSign()
	if err != nil {
		return err
	}

	nodeSig, err := a.Signature.ToSignature()
	if err != nil {
		return err
	}
	nodeKey, err := btcec.ParsePubKey(a.NodeID[:], btcec.S256())
	if err != nil {
		return err
	}

	// Finally ensure that the passed signature is valid, if not we'll
	// return an error so this node announcement can be rejected.
	dataHash := chainhash.DoubleHashB(data)
	if !nodeSig.Verify(dataHash, nodeKey) {
		var msgBuf bytes.Buffer
		if _, err := lnwire.WriteMessage(&msgBuf, a, 0); err != nil {
			return err
		}

		return errors.Errorf("signature on NodeAnnouncement(%x) is "+
			"invalid: %x", nodeKey.SerializeCompressed(),
			msgBuf.Bytes())
	}

	return nil
}

// ValidateChannelUpdateAnn validates the channel update announcement by
// checking (1) that the included signature covers the announcement and has been
// signed by the node's private key, and (2) that the announcement's message
// flags and optional fields are sane.
func ValidateChannelUpdateAnn(pubKey *btcec.PublicKey, capacity bronutil.Amount,
	a *lnwire.ChannelUpdate) error {

	if err := validateOptionalFields(capacity, a); err != nil {
		return err
	}

	return VerifyChannelUpdateSignature(a, pubKey)
}

// VerifyChannelUpdateSignature verifies that the channel update message was
// signed by the party with the given node public key.
func VerifyChannelUpdateSignature(msg *lnwire.ChannelUpdate,
	pubKey *btcec.PublicKey) error {

	data, err := msg.DataToSign()
	if err != nil {
		return fmt.Errorf("unable to reconstruct message data: %v", err)
	}
	dataHash := chainhash.DoubleHashB(data)

	nodeSig, err := msg.Signature.ToSignature()
	if err != nil {
		return err
	}

	if !nodeSig.Verify(dataHash, pubKey) {
		return fmt.Errorf("invalid signature for channel update %v",
			spew.Sdump(msg))
	}

	return nil
}

// validateOptionalFields validates a channel update's message flags and
// corresponding update fields.
func validateOptionalFields(capacity bronutil.Amount,
	msg *lnwire.ChannelUpdate) error {

	if msg.MessageFlags.HasMaxHtlc() {
		maxHtlc := msg.HtlcMaximumMsat
		if maxHtlc == 0 || maxHtlc < msg.HtlcMinimumMsat {
			return errors.Errorf("invalid max htlc for channel "+
				"update %v", spew.Sdump(msg))
		}

		// For light clients, the capacity will not be set so we'll skip
		// checking whether the MaxHTLC value respects the channel's
		// capacity.
		capacityMsat := lnwire.NewMSatFromSatoshis(capacity)
		if capacityMsat != 0 && maxHtlc > capacityMsat {
			return errors.Errorf("max_htlc(%v) for channel "+
				"update greater than capacity(%v)", maxHtlc,
				capacityMsat)
		}
	}

	return nil
}
