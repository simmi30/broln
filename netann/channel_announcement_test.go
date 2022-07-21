package netann

import (
	"bytes"
	"testing"

	"github.com/brsuite/brond/chaincfg/chainhash"
	"github.com/brsuite/brond/wire"
	"github.com/brsuite/bronutil"
	"github.com/brsuite/broln/channeldb"
	"github.com/brsuite/broln/lnwire"
	"github.com/stretchr/testify/assert"
)

func TestCreateChanAnnouncement(t *testing.T) {
	t.Parallel()

	key := [33]byte{0x1}
	sig := lnwire.Sig{0x1}
	features := lnwire.NewRawFeatureVector(lnwire.AnchorsRequired)
	var featuresBuf bytes.Buffer
	if err := features.Encode(&featuresBuf); err != nil {
		t.Fatalf("unable to encode features: %v", err)
	}

	expChanAnn := &lnwire.ChannelAnnouncement{
		ChainHash:       chainhash.Hash{0x1},
		ShortChannelID:  lnwire.ShortChannelID{BlockHeight: 1},
		NodeID1:         key,
		NodeID2:         key,
		NodeSig1:        sig,
		NodeSig2:        sig,
		BrocoinKey1:     key,
		BrocoinKey2:     key,
		BrocoinSig1:     sig,
		BrocoinSig2:     sig,
		Features:        features,
		ExtraOpaqueData: []byte{0x1},
	}

	chanProof := &channeldb.ChannelAuthProof{
		NodeSig1Bytes:    expChanAnn.NodeSig1.ToSignatureBytes(),
		NodeSig2Bytes:    expChanAnn.NodeSig2.ToSignatureBytes(),
		BrocoinSig1Bytes: expChanAnn.BrocoinSig1.ToSignatureBytes(),
		BrocoinSig2Bytes: expChanAnn.BrocoinSig2.ToSignatureBytes(),
	}
	chanInfo := &channeldb.ChannelEdgeInfo{
		ChainHash:        expChanAnn.ChainHash,
		ChannelID:        expChanAnn.ShortChannelID.ToUint64(),
		ChannelPoint:     wire.OutPoint{Index: 1},
		Capacity:         bronutil.SatoshiPerBrocoin,
		NodeKey1Bytes:    key,
		NodeKey2Bytes:    key,
		BrocoinKey1Bytes: key,
		BrocoinKey2Bytes: key,
		Features:         featuresBuf.Bytes(),
		ExtraOpaqueData:  expChanAnn.ExtraOpaqueData,
	}
	chanAnn, _, _, err := CreateChanAnnouncement(
		chanProof, chanInfo, nil, nil,
	)
	if err != nil {
		t.Fatalf("unable to create channel announcement: %v", err)
	}

	assert.Equal(t, chanAnn, expChanAnn)
}
