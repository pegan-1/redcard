/*
blog_engine.go
Manages the blog for the redcard instance.

@author  Peter Egan
@since   2021-08-17
@lastUpdated 2021-08-17

Copyright (c) 2021 kiercam llc
*/

package main

import "fmt"

type blog_post struct {
	title   string
	content string
}

func (b blog_post) post() {
	fmt.Println("Posting to the Blog!!!")
	fmt.Println(b.title)
	fmt.Println(b.content)

	// // Note: The blog does not exist until the first blog post.
	// if _, err := os.Stat("/path/to/whatever"); err == nil {
	// 	// path/to/whatever exists
	// 	fmt.Println("Blog file already exists!!!")

	// } else if os.IsNotExist(err) {
	// 	// path/to/whatever does *not* exist
	// 	fmt.Println("The ")

	// } else {
	// 	// Schrodinger: file may or may not exist. See err for details.

	// 	// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence

	// }
}
