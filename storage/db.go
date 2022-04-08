package storage

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rsksmart/liquidity-provider/types"
	log "github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

const (
	driver = "sqlite"
)

type DBConnector interface {
	CheckConnection() error
	Close() error

	InsertQuote(id string, q *types.Quote) error
	GetQuote(quoteHash string) (*types.Quote, error) // returns nil if not found
	DeleteExpiredQuotes(expTimestamp int64) error

	RetainQuote(entry *types.RetainedQuote) error
	GetRetainedQuotes(filter []types.RQState) ([]*types.RetainedQuote, error)
	GetRetainedQuote(hash string) (*types.RetainedQuote, error) // returns nil if not found
	UpdateRetainedQuoteState(hash string, oldState types.RQState, newState types.RQState) error
	GetLockedLiquidity() (*types.Wei, error)
}

type DB struct {
	db *sqlx.DB
}

type QuoteHash struct {
	QuoteHash string `db:"quote_hash"`
}

type UpdateQuoteState struct {
	QuoteHash string        `db:"quote_hash"`
	OldState  types.RQState `db:"old_state"`
	NewState  types.RQState `db:"new_state"`
}

func (db *DB) CheckConnection() error {
	return db.db.Ping()
}

func Connect(dbPath string) (*DB, error) {
	log.Debug("connecting to DB: ", dbPath)
	db, err := sqlx.Connect(driver, dbPath)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(enableForeignKeysCheck); err != nil {
		return nil, err
	}
	if _, err := db.Exec(createQuoteTable); err != nil {
		return nil, err
	}
	if _, err := db.Exec(createRetainedQuoteTable); err != nil {
		return nil, err
	}
	if _, err := db.Exec(createRetainedQuoteIndexes); err != nil {
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
	log.Debug("inserting quote{", id, "}", ": ", q)
	query, args, _ := sqlx.Named(insertQuote, q)
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
	quote := types.Quote{}
	err := db.db.Get(&quote, selectQuoteByHash, quoteHash)
	switch err {
	case nil:
		return &quote, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (db *DB) DeleteExpiredQuotes(expTimestamp int64) error {
	log.Debug("deleting expired quotes...")
	res, err := db.db.Exec(deleteExpiredQuotes, expTimestamp)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected > 0 {
		log.Infof("deleted %v expired quote(s)", rowsAffected)
	} else {
		log.Debug("no expired quotes found; nothing to delete")
	}

	return nil
}

func (db *DB) RetainQuote(entry *types.RetainedQuote) error {
	log.Debug("inserting retained quote:", entry.QuoteHash)
	query, args, _ := sqlx.Named(insertRetainedQuote, entry)

	_, err := db.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetRetainedQuotes(filter []types.RQState) ([]*types.RetainedQuote, error) {
	log.Debug("retrieving retained quotes")
	var retainedQuotes []*types.RetainedQuote

	query, args, err := sqlx.In(selectRetainedQuotes, filter)
	if err != nil {
		return nil, err
	}
	rows, err := db.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		entry := types.RetainedQuote{}

		err = rows.StructScan(&entry)
		if err != nil {
			return nil, err
		}

		retainedQuotes = append(retainedQuotes, &entry)
	}

	return retainedQuotes, nil
}

func (db *DB) GetRetainedQuote(hash string) (*types.RetainedQuote, error) {
	log.Debug("getting retained quote: ", hash)
	rows, err := db.db.Queryx(getRetainedQuote, hash)
	if err != nil {
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)

	if rows.Next() {
		entry := types.RetainedQuote{}

		err = rows.StructScan(&entry)
		if err != nil {
			return nil, err
		}

		return &entry, nil
	}

	return nil, nil
}

func (db *DB) UpdateRetainedQuoteState(hash string, oldState types.RQState, newState types.RQState) error {
	log.Debugf("updating state from %v to %v for retained quote: %v", oldState, newState, hash)

	query, args, _ := sqlx.Named(updateRetainedQuoteState, UpdateQuoteState{QuoteHash: hash, OldState: oldState, NewState: newState})
	res, err := db.db.Exec(query, args...)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fmt.Errorf("error updating retained quote: %v; oldState: %v; newState: %v", hash, oldState, newState)
	}

	return nil
}

func (db *DB) GetLockedLiquidity() (*types.Wei, error) {
	log.Debug("retrieving locked liquidity")

	filter := []types.RQState{types.RQStateWaitingForDeposit, types.RQStateCallForUserFailed}
	query, args, err := sqlx.In(selectRetainedQuotesReqLiq, filter)
	if err != nil {
		return nil, err
	}
	rows, err := db.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)

	lockedLiq := types.NewWei(0)
	for rows.Next() {
		reqLiq := new(types.Wei)

		err = rows.Scan(reqLiq)
		if err != nil {
			return nil, err
		}

		lockedLiq.Add(lockedLiq, reqLiq)
	}

	return lockedLiq, nil
}
