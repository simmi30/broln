package channeldb

import (
	"testing"

	"github.com/brsuite/broln/kvdb"
)

func TestMain(m *testing.M) {
	kvdb.RunTests(m)
}
