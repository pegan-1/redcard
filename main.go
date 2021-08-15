package main

import "fmt"

func main() {
	fmt.Println("Welcome to redcard!")

	// Instantiate a Web Server
	webServer := WebServer{}

	// And start listening for http requests.
	webServer.listen()
}
