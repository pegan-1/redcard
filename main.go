/*
main.go
Entry point for the redcard web platform

@author  Peter Egan
@since   2021-08-15
@lastUpdated 2021-08-15

Copyright (c) 2021 kiercam llc
*/

package main

import "fmt"

func main() {
	fmt.Println("Welcome to redcard web platform!")

	// Instantiate a Web Server
	webServer := WebServer{}

	// And start listening for http requests.
	webServer.listen()
}
