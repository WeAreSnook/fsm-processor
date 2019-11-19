package main

import (
	"fmt"

	"github.com/addjam/fsm-processor/people"
	"github.com/addjam/fsm-processor/spreadsheet"
)

// InputData represents all options and files received
type InputData struct {
	// Options
	rolloverMode       bool // when NLC wipes out the data for the previous year and prepares the award for the next school year.
	benefitAmountPence int

	// File paths
	benefitExtract spreadsheet.ParserInput
	dependentsSHBE spreadsheet.ParserInput
	hbucd          spreadsheet.ParserInput
	fsmCgAwards    spreadsheet.ParserInput
	schoolRoll     spreadsheet.ParserInput
	consent360     spreadsheet.ParserInput
}

func main() {
	store := people.Store{}

	inputData := InputData{
		rolloverMode:       false,
		benefitAmountPence: 61000, // £610

		benefitExtract: spreadsheet.ParserInput{
			Path:       "./private-data/Benefit Extract_06-09-19.txt",
			HasHeaders: true,
		},
		dependentsSHBE: spreadsheet.ParserInput{
			Path:       "./private-data/dependants SHBE_06-09-19-2.xlsx",
			HasHeaders: true,
		},
		hbucd: spreadsheet.ParserInput{
			Path:       "./private-data/hb-uc.d-06-09-19.txt",
			HasHeaders: false,
			Format:     spreadsheet.Ssv,
		},
		fsmCgAwards: spreadsheet.ParserInput{
			Path:       "./private-data/FSM&CGawards_06-09-19.xlsx",
			HasHeaders: true,
		},
		schoolRoll: spreadsheet.ParserInput{
			Path:       "./private-data/School Roll Pupil Data_06-09-19-2.xlsx",
			HasHeaders: true,
		},
		consent360: spreadsheet.ParserInput{
			Path:       "./private-data/Consent Report W360.xls",
			HasHeaders: true,
			RequiredHeaders: []string{
				"DocDesc",
				"DocDate",
				"CLAIMREFERENCE",
			},
		},
	}

	AddPeopleWithConsent(inputData, &store)
	fmt.Printf("%d people with consent\n", len(store.People))
	store.People = PeopleInHouseholdsWithChildren(inputData, store)
	fmt.Printf("%d people after household check\n", len(store.People))

	people.RespondWith(store, nil)
}
