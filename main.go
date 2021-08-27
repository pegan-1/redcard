/*
main.go
Entry point for the redcard web platform

@author  Peter Egan
@since   2021-08-15
@lastUpdated 2021-08-22

Copyright (c) 2021 kiercam llc
*/

package main

import "fmt"

// Global instances...
var db = DB{}
var webServer = WebServer{}

func main() {
	fmt.Println("Welcome to redcard web platform!")

	// Initialize/Start the Database
	db.start() // Initialize the Database

	// Start listening or http requests.
	webServer.listen()
}
