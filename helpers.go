package main

import (
	"regexp"
	"strings"

	"github.com/jamesturk/go-jellyfish"
)

// CleanRegex is used for cleaning strings by removing punctuation and spaces
var CleanRegex *regexp.Regexp = regexp.MustCompile(`[^a-zA-Z\d+]`)

// CleanString replaces puncutation and spaces, and lowercases the string
func CleanString(str string) string {
	return strings.ToLower(CleanRegex.ReplaceAllString(str, ""))
}

// CompareStrings returns the jaro winkler distance from 0 (no similarity) to 1 (identical) between two strings
func CompareStrings(a, b string) float64 {
	return jellyfish.JaroWinkler(a, b)
}

// CompareCleanedStrings cleans inputs and passes to CompareStrings
func CompareCleanedStrings(a, b string) float64 {
	return CompareStrings(CleanString(a), CleanString(b))
}
