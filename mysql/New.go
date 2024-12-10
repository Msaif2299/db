package mysql

import (
	"database/sql"
	"errors"
	"sync"
	"time"

	"log/slog"
)

type Database struct {
	dbConnection               *sql.DB
	DSN                        string
	Logger                     *slog.Logger
	ShowSQL                    bool
	Timed                      bool
	Lock                       sync.Mutex
	connected                  bool
	MaxDatabaseOpenConnections int
	MaxDatabaseIdleConnections int
	DatabaseIdleTimeout        time.Duration
}

var DB *Database

func New(newDSN string, L *slog.Logger) {

	DB = &Database{
		connected: false,
		DSN:       newDSN,
		Logger:    L,
		ShowSQL:   false,
	}
}

func getConnection() (*sql.DB, error) {

	DB.Lock.Lock()
	// check once more - in case a prev goroutine has established a connection
	if DB.connected && DB.dbConnection != nil {
		DB.Lock.Unlock()
		return DB.dbConnection, nil
	}

	if DB.DSN == "" {
		return nil, errors.New("empty database dsn")
	}

	var err error

	// attempt 3 times to connect, then give up
	for i := 0; i < 3; i++ {
		DB.dbConnection, err = sql.Open("mysql", DB.DSN)

		if err == nil {
			// Open may just validate its arguments without creating a connection to the database.
			// To verify that the data source name is valid, call Ping.
			err = DB.dbConnection.Ping()
			if err == nil {
				break // connection was fine
			}
			DB.Logger.With("attempt", i).With("error", err.Error()).Error("Unable to Ping Database")
			time.Sleep(500 * time.Millisecond) // wait a short while before trying again
			continue
		}
		time.Sleep(500 * time.Millisecond)
		DB.Logger.With("attempt", i).With("error", err.Error()).Error("Unable to Ping Database")
	}

	if err != nil {
		return nil, err
	}

	DB.dbConnection.SetMaxOpenConns(25)
	DB.dbConnection.SetMaxIdleConns(25)
	DB.dbConnection.SetConnMaxIdleTime(5 * time.Minute)
	DB.connected = true
	DB.Lock.Unlock()

	return DB.dbConnection, nil
}

func (db *Database) getConnectionV2() (*sql.DB, error) {

	db.Lock.Lock()
	// check once more - in case a prev goroutine has established a connection
	if db.connected && db.dbConnection != nil {
		db.Lock.Unlock()
		return db.dbConnection, nil
	}

	if db.DSN == "" {
		return nil, errors.New("empty database dsn")
	}

	var err error

	// attempt 3 times to connect, then give up
	for i := 0; i < 3; i++ {
		db.dbConnection, err = sql.Open("mysql", db.DSN)

		if err == nil {
			// Open may just validate its arguments without creating a connection to the database.
			// To verify that the data source name is valid, call Ping.
			err = db.dbConnection.Ping()
			if err == nil {
				break // connection was fine
			}
			db.Logger.With("attempt", i).With("error", err.Error()).Error("Unable to Ping Database")
			time.Sleep(500 * time.Millisecond) // wait a short while before trying again
			continue
		}
		time.Sleep(500 * time.Millisecond)
		db.Logger.With("attempt", i).With("error", err.Error()).Error("Unable to Ping Database")
	}

	if err != nil {
		return nil, err
	}

	db.dbConnection.SetMaxOpenConns(25)
	db.dbConnection.SetMaxIdleConns(25)
	db.dbConnection.SetConnMaxIdleTime(5 * time.Minute)
	db.connected = true
	db.Lock.Unlock()

	return db.dbConnection, nil
}
