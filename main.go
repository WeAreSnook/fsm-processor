package main

import (
	"fmt"
	"fsm-processor/people"
	"fsm-processor/spreadsheet"
)

// InputData represents all options and files received
type InputData struct {
	// Options
	rolloverMode       bool // when NLC wipes out the data for the previous year and prepares the award for the next school year.
	benefitAmountPence int

	// File paths
	benefitExtractPath spreadsheet.ParserInput // .txt (formatted as CSV)
	dependentsSHBEPath spreadsheet.ParserInput // .xlsx
	hbucdPath          spreadsheet.ParserInput // .txt (formatted as CSV)
	fsmCgAwardsPath    spreadsheet.ParserInput // .xlsx
	schoolRollPath     spreadsheet.ParserInput // .xlsx
	consent360Path     spreadsheet.ParserInput // .xls
}

func main() {
	store := people.Store{}

	inputData := InputData{
		rolloverMode:       false,
		benefitAmountPence: 61000, // £610

		benefitExtractPath: spreadsheet.ParserInput{
			Path:       "./private-data/Benefit Extract_06-09-19.txt",
			HasHeaders: true,
			RequiredHeaders: []string{
				"DocDesc",
				"DocDate",
				"CLAIMREFERENCE",
			},
		},
		dependentsSHBEPath: spreadsheet.ParserInput{
			Path:       "./private-data/dependants SHBE_06-09-19-2.xlsx",
			HasHeaders: true,
		},
		hbucdPath: spreadsheet.ParserInput{
			Path:       "./private-data/hb-uc.d-06-09-19.txt",
			HasHeaders: false,
			Format:     spreadsheet.Ssv,
		},
		fsmCgAwardsPath: spreadsheet.ParserInput{
			Path:       "./private-data/FSM&CGawards_06-09-19.xlsx",
			HasHeaders: true,
		},
		schoolRollPath: spreadsheet.ParserInput{
			Path:       "./private-data/School Roll Pupil Data_06-09-19-2.xlsx",
			HasHeaders: true,
		},
		consent360Path: spreadsheet.ParserInput{
			Path:       "./private-data/Consent Report W360.xls",
			HasHeaders: true,
		},
	}

	AddPeopleWithConsent(inputData, &store)
	fmt.Printf("%d people with consent\n", len(store.People))
	store.People = PeopleInHouseholdsWithChildren(inputData, store)
	fmt.Printf("%d people after household check\n", len(store.People))

	people.RespondWith(store, nil)
}
