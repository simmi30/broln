//go:build brocoind && notxindex
// +build brocoind,notxindex

package lntest

import (
	"github.com/brsuite/brond/chaincfg"
)

// NewBackend starts a brocoind node without the txindex enabled and returns a
// BitoindBackendConfig for that node.
func NewBackend(miner string, netParams *chaincfg.Params) (
	*BrocoindBackendConfig, func() error, error) {

	extraArgs := []string{
		"-debug",
		"-regtest",
		"-disablewallet",
	}

	return newBackend(miner, netParams, extraArgs)
}
