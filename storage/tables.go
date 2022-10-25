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

const createPegoutQuoteTable = `
CREATE TABLE IF NOT EXISTS pegout_quotes (
	hash TEXT PRIMARY KEY,
	lbc_addr TEXT,
	lp_rsk_addr TEXT,
	rsk_refund_addr TEXT,
	fee TEXT,
	penalty_fee TEXT,
	nonce INTEGER,
	value TEXT,
	agreement_timestamp INTEGER,
	deposit_date_limit INTEGER,
	deposit_confirmations INTEGER,
	transfer_confirmations INTEGER,
	transfer_time INTEGER,
	expire_date INTEGER,
	expire_blocks INTEGER,
	derivation_address TEXT
)
`

const createRetainedPegoutQuoteTable = `
CREATE TABLE IF NOT EXISTS retained_pegout_quotes (
	quote_hash TEXT PRIMARY KEY NOT NULL,
	deposit_addr TEXT NOT NULL,
	signature TEXT NOT NULL,
	req_liq TEXT NOT NULL,
	state INTEGER NOT NULL,
	FOREIGN KEY(quote_hash) REFERENCES pegout_quotes(hash)
)
`

const createRetainedPegoutQuoteIndexes = `
CREATE INDEX IF NOT EXISTS retained_pegout_quotes_state_idx
ON retained_pegout_quotes (state)
`
