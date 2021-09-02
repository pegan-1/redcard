/*
web_server.go
Web Server for the redcard web platform.

@author  Peter Egan
@since   2021-08-15
@lastUpdated 2021-08-26

Copyright (c) 2021 kiercam llc
*/

package main

import (
	"encoding/json"
	"errors"
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
	http.HandleFunc("/admin", adminHandler)       // Handles the admin panel...
	http.HandleFunc("/admin.html", admin2Handler) // Don't allow to be accessed directly.
	http.HandleFunc("/blog", blogHandler)         // Handles the blog...
	http.HandleFunc("/login", loginHandler)       // Handles the login...

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
		http.Error(w, "404 not found.", http.StatusNotFound)
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

		// TBD: Have the blog post return an error (possibly)

		// After posting the blog, redirect the User to the blog.
		http.Redirect(w, r, "blog.html", http.StatusSeeOther)

		// Reply back to the web
		fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		fmt.Fprintf(w, "Title = %s\n", title)
		fmt.Fprintf(w, "Post = %s\n", post)
	case "PUT": // Here temporarily to handle the Quill
		// Ensure the correct content type has been sent to the server.
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			// Put an error code here. Do nothing but notify the screen.
			fmt.Println("Incorrect Content Type")
		}
		var bp blog_post2 // Declare the blog post.
		var unmarshalErr *json.UnmarshalTypeError
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&bp)
		log.Println("decoding the body!")
		log.Println(bp)
		if err != nil {
			fmt.Println("There was a decoding error!")
			if errors.As(err, &unmarshalErr) {
				errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
			} else {
				errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
			}
			return
		}

		// Add the post to the blog.
		bp.postQuill()

		// Return a success back to the client.
		// Modify this later.
		errorResponse(w, "Success", http.StatusOK)
		return
	default:
		fmt.Fprintf(w, "Sorry, only POST method is supported.")
	}
}

// Only allow admin.html if it is from the login screen (a redirect)
func admin2Handler(w http.ResponseWriter, r *http.Request) {
	// Build the valid referrer string.
	referrerStr := "http://" + r.Host + "/login"
	if r.Referer() == referrerStr {
		// login was successful.  Show the User the admin screen.
		http.ServeFile(w, r, "static/admin.html")
	} else { // User has not logged in. Don't allow access.
		http.Error(w, "404 not found.", http.StatusNotFound)
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
	if r.URL.Path != "/login" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// Process 'Get' and 'Post' calls
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "static/login.html")
	case "POST":
		// The User is attempting to login.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseLoginPost() err: %v", err)
			return
		}

		// Grab the login information from the POST
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Verify the login credentials (very rudimentary for now)
		dbPassword := db.read(username)
		if password == dbPassword { // Credentials match
			// Redirect the User to the admin screen.
			http.Redirect(w, r, "admin.html", http.StatusSeeOther)
		} else { // Credentials do no match.
			// Inform the User that the passwords don't match.
			fmt.Fprintf(w, "Failed Login Attempt. Login credentials are incorrect.")
			// fmt.Fprintf(w, "Attempted User Login! r.PostFrom = %v\n", r.PostForm)
			// fmt.Fprintf(w, "Username = %s\n", username)
			// fmt.Fprintf(w, "Password = %s\n", password)
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

// An error has been encountered in processing HTTP requests. Send a response back
// to the client.
func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
