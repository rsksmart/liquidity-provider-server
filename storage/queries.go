package storage

const selectQuoteByHash = `
SELECT 
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
	confirmations,
    call_on_register
FROM quotes 
WHERE hash = ?
LIMIT 1`

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
	confirmations,
	call_on_register
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
	:confirmations,
    :call_on_register
)
`

const deleteExpiredQuotes = `
DELETE FROM quotes
WHERE hash NOT IN (SELECT quote_hash FROM retained_quotes)
AND agreement_timestamp + time_for_deposit < ?
`

const getRetainedQuote = `
SELECT
	quote_hash,
	deposit_addr,
	signature,
	req_liq,
	state
FROM retained_quotes
WHERE quote_hash = ?
LIMIT 1`

const insertRetainedQuote = `
INSERT INTO retained_quotes (
    quote_hash,
	deposit_addr,
	signature,
	req_liq,
	state
)
VALUES (
    :quote_hash,
	:deposit_addr,
	:signature,
	:req_liq,
	:state
)
`

const updateRetainedQuoteState = `
UPDATE retained_quotes
SET state = :new_state
WHERE quote_hash = :quote_hash AND state = :old_state
`

const selectRetainedQuotes = `
SELECT
	quote_hash,
	deposit_addr,
	signature,
	req_liq,
	state
FROM retained_quotes
WHERE state IN (?)
`

const selectRetainedQuotesReqLiq = `
SELECT
	req_liq
FROM retained_quotes
WHERE state IN (?)
`

const insertPegOutQuote = `
INSERT INTO pegout_quotes (
	hash,
	lbc_addr,
	lp_rsk_addr,
	rsk_refund_addr,
	fee,
	penalty_fee,
	nonce,
	value,
	agreement_timestamp,
	deposit_date_limit,
	deposit_confirmations,
	transfer_confirmations,
	transfer_time,
	expire_date,
	expire_blocks,
    derivation_address
)
VALUES (
    ?,
	:lbc_addr,
	:lp_rsk_addr,
	:rsk_refund_addr,
	:fee,
	:penalty_fee,
	:nonce,
	:value,
	:agreement_timestamp,
	:deposit_date_limit,
	:deposit_confirmations,
	:transfer_confirmations,
    :transfer_time,
    :expire_date,
    :expire_blocks,
    ?
)
`

const insertRetainedPegOutQuote = `
INSERT INTO retained_pegout_quotes (
    quote_hash,
	deposit_addr,
	signature,
	req_liq,
	state
)
VALUES (
    :quote_hash,
	:deposit_addr,
	:signature,
	:req_liq,
	:state
)
`

const updateRetainedPegOutQuoteState = `
UPDATE retained_pegout_quotes
SET state = :new_state
WHERE quote_hash = :quote_hash AND state = :old_state
`

const getRetainedPegOutQuote = `
SELECT
	quote_hash,
	deposit_addr,
	signature,
	req_liq,
	state
FROM retained_pegout_quotes
WHERE quote_hash = ?
LIMIT 1`

const deleteExpiredPegOutQuotes = `
DELETE FROM pegout_quotes
WHERE hash NOT IN (SELECT quote_hash FROM retained_quotes)
AND agreement_timestamp + time_for_deposit < ?
`

const selectPegOutQuoteByHash = `
SELECT 
	lbc_addr,
	lp_rsk_addr,
	rsk_refund_addr,
	fee,
	penalty_fee,
	nonce,
	value,
	agreement_timestamp,
	deposit_date_limit,
	deposit_confirmations,
	transfer_confirmations,
	transfer_time,
	expire_date,
	expire_blocks
FROM pegout_quotes 
WHERE hash = ?
LIMIT 1`
