/*
utilities.go
Utility functions for redcard.

@author  Peter Egan
@since   2021-09-25
@lastUpdated 2021-09-26

Copyright (c) 2021 kiercam llc
*/

package main

import (
	"strconv"
	"strings"
)

// Given a string, generate a file name.
// A file name should be lowercase, alphanumeric. And it should be unique.
func generateFileName(s string) string {
	// 1) Create a lowercase alphanumeric filename, with hyphens substituting for spaces.
	titleMinusPunct := stripPunctuation(s)
	blogFileName := strings.ReplaceAll(strings.ToLower(titleMinusPunct), " ", "-")

	// 2) Ensure that the filename is unique
	// Will simply add an incremental number to the end of the file name.
	if db.doesKeyExist(blogFileName) {
		fileNameExists := true // Stays true until a unique filename is found
		increment := 1         // Increment each time filename is not unique
		for fileNameExists {
			newFileName := blogFileName + "-" + strconv.Itoa(increment)
			if db.doesKeyExist(newFileName) {
				increment++
			} else {
				blogFileName = blogFileName + "-" + strconv.Itoa(increment) // Set new filename
				fileNameExists = false
			}
		}
	}
	return blogFileName
}

// Given a string, strip out all punctuation.
func stripPunctuation(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			result.WriteByte(b)
		}
	}
	return result.String()
}
