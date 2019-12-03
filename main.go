package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/addjam/fsm-processor/spreadsheet"
)

// InputData represents all options and files received
type InputData struct {
	// Options
	rolloverMode  bool // when NLC wipes out the data for the previous year and prepares the award for the next school year.
	benefitAmount float32

	// File paths
	benefitExtract  spreadsheet.ParserInput
	dependentsSHBE  spreadsheet.ParserInput
	universalCredit spreadsheet.ParserInput
	fsmCgAwards     spreadsheet.ParserInput
	schoolRoll      spreadsheet.ParserInput
	consent360      spreadsheet.ParserInput
}

var privateInputData = InputData{
	rolloverMode:  false,
	benefitAmount: 610.0, // £610

	benefitExtract: spreadsheet.ParserInput{
		Path:       "./private-data/Benefit Extract_06-09-19.txt",
		HasHeaders: true,
		RequiredHeaders: []string{
			// Extracted in consent check
			"Claim Number",
			"Clmt First Forename",
			"Clmt Surname",

			// Tax credit step one
			"Clmt Personal Pension",
			"Clmt State Retirement Pension (incl SERP's graduated pension etc)",
			"Ptnr Personal Pension",
			"Ptnr State Retirement Pension (incl SERP's graduated pension etc)",
			"Clmt Occupational Pension",
			"Ptnr Occupational Pension",

			// Tax credit step two
			"Clmt AIF",
			"Clmt Employment (gross)",
			"Clmt Self-employment (gross)",
			"Clmt Student Grant/Loan",
			"Clmt Sub-tenants",
			"Clmt Boarders",
			"Clmt Government Training",
			"Clmt Statutory Sick Pay",
			"Clmt Widowed Parent's Allowance",
			"Clmt Apprenticeship",
			"Clmt Statutory Sick Pay",
			"Other weekly Income including In-Work Credit",
			"Ptnr AIF",
			"Ptnr Employment (gross)",
			"Ptnr Self-employment (gross)",
			"Ptnr Student Grant/Loan",
			"Ptnr Sub-tenants",
			"Ptnr Boarders",
			"Ptnr Training for Work/Community Action",
			"Ptnr New Deal 50+ Employment Credit",
			"Ptnr Government Training",
			"Ptnr Carer's Allowance",
			"Ptnr Statutory Sick Pay",
			"Ptnr Widowed Parent's Allowance",
			"Ptnr Apprenticeship",
			"Other weekly Income including In-Work Credit",
			"Clmt Savings Credit",
			"Ptnr Savings Credit",
			"Clmt Widows Benefit",
			"Ptnr Widows Benefit",
		},
	},
	dependentsSHBE: spreadsheet.ParserInput{
		Path:       "./private-data/dependants SHBE_06-09-19-2.xlsx",
		HasHeaders: true,
	},
	universalCredit: spreadsheet.ParserInput{
		Path:       "./private-data/hb-uc.d_06-09-19.txt",
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

func main() {
	store := PeopleStore{}

	err := AddPeopleWithConsent(privateInputData, &store)
	handleErr(err, store)
	fmt.Printf("%d people with consent\n", len(store.People))

	store.People, err = PeopleInHouseholdsWithChildren(privateInputData, store)
	handleErr(err, store)
	fmt.Printf("%d people after household check\n", len(store.People))

	store.People, err = PeopleWithQualifyingIncomes(privateInputData, store)
	handleErr(err, store)
	fmt.Printf("%d people after income qualifying\n", len(store.People))

	nlcDependents, nonNlcDependents, err := PeopleWithChildrenAtNlcSchool(privateInputData, store)
	handleErr(err, store)
	store.ReportForEducationDependents = nonNlcDependents
	store.NlcDependents = nlcDependents
	fmt.Printf("%d dependents in NLC schools, %d unamtched\n", len(nlcDependents), len(nonNlcDependents))

	writeOutput(store)

	RespondWith(&store, nil)
}

func handleErr(err error, store PeopleStore) {
	if err != nil {
		RespondWith(&store, err)
		return
	}
}

// TODO update this to the expected format
func writeOutput(store PeopleStore) {
	file, err := os.Create("report_people.csv")
	if err != nil {
		fmt.Println("Error creating output")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"claim number", "forename", "surname", "addr", "postcode"})
	for _, person := range store.People {
		claimStr := fmt.Sprintf("%d", person.ClaimNumber)
		err := writer.Write([]string{claimStr, person.Forename, person.Surname, person.AddressStreet, person.Postcode})
		if err != nil {
			fmt.Println("Error Writing line")
		}
	}
}
