/*
web_server.go
Web Server for the redcard web platform.

@author  Peter Egan
@since   2021-08-15
@lastUpdated 2021-08-17

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
func (WebServer) listen() {
	// Set up the file server to serve the static web pages
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)

	// Set up the route handlers
	http.HandleFunc("/admin", adminHandler) // Handles the admin panel...
	http.HandleFunc("/blog", blogHandler)   // Handles the blog...

	// Start up the Web Server
	fmt.Println("Starting Web Server on port 8080")
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
