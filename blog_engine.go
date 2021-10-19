/*
blog_engine.go
Manages the blog for the redcard instance.

@author  Peter Egan
@since   2021-08-17
@lastUpdated 2021-10-18

Copyright (c) 2021 kiercam llc
*/

package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"image/jpeg"
	"image/png"

	"github.com/nfnt/resize"
)

type blog_post struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Content string `json:"content"`
}

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

	// 5) Add post to the Blog Summary page
	postToBlogSummary(b.Title, b.Summary, content, blogPostTime, url)

	// 6) Add the post the homepage (TODO)
}

// Post the blog to the summary page
func postToBlogSummary(title string, summary string, content string, postTime time.Time, url string) {
	// TODO  - Need to deal with images.

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
	newPostString := `    <div class="post-summary">
        <div class="title">%s</div>
        <div class="date">%s</div>
        <div class="summary">
          <p> %s </p>
        </div>
        <div class="post-link">
          <a href="%s">Read More...</a>
        </div>
        <hr class="solid">
      </div>`
	newPost := fmt.Sprintf(newPostString, title, postTime.Format("January 02, 2006"), summary, url)

	// d) Split the blog byte slice.
	s := strings.SplitAfter(blogString, "<div id=\"blog\">")

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
	postImages := imageRegex.FindAllStringSubmatch(blogContent, -1)

	// Default thumbnail URL (used by default if a blog post doesn't contain an image.)
	// tnURL := "defaultURL" // TODO...

	// Save each image and set the image URL
	for i, postImage := range postImages {
		// Grab the original URL
		origImageURL := postImage[0]

		// Grab the image
		imageSlice := strings.Split(origImageURL, ",")
		base64Image := strings.TrimSuffix(imageSlice[1], `">`) //image.

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
		dec, err := base64.StdEncoding.DecodeString(base64Image)
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

		// Create the thumbnail, save to /static/images/blog location and return url.
		if i == 0 { // Save the first image as a thumbnail.
			// Generate the file location and file name...
			tnLocation := "static/images/blog/"
			tnName := "thumbnail_" + fileName

			// Resize the first image in the blog to use as the thumnail
			tnImagePre, _, err := image.Decode(bytes.NewReader(dec)) // Turn from [] byte to image
			if err != nil {
				fmt.Println("Was not able to resize the image!")
				panic(err)
			}
			tnImagePost := resize.Thumbnail(200, 200, tnImagePre, resize.Lanczos3) // Create thumnail

			// Convert the image to a []byte for saving to file system.
			tnBuf := new(bytes.Buffer)

			// Encode the buffer in the correct image format.
			if imageType == ".png" {
				fmt.Println("Formatting thumbnail in .png")
				err = png.Encode(tnBuf, tnImagePost)
				if err != nil {
					panic(err)
				}
			} else if imageType == ".jpeg" {
				fmt.Println("Formatting thumbnail in .jpeg")
				err = jpeg.Encode(tnBuf, tnImagePost, nil)
				if err != nil {
					panic(err)
				}
			} else {
				panic("Can't process thumbnail of image type " + imageType)
			}

			// Save the thumbnail to the file system...
			f, err := os.Create(tnLocation + tnName)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			//TODO - How to save the image file properly....
			if _, err := f.Write(tnBuf.Bytes()); err != nil {
				panic(err)
			}
			if err := f.Sync(); err != nil {
				panic(err)
			}
			// TODO... Review and determine if the above library dependency can be removed.
			// https://stackoverflow.com/questions/22940724/go-resizing-images
		}
	}
	return (blogContent)
}
