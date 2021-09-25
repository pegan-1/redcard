/*
web_server.go
Web Server for the redcard web platform.

@author  Peter Egan
@since   2021-08-15
@lastUpdated 2021-09-14

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
	http.HandleFunc("/blog", blogHandler)         // Handles the blog summary...
	http.HandleFunc("/login", loginHandler)       // Handles the login...

	// Start up the Web Server
	// fmt.Println("Starting Web Server on port 8080")
	fmt.Println("Web Server running on port 80")
	if err := http.ListenAndServe(":80", nil); err != nil {
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
		// Ensure the correct content type has been sent to the server.
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			// Put an error code here. Do nothing but notify the screen.
			fmt.Println("Incorrect Content Type")
		}
		var bp blog_post // Declare the blog post.
		var unmarshalErr *json.UnmarshalTypeError
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&bp)
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
		// bp.post()

		// PLACEHOLDER for new blog posting.
		// Once in place, will replace the bp.post() with this.
		bp.postToBlog()

		// The blog has posted, send the response back to the User.
		fmt.Fprintf(w, `OK`)

		//TODO: Better error handling.
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
	// fmt.Println("Received a request for /blog")
	// fmt.Println("Request URI " + r.RequestURI)
	// fmt.Println("Method: " + r.Method)

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
