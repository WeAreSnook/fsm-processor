package main

import (
	"testing"

	"github.com/addjam/fsm-processor/spreadsheet"
)

func TestExtractPeopleWithConsent(t *testing.T) {
	store := PeopleStore{
		People: []Person{},
	}

	inputData := InputData{
		rolloverMode:  false,
		benefitAmount: 610, // £610

		benefitExtract:  spreadsheet.ParserInput{Path: "./testdata/Benefit Extract_06_09_19.txt", HasHeaders: true},
		dependentsSHBE:  spreadsheet.ParserInput{Path: "./testdata/dependants SHBE_06-09-19-2.xlsx"},
		universalCredit: spreadsheet.ParserInput{Path: "./testdata/hb-uc.d-06-09-19.txt"},
		fsmCgAwards:     spreadsheet.ParserInput{Path: "./testdata/FSM&CGawards_06-09-19.xlsx"},
		schoolRoll:      spreadsheet.ParserInput{Path: "./testdata/School Roll Pupil Data_06-09-19-2.xlsx"},
		consent360:      spreadsheet.ParserInput{Path: "./testdata/Consent Report W360.xls"},
	}

	t.Run("finds the correct matches", func(t *testing.T) {
		AddPeopleWithConsent(inputData, &store)

		if len(store.People) != 3 {
			t.Errorf("Expected 3 people in store, got %d", len(store.People))
		}
	})
}
