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

const getRetainedQuote = `
SELECT
	quote_hash,
	deposit_addr,
	signature,
	called_for_user,
	req_liq
FROM retained_quotes
WHERE quote_hash = ?
LIMIT 1`

const insertRetainedQuote = `
INSERT INTO retained_quotes (
    quote_hash,
	deposit_addr,
	signature,
	called_for_user,
	req_liq
)
VALUES (
    :quote_hash,
	:deposit_addr,
	:signature,
	:called_for_user,
	:req_liq
)
`

const setRetainedQuoteCalledForUserFlag = `
UPDATE retained_quotes
SET called_for_user = 1
WHERE quote_hash = :quote_hash
`

const deleteRetainedQuote = `
DELETE FROM retained_quotes
WHERE quote_hash = :quote_hash
`

const selectRetainedQuotes = `
SELECT
	quote_hash,
	deposit_addr,
	signature,
	called_for_user,
	req_liq
FROM retained_quotes`
