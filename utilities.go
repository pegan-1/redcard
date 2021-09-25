/*
utilities.go
Utility functions for redcard.

@author  Peter Egan
@since   2021-09-25
@lastUpdated 2021-09-25

Copyright (c) 2021 kiercam llc
*/

package main

import "strings"

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
