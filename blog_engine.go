/*
blog_engine.go
Manages the blog for the redcard instance.

@author  Peter Egan
@since   2021-08-17
@lastUpdated 2021-09-10

Copyright (c) 2021 kiercam llc
*/

package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

type blog_post struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Create the new blog post.
func (b blog_post) post() {
	// Given a valid blog post, write the post to the blog file.
	// 1) Read in the blog file...
	blog, err := ioutil.ReadFile("static/blog.html")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return // Todo - better error handling.
	}

	// 2) Pre-process the blog post.
	content := processImages(b.Content)

	// 3) Insert the blog post in the file ...
	// a) Grab the datetime
	currentTime := time.Now()

	// b) Covert the blog slice into a string
	blogString := string(blog)

	// c) Create the new HMTL for the blog post
	newPostString := `<div class="post">
	<h2>%s</h2>
	<h5>%s</h5>
	%s
	<hr class="solid">
	</div>`
	// newPost := fmt.Sprintf(newPostString, b.Title, currentTime.Format("2006-January-02"), b.Content)
	newPost := fmt.Sprintf(newPostString, b.Title, currentTime.Format("2006-January-02"), content)

	// d) Split the blog byte slice.
	s := strings.SplitAfter(blogString, "<div class=\"blog\">")

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

// Given a new blog post, scan the post and process any images.
func processImages(blogContent string) string {
	// Find images in the blog post and store each image in the slice.
	imageRegex := regexp.MustCompile(`<img(.*?)>`)
	images := imageRegex.FindAllStringSubmatch(blogContent, -1)

	// Save each image and set the image URL
	for _, image := range images {
		// Grab the original URL
		origImageURL := image[0]

		// Grab the image
		imageSlice := strings.Split(origImageURL, ",")
		image := strings.TrimSuffix(imageSlice[1], `">`) //image.

		// And determine the image type
		var imageType string
		if strings.Contains(imageSlice[0], "png") {
			imageType = ".png"
		} else if (strings.Contains(imageSlice[0], "jpeg")) ||
			(strings.Contains(imageSlice[0], "jpg")) {
			imageType = ".jpeg"
		} else {
			panic("Can't process image type " + imageSlice[0])
		}

		// Generate the file name (using current time)
		t := time.Now()
		fileLocation := "static/images/blog/"
		fileName := t.Format(time.RFC3339Nano) + imageType

		// Save the image
		dec, err := base64.StdEncoding.DecodeString(image)
		if err != nil {
			// TODO.. Figure out proper error logging...
			panic(err)
		}
		f, err := os.Create(fileLocation + fileName)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if _, err := f.Write(dec); err != nil {
			panic(err)
		}
		if err := f.Sync(); err != nil {
			panic(err)
		}

		// Determine new URL
		newImageURL := "<img src=images/blog/" + fileName + ">"

		//Replace the original URL with the new URL
		blogContent = strings.Replace(blogContent, origImageURL, newImageURL, -1)
	}

	return (blogContent)
}
