package main

import (
	"testing"

	"github.com/addjam/fsm-processor/spreadsheet"
)

func TestIncomeCheck(t *testing.T) {
	store := PeopleStore{
		People: []Person{},
	}

	performIncomeCheck := func(t *testing.T, inputData InputData, store *PeopleStore) error {
		t.Helper()

		err := AddPeopleWithConsent(inputData, store)
		if err != nil {
			return err
		}

		store.People, err = PeopleInHouseholdsWithChildren(inputData, *store)
		if err != nil {
			return err
		}

		store.People, err = PeopleWithQualifyingIncomes(inputData, *store)
		if err != nil {
			return err
		}

		return nil
	}

	expectedParserInput := spreadsheet.ParserInput{Path: "./private-data/testdata/expected qualifying cases.csv", HasHeaders: true}

	t.Run("extracts the correct number of qualifying people", func(t *testing.T) {
		err := performIncomeCheck(t, privateInputData, &store)
		if err != nil {
			t.Fatalf("Error performing income check: %s", err)
		}

		gotNumClaims := len(store.People)
		wantNumClaims := spreadsheet.CountRows(expectedParserInput)

		if wantNumClaims != gotNumClaims {
			t.Fatalf("Expected %d claims but got %d", wantNumClaims, gotNumClaims)
		}
	})
}
