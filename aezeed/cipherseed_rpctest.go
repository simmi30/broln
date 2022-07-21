//go:build rpctest
// +build rpctest

package aezeed

import "github.com/brsuite/bronwallet/waddrmgr"

func init() {
	// For the purposes of our itest, we'll crank down the scrypt params a
	// bit.
	scryptN = waddrmgr.FastScryptOptions.N
	scryptR = waddrmgr.FastScryptOptions.R
	scryptP = waddrmgr.FastScryptOptions.P
}
