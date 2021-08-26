/*
blog_engine.go
Manages the blog for the redcard instance.

@author  Peter Egan
@since   2021-08-17
@lastUpdated 2021-08-22

Copyright (c) 2021 kiercam llc
*/

package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

type blog_post struct {
	title   string
	content string
}

func (b blog_post) post() {
	// Given a valid blog post, write the post to the blog file.
	// 1) Read in the blog file...
	blog, err := ioutil.ReadFile("static/blog.html")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return // Todo - better error handling.
	}

	// 2) Insert the blog post in the file ...
	// a) Grab the datetime
	currentTime := time.Now()

	// b) Covert the blog slice into a string
	blogString := string(blog)

	// c) Create the new HMTL for the blog post
	newPostString := `<div class="post" id="placeholder">
	<h2>%s</h2>
	<h5>%s</h5>
 	<!-- <div class="fakeimg" style="height:200px;">Image</div> -->
	<p>%s</p>
	<hr class="solid">
	</div>`
	newPost := fmt.Sprintf(newPostString, b.title, currentTime.Format("2006-January-02"), b.content)

	// d) Split the blog byte slice.
	s := strings.SplitAfter(blogString, "<div class=\"row\">")

	// e) Recombine: slice1 + new post + slice2
	newBlog := s[0] + "\n\t" + newPost + s[1]

	// 3) Write out the new blog...
	newBlogFile := []byte(newBlog)
	errFile := ioutil.WriteFile("static/blog.html", newBlogFile, 0644)
	if errFile != nil {
		fmt.Printf("err: %v\n", errFile)
		return // Todo - better error handling.
	}
}
