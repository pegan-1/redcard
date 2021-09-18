/*
blog_engine.go
Manages the blog for the redcard instance.

@author  Peter Egan
@since   2021-08-17
@lastUpdated 2021-09-18

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

// Create the new blog post.
// Given a valid blog post -
//	 - Create a file containing the new post.
//   - Add the new post to the summary blog page.
//   - (TBD) Add to the hompage.
//   - (TBD) Archive current home page/blog summary page before posting blog.
func (b blog_post) postNew() {
	db.printKeys()
	// 1) Need to check for duplicate titles...
	// Two possible design choices: (TBD)
	//	- Just check for a duplicate when poster has submitted and add a number to the title. (MVP choice)
	//  - Check in the blog post editor and notify the User if it is a duplicate.

	// 2) Pre-process the blog post.
	content := processImages(b.Content)

	// 3) Snapshot the blog post time
	blogPostTime := time.Now()

	// 4) Create the stand-alone blog post (START HERE NEXT!)
	blogPostString := `<html>
	<head>
		<title>%s</title>
			<link rel="stylesheet" type="text/css" href="css/blog.css">
	</head>
	<body>
		<div class="post">
			<h2>%s</h2>
			<h5>%s</h5>
			%s
			<hr class="solid">
		</div>
		<div id="footer">
			<p id="footer_logo">powered by redcard</p>
		</div>
	</body>
</html>`
	blogPost := fmt.Sprintf(blogPostString, b.Title, b.Title, blogPostTime.Format("2006-January-02"), content)

	// 3) Save the blog post to its own html file
	// a) Generate file name (convert the title to a file name)
	blogFileName := strings.ReplaceAll(strings.ToLower(b.Title), " ", "-")
	fmt.Println("The blog file name is " + blogFileName)

	// b) Create the blog file
	blogFile, err := os.Create("static/blog/" + blogFileName + ".html")
	if err != nil {
		// TODO: Come up with better error handling.
		fmt.Printf("Unable to create file: %v", err)
	}

	// c) Write the post to the file
	n, err := blogFile.WriteString(blogPost)
	if err != nil {
		// TODO: Come up with better error handling.
		fmt.Printf("Unable to write to the file: %v", err)
	}
	fmt.Printf("wrote %d bytes\n", n)

	// 4) Add post to the Blog Summary page (TBD)

	// 5) Add the post the homepage (TBD)

	// Create the new blog post
	// START HERE NEXT!
	// https://stackoverflow.com/questions/46748636/how-to-create-new-file-using-go-script

	// Read in the blog summary page.
	// blog_summary, err := ioutil.ReadFile("static/blog.html")
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// 	return // Todo - better error handling.
	// }

	// 	<html>
	// 	<head>
	// 	  <title>Blog</title>
	// 	  <link rel="stylesheet" type="text/css" href="css/blog.css">
	// 	</head>

	// 	<body>
	// 	  <div class="topnav">
	// 		<a href="./index.html">Home</a>
	// 		<a class="active">Blog</a>
	// 	  </div>
	// 	  <div class="header">
	// 		<!-- <h2>Blog</h2> -->
	// 	  </div>
	// 	  <div class="blog">
	// 	  <div class="post">
	// 	  <h2>Trying to Test the Blog Again</h2>
	// 	  <h5>2021-September-18</h5>
	// 	  <p>Testing the blog name.</p>
	// 	  <hr class="solid">
	// 	  </div>
	// 	  <div class="post">
	// 	  <h2>Testing the Title to Filename</h2>
	// 	  <h5>2021-September-18</h5>
	// 	  <p>Here I am, testing the title to a filename.</p>
	// 	  <hr class="solid">
	// 	  </div>
	// 	  <div class="post">
	// 	  <h2>Here's another test</h2>
	// 	  <h5>2021-September-16</h5>
	// 	  <p>Yet another test!</p>
	// 	  <hr class="solid">
	// 	  </div>
	// 	  <div class="post">
	// 	  <h2>Testing the blog post</h2>
	// 	  <h5>2021-September-16</h5>
	// 	  <p>Will I see the keys?</p>
	// 	  <hr class="solid">
	// 	  </div>
	// 	  <div class="post">
	// 	  <h2>Do I fire the new post?</h2>
	// 	  <h5>2021-September-16</h5>
	// 	  <p>Checking if I fire the new post!</p>
	// 	  <hr class="solid">
	// 	  </div>
	// 	  <div class="post">
	// 	  <h2>Is the blog still working?</h2>
	// 	  <h5>2021-September-15</h5>
	// 	  <p>I believe it is still working.</p>
	// 	  <hr class="solid">
	// 	  </div>
	// 	  <div class="post">
	// 	  <h2>Blog Summary</h2>
	// 	  <h5>2021-September-15</h5>
	// 	  <p>Testing the blog summary page...</p>
	// 	  <hr class="solid">
	// 	  </div></div>
	// 	  <div class="footer">
	// 		<p class="footer_logo">powered by redcard</p>
	// 	  </div>
	// 	</body>
	//   </html>

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
