package postgres

import "github.com/brsuite/bronwallet/walletdb"

type Fixture interface {
	DB() walletdb.DB
	Dump() (map[string]interface{}, error)
}
