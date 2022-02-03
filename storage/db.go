package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/rsksmart/liquidity-provider/providers"
	"github.com/rsksmart/liquidity-provider/types"
	log "github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

const (
	driver = "sqlite"
)

type DBConnector interface {
	providers.RetainedQuotesRepository

	CheckConnection() error
	Close() error

	InsertQuote(id string, q *types.Quote) error
	GetQuote(quoteHash string) (*types.Quote, error)

	GetRetainedQuotes() ([]*types.RetainedQuote, error)
	SetRetainedQuoteCalledForUserFlag(hash string) error
}

type DB struct {
	db *sqlx.DB
}

type QuoteHash struct {
	QuoteHash string `db:"quote_hash"`
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
	if err != nil {
		return nil, err
	}

	return &quote, nil
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

func (db *DB) GetRetainedQuotes() ([]*types.RetainedQuote, error) {
	log.Debug("retrieving retained quotes")
	var retainedQuotes []*types.RetainedQuote
	rows, err := db.db.Queryx(selectRetainedQuotes)
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

func (db *DB) SetRetainedQuoteCalledForUserFlag(hash string) error {
	log.Debug("setting retained quote's calledForUser flag")
	query, args, _ := sqlx.Named(setRetainedQuoteCalledForUserFlag, QuoteHash{QuoteHash: hash})
	if _, err := db.db.Exec(query, args...); err != nil {
		return err
	}
	return nil
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

func (db *DB) DeleteRetainedQuote(hash string) error {
	log.Debug("deleting retained quote:", hash)

	query, args, _ := sqlx.Named(deleteRetainedQuote, QuoteHash{QuoteHash: hash})
	if _, err := db.db.Exec(query, args...); err != nil {
		return err
	}
	return nil
}
