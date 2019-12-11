package main

import (
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
	filter          spreadsheet.ParserInput
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
	filter: spreadsheet.ParserInput{
		Path:       "./private-data/filter.xlsx",
		HasHeaders: true,
		RequiredHeaders: []string{
			"claim ref",
			"seemis ID",
		},
	},
}

func main() {
	fsmStore := GenerateFsmAwards(privateInputData)
	ctrStore := GenerateCtrBasedAwards(privateInputData, fsmStore)

	RespondWith(&ctrStore, nil)
}
