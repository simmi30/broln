package brond_test

import (
	"testing"

	lnwallettest "github.com/brsuite/broln/lnwallet/test"
)

// TestLightningWallet tests LightningWallet powered by brond against our suite
// of interface tests.
func TestLightningWallet(t *testing.T) {
	lnwallettest.TestLightningWallet(t, "brond")
}
