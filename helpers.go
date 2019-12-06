package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/addjam/fsm-processor/spreadsheet"
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

// FilterOnlyNewEntitlements filters dependents to ones which have a change in FSM/CG entitlements
func FilterOnlyNewEntitlements(dependents []Dependent) []Dependent {
	withNewEntitlements := []Dependent{}

	for _, d := range dependents {

		if d.HasNewEntitlements() {
			withNewEntitlements = append(withNewEntitlements, d)
		}
	}

	return withNewEntitlements
}

// FilterMinimumP1 returns only the dependents that are in at least P1
func FilterMinimumP1(dependents []Dependent) []Dependent {
	result := []Dependent{}

	for _, d := range dependents {
		if d.IsAtLeastP1() {
			result = append(result, d)
		}
	}

	return result
}

// SplitByMinimumAge splits the dependents into an array >= 16 and an array < 16 years old
func SplitByMinimumAge(inputData InputData, dependents []Dependent) (atThreshold []Dependent, belowThreshold []Dependent) {
	for _, d := range dependents {
		if d.IsAtLeast16(inputData.rolloverMode) {
			atThreshold = append(atThreshold, d)
		} else {
			belowThreshold = append(belowThreshold, d)
		}
	}

	return
}

// FilterUsingExclusionList returns only the dependents that aren't in the filter list
func FilterUsingExclusionList(inputData InputData, dependents []Dependent) []Dependent {
	result := []Dependent{}

	index, err := spreadsheet.CreateIndex(inputData.filter, "claim ref", func(cellValue string) string {
		return cellValue
	})

	if err != nil {
		return result
	}

	for _, d := range dependents {
		claimStr := fmt.Sprintf("%d", d.Person.ClaimNumber)
		isFiltered := false
		if rows, ok := index[claimStr]; ok {
			for _, row := range rows {
				seemis := spreadsheet.ColByName(row, "seemis ID")
				filteredByThisRow := seemis != "" && seemis == d.Seemis
				isFiltered = isFiltered || filteredByThisRow
			}
		}

		if !isFiltered {
			result = append(result, d)
		}
	}

	return result
}
