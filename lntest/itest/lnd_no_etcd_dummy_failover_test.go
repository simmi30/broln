//go:build !kvdb_etcd
// +build !kvdb_etcd

package itest

import (
	"github.com/brsuite/broln/lntest"
)

// testEtcdFailover is an empty itest when broln is not compiled with etcd
// support.
func testEtcdFailover(net *lntest.NetworkHarness, ht *harnessTest) {}
