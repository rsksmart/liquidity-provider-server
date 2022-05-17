package storage

const enableForeignKeysCheck = `
PRAGMA foreign_keys = ON
`

const createQuoteTable = `
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
	confirmations INTEGER,
	call_on_register INTEGER
)
`

const createRetainedQuoteTable = `
CREATE TABLE IF NOT EXISTS retained_quotes (
	quote_hash TEXT PRIMARY KEY NOT NULL,
	deposit_addr TEXT NOT NULL,
	signature TEXT NOT NULL,
	req_liq TEXT NOT NULL,
	state INTEGER NOT NULL,
	FOREIGN KEY(quote_hash) REFERENCES quotes(hash)
)
`

const createRetainedQuoteIndexes = `
CREATE INDEX IF NOT EXISTS retained_quotes_state_idx
ON retained_quotes (state)
`
