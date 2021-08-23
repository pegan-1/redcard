/*
web_server.go
Web Server for the redcard web platform.

@author  Peter Egan
@since   2021-08-15
@lastUpdated 2021-08-22

Copyright (c) 2021 kiercam llc
*/

package main

import (
	"fmt"
	"log"
	"net/http"
)

type WebServer struct{}

// listen() sets up all of the route handlers and launches the
// redcard web server.
func (ws WebServer) listen() {
	// Set up the file server to serve the static web pages
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)

	// Set up the route handlers
	http.HandleFunc("/admin", adminHandler) // Handles the admin panel...
	http.HandleFunc("/blog", blogHandler)   // Handles the blog...
	http.HandleFunc("/login", loginHandler) // Handles the login...

	// Start up the Web Server
	// fmt.Println("Starting Web Server on port 8080")
	fmt.Println("Web Server running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// Route Handlers..........
// Handles the redcard admin panel
func adminHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// Process 'Get' and 'Post' calls
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "static/admin.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseAdminPost() err: %v", err)
			return
		}

		// Grab the blog information from the POST
		title := r.FormValue("blog_title")
		post := r.FormValue("blog_body")

		// Add the blog post to the blog.
		bp := blog_post{title: title, content: post}
		bp.post()

		// Reply back to the web
		fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		fmt.Fprintf(w, "Title = %s\n", title)
		fmt.Fprintf(w, "Post = %s\n", post)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

// Handles the Blog
func blogHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/blog" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// Only the GET method is supported.
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	// Serve up the blog (it's that simple!)
	http.ServeFile(w, r, "static/blog.html")
}

// Handles the User Login
func loginHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("The Login Handler has been called!")

	if r.URL.Path != "/login" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// Process 'Get' and 'Post' calls
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "static/login.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseLoginPost() err: %v", err)
			return
		}

		// Grab the login information from the POST
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Verify the login credentials (TBD)
		// Just print for now to see that the server is receiving the information.
		fmt.Println("The username is " + username)
		fmt.Println("The password is " + password)

		// Reply back to the web
		fmt.Fprintf(w, "Attempted User Login! r.PostFrom = %v\n", r.PostForm)
		fmt.Fprintf(w, "Username = %s\n", username)
		fmt.Fprintf(w, "Password = %s\n", password)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
