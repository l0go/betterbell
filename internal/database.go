package internal

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
)

// Creates the tables
func CreateTables(db *sql.DB) error {
	// This is where all the different peers are held
	_, err1 := db.Exec("CREATE TABLE IF NOT EXISTS Peers (id INTEGER PRIMARY KEY AUTOINCREMENT, endpoint TEXT, secret TEXT, UNIQUE(endpoint));")

	// Create jobs table
	_, err2 := db.Exec("CREATE TABLE IF NOT EXISTS CronJobs (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT, expression TEXT, enabled BOOLEAN);")

	// Create user table
	_, err3 := db.Exec("CREATE TABLE IF NOT EXISTS Users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, hash TEXT, salt TEXT, UNIQUE(username));")

	// We can now generate a peer that is utilized as this instance's secret
	key, err := generateKey(16)
	if err != nil {
		return fmt.Errorf("Could not generate key: %s", err)
	}
	secret := base64.URLEncoding.EncodeToString(key)

	_, err = db.Exec("INSERT OR IGNORE INTO Peers (endpoint, secret) VALUES(?, ?);", "THIS", secret)
	if err != nil {
		return fmt.Errorf("Could not create peer: %s", err)
	}

	return errors.Join(err1, err2, err3)
}
