package main

import (
	"fsm-processor/people"
	"testing"
)

func TestExtractPeopleWithConsent(t *testing.T) {
	store := people.Store{
		People: []people.Person{},
	}

	inputData := InputData{
		rolloverMode:       false,
		benefitAmountPence: 61000, // £610

		benefitExtractPath: "./testdata/Benefit Extract_06_09_19.txt",
		dependentsSHBEPath: "./testdata/dependants SHBE_06-09-19-2.xlsx",
		hbucdPath:          "./testdata/hb-uc.d-06-09-19.txt",
		fsmCgAwardsPath:    "./testdata/FSM&CGawards_06-09-19.xlsx",
		schoolRollPath:     "./testdata/School Roll Pupil Data_06-09-19-2.xlsx",
		consent360Path:     "./testdata/Consent Report W360.xls",
	}

	t.Run("finds the correct matches", func(t *testing.T) {
		AddPeopleWithConsent(inputData, &store)

		if len(store.People) != 3 {
			t.Errorf("Expected 3 people in store, got %d", len(store.People))
		}
	})
}
