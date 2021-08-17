package main

import (
	"fmt"
	"log"
	"net/http"
)

type WebServer struct{}

func (WebServer) listen() {
	// Set up the route handlers
	http.HandleFunc("/hello", helloHandler)

	// Set up the admin handler...
	http.HandleFunc("/redcard", redcardHandler)

	// Set up the file server to serve the static web pages
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)

	fmt.Println("Starting Web Server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// Route Handlers
// /hello - Placeholder for now.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello There!")
}

// Manages the redcard admin panel
func redcardHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/redcard" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// Process 'Get' and 'Post' calls
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "static/redcard.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseRedcardPost() err: %v", err)
			return
		}
		fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		title := r.FormValue("blog_title")
		post := r.FormValue("blog_body")
		fmt.Fprintf(w, "Title = %s\n", title)
		fmt.Fprintf(w, "Post = %s\n", post)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
