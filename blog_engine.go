/*
blog_engine.go
Manages the blog for the redcard instance.

@author  Peter Egan
@since   2021-08-17
@lastUpdated 2021-10-10

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
	Summary string `json:"summary"`
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
//   - (TODO) Add to the hompage.
//   - (TODO) Archive current home page/blog summary page before posting blog.
func (b blog_post) postToBlog() {
	// 1) Pre-process the blog post.
	content := processImages(b.Content)

	// 2) Snapshot the blog post time
	blogPostTime := time.Now()
	blogPostTimeString := "Posted on " + blogPostTime.Format("January 02, 2006")

	// 3) Create the stand-alone blog post
	blogPostString := `<html>
	<head>
		<title>%s</title>
			<link rel="stylesheet" type="text/css" href="../css/post.css">
	</head>
	<body>
		<div id="topnav">
			<a href="../index.html">Home</a>
			<a href="../blog">Blog</a>
  		</div>
		<div id="post">
			<div id="title">%s</div>
			<div id="date">%s</div>
			<div id="content">
				%s
			</div>
			<a href="../blog">
				<div id="blog-return">Back To Blog</div>
			</a>
			<hr id="solid">
		</div>
		<div id="footer">
			<p id="footer_logo">powered by redcard</p>
		</div>
	</body>
</html>`

	blogPost := fmt.Sprintf(blogPostString, b.Title, b.Title, blogPostTimeString, content)

	// 4) Save the blog post to its own html file
	// a) Generate file name (convert the title to a file name)
	blogFileName := generateFileName(b.Title)

	// b) Create the blog file
	url := "posts/" + blogFileName + ".html" // URL of new post
	file := "static/" + url                  // File of new post
	blogFile, err := os.Create(file)
	if err != nil {
		// TODO: Come up with better error handling.
		fmt.Printf("Unable to create file: %v", err)
	}

	// c) Write the post to the file
	_, err2 := blogFile.WriteString(blogPost)
	if err2 != nil {
		// TODO: Come up with better error handling.
		fmt.Printf("Unable to write to the file: %v", err)
	}

	// d) Post has been successfully created.
	//    Add the filename to the database.
	payload := "{blog: true}"
	db.write(blogFileName, payload)

	// 5) Add post to the Blog Summary page (TODO)
	postToBlogSummary(b.Title, b.Summary, content, blogPostTime, url)

	// 6) Add the post the homepage (TODO)
}

// Post the blog to the summary page
func postToBlogSummary(title string, summary string, content string, postTime time.Time, url string) {
	// Would like to post the following to the summary page...
	// Picture
	// Date
	// Title
	// Brief Content
	// Read More

	// TODO
	// a) How to come up with a content summary.
	// b) How to manage images.

	// TODO  - Need to deal with images.
	fmt.Println("Post to the Summary")
	fmt.Println(title)
	fmt.Println(summary)
	fmt.Println(content)
	fmt.Println(url)

	// 1) Read in the blog file...
	blog, err := ioutil.ReadFile("static/blog.html")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return // Todo - better error handling.
	}

	// 2) Pre-process the blog post.
	// TODO - Not sure what to do with images.
	// content := processImages(b.Content)

	// 3) Insert the blog summary in the file ...
	// b) Covert the blog slice into a string
	blogString := string(blog)

	// c) Create the new HMTL for the blog post
	newPostString := `<div class="post">
	<h2>%s</h2>
	<h5>%s</h5>
	%s
	<a href="%s">Read More...</a>
	<hr class="solid">
	</div>`
	// newPost := fmt.Sprintf(newPostString, b.Title, currentTime.Format("2006-January-02"), b.Content)
	newPost := fmt.Sprintf(newPostString, title, postTime.Format("2006-January-02"), content, url)

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

// func postBlogToHomepage() {

// }

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
		// fileLocation := "static/images/blog/"
		fileLocation := "static/posts/images/"
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
		// newImageURL := "<img src=images/blog/" + fileName + ">"
		// newImageURL := "<img src=posts/images/" + fileName + ">"
		newImageURL := "<img src=images/" + fileName + ">"

		//Replace the original URL with the new URL
		blogContent = strings.Replace(blogContent, origImageURL, newImageURL, -1)
	}

	return (blogContent)
}
