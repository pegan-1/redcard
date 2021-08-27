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
	"log"
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
	dbFile, fileOpenErr := os.Open("rc.db")
	if fileOpenErr != nil {
		log.Fatal(fileOpenErr)
	}
	defer dbFile.Close()

	// Update cache with database file.
	jsonDB, readFileError := ioutil.ReadAll(dbFile)
	if readFileError != nil {
		log.Fatal(readFileError)
	}
	cacheError := json.Unmarshal([]byte(jsonDB), &db.cache)
	if cacheError != nil {
		log.Fatal(cacheError)
	}

	// No errors in loading the database. Report that the db is running.
	fmt.Println("Database running")
}

// Given a key, read the associated value from the database.
func (DB) read(key string) string {
	return db.cache[key]
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
