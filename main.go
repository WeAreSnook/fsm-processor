package main

import (
	"fmt"
	"fsm-processor/people"
)

// InputData represents all options and files received
type InputData struct {
	// Options
	rolloverMode       bool // when NLC wipes out the data for the previous year and prepares the award for the next school year.
	benefitAmountPence int

	// File paths
	benefitExtractPath string // .txt (formatted as CSV)
	dependentsSHBEPath string // .xlsx
	hbucdPath          string // .txt (formatted as CSV)
	fsmCgAwardsPath    string // .xlsx
	schoolRollPath     string // .xlsx
	consent360Path     string // .xls
}

func main() {
	store := people.Store{}

	inputData := InputData{
		rolloverMode:       false,
		benefitAmountPence: 61000, // £610

		benefitExtractPath: "./private-data/Benefit Extract_06-09-19.txt",
		dependentsSHBEPath: "./private-data/dependants SHBE_06-09-19-2.xlsx",
		hbucdPath:          "./private-data/hb-uc.d-06-09-19.txt",
		fsmCgAwardsPath:    "./private-data/FSM&CGawards_06-09-19.xlsx",
		schoolRollPath:     "./private-data/School Roll Pupil Data_06-09-19-2.xlsx",
		consent360Path:     "./private-data/Consent Report W360.xls",
	}

	AddPeopleWithConsent(inputData, &store)
	fmt.Printf("%d people with consent\n", len(store.People))
	store.People = PeopleInHouseholdsWithChildren(inputData, store)
	fmt.Printf("%d people after household check\n", len(store.People))

	people.RespondWith(store, nil)
}
