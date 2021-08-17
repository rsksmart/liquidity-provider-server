package storage

import (
	"math/big"

	"github.com/jmoiron/sqlx"
	"github.com/rsksmart/liquidity-provider/types"
	log "github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

const (
	driver     = "sqlite"
	feePos     = 6
	penaltyPos = 7
	valuePos   = 12
)

const insertQuote = `
INSERT INTO quotes (
    hash,
	fed_addr,
	lbc_addr,
	lp_rsk_addr,
	btc_refund_addr,
	rsk_refund_addr,
	lp_btc_addr,
	call_fee,
	penalty_fee,
	contract_addr,
	data,
	gas_limit,
	nonce,
	value,
	agreement_timestamp,
	time_for_deposit,
	call_time,
	confirmations
)
VALUES (
    ?,
	:fed_addr,
	:lbc_addr,
	:lp_rsk_addr,
	:btc_refund_addr,
	:rsk_refund_addr,
	:lp_btc_addr,
	:call_fee,
	:penalty_fee,
	:contract_addr,
	:data,
	:gas_limit,
	:nonce,
	:value,
	:agreement_timestamp,
	:time_for_deposit,
	:call_time,
	:confirmations
)
`
const createTable = `
CREATE TABLE IF NOT EXISTS quotes (
	hash TEXT PRIMARY KEY,
	fed_addr TEXT ,
	lbc_addr TEXT,
	lp_rsk_addr TEXT,
	btc_refund_addr TEXT,
	rsk_refund_addr TEXT,
	lp_btc_addr TEXT,
	call_fee TEXT,
	penalty_fee TEXT,
	contract_addr TEXT,
	data TEXT,
	gas_limit INTEGER,
	nonce INTEGER,
	value TEXT,
	agreement_timestamp INTEGER,
	time_for_deposit INTEGER,
	call_time INTEGER,
	confirmations INTEGER
)
`

type DB struct {
	db *sqlx.DB
}

func Connect(dbPath string) (*DB, error) {
	log.Debug("connecting to DB: ", dbPath)
	db, err := sqlx.Connect(driver, dbPath)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(createTable); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) Close() error {
	log.Debug("closing connection to DB")
	err := db.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) InsertQuote(id string, q *types.Quote) error {
	log.Debug("inserting quote: ", q)
	query, args, _ := sqlx.Named(insertQuote, q)

	callFee := args[feePos].(big.Int)
	penaltyFee := args[penaltyPos].(big.Int)
	value := args[valuePos].(big.Int)
	args[feePos] = callFee.String()
	args[penaltyPos] = penaltyFee.String()
	args[valuePos] = value.String()
	args = append(args, 0)
	copy(args[1:], args)
	args[0] = id

	if _, err := db.db.Exec(query, args...); err != nil {
		return err
	}
	return nil
}

func (db *DB) GetQuote(quoteHash string) (*types.Quote, error) {
	log.Debug("retrieving quote: ", quoteHash)
	var quote types.Quote
	err := db.db.Select(quote, "SELECT * FROM quotes where hash = $1", quoteHash)
	if err != nil {
		return nil, err
	}
	return &quote, nil
}
