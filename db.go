/*
db.go
Manages the redcard database.

@author  Peter Egan
@since   2021-08-22
@lastUpdated 2021-08-23

Copyright (c) 2021 kiercam llc
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type DB struct {
	cache map[string]string
}

// Start the database
func (DB) start() {
	// Instantiate the cache.
	db.cache = map[string]string{}

	// Read in database file.

	// Update cache with database file.

	// One time only - write to database file to bootstrap...
	db.write("admin", "admin")

	fmt.Println("Database running")

}

// Given a key, read the associated value from the database.
func (DB) read(key string) {

}

// Given a key, value -> store the information to the database.
func (DB) write(key string, value string) {
	// Place new entry into the cache.
	db.cache[key] = value

	// Convert the cache to JSON format.
	cacheJSON, err := json.Marshal(db.cache)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Write the updated cache to the database file
	errFile := ioutil.WriteFile("rc.db", cacheJSON, os.ModePerm)
	if errFile != nil {
		fmt.Printf("err: %v\n", errFile)
		return // Todo - better error handling.
	}
}

// Stop the database
func (DB) stop() {

}

// true -> database is running, false -> otherwise
func (DB) isRunning() bool {
	// TBD... Actually create a check that verifies that the db is up.
	return true
}
