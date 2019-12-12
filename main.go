package main

import (
	"flag"
	"log"

	"github.com/addjam/fsm-processor/spreadsheet"
)

// InputData represents all options and files received
type InputData struct {
	// Options
	rolloverMode  bool // when NLC wipes out the data for the previous year and prepares the award for the next school year.
	benefitAmount float32
	outputFolder  string
	devMode       bool

	// File paths
	benefitExtract  spreadsheet.ParserInput
	dependentsSHBE  spreadsheet.ParserInput
	universalCredit spreadsheet.ParserInput
	fsmCgAwards     spreadsheet.ParserInput
	schoolRoll      spreadsheet.ParserInput
	consent360      spreadsheet.ParserInput
	filter          spreadsheet.ParserInput
}

func main() {
	inputData := parseInputData()

	fsmStore := GenerateFsmAwards(inputData)
	ctrStore := GenerateCtrBasedAwards(inputData, fsmStore)

	RespondWith(&ctrStore, nil)
}

func parseInputData() InputData {
	outputFolderPtr := flag.String("output", "./", "path of the folder outputs should be stored in")
	benefitExtractPtr := flag.String("benefitextract", "", "filepath for benefit extract spreadsheet")
	dependentsSHBEPtr := flag.String("dependents", "", "filepath for dependents SHBE spreadsheet")
	universalCreditPtr := flag.String("universalcredit", "", "filepath for universal credit spreadsheet")
	fsmCgAwardsPtr := flag.String("awards", "", "filepath for current awards spreadsheet")
	schoolRollPtr := flag.String("schoolroll", "", "filepath for school roll spreadsheet")
	consent360Ptr := flag.String("consent", "", "filepath for consent spreadsheet")
	filterPtr := flag.String("filter", "", "filepath for filter spreadsheet")
	rolloverModePtr := flag.Bool("rollover", false, "rollover mode")
	developmentModePtr := flag.Bool("dev", false, "development mode, use private-data")
	benefitAmountPtr := flag.Float64("benefitamount", 610.0, "benefit amount") // default £610
	flag.Parse()

	path := func(inputPath, devModePath string) string {
		var outputPath string
		if inputPath != "" {
			outputPath = inputPath
		} else if *developmentModePtr {
			outputPath = devModePath
		}

		if outputPath == "" {
			log.Fatalf("Error, missing input for path\n")
		}

		return outputPath
	}

	return InputData{
		rolloverMode:  *rolloverModePtr,
		benefitAmount: float32(*benefitAmountPtr),
		outputFolder:  *outputFolderPtr,
		devMode:       *developmentModePtr,

		benefitExtract: spreadsheet.ParserInput{
			Path:       path(*benefitExtractPtr, "./private-data/Benefit Extract_06-09-19.txt"),
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
			Path:       path(*dependentsSHBEPtr, "./private-data/dependants SHBE_06-09-19-2.xlsx"),
			HasHeaders: true,
		},
		universalCredit: spreadsheet.ParserInput{
			Path:       path(*universalCreditPtr, "./private-data/hb-uc.d_06-09-19.txt"),
			HasHeaders: false,
			Format:     spreadsheet.Ssv,
		},
		fsmCgAwards: spreadsheet.ParserInput{
			Path:       path(*fsmCgAwardsPtr, "./private-data/FSM&CGawards_06-09-19.xlsx"),
			HasHeaders: true,
		},
		schoolRoll: spreadsheet.ParserInput{
			Path:       path(*schoolRollPtr, "./private-data/School Roll Pupil Data_06-09-19-2.xlsx"),
			HasHeaders: true,
		},
		consent360: spreadsheet.ParserInput{
			Path:       path(*consent360Ptr, "./private-data/Consent Report W360.xls"),
			HasHeaders: true,
			RequiredHeaders: []string{
				"DocDesc",
				"DocDate",
				"CLAIMREFERENCE",
			},
		},
		filter: spreadsheet.ParserInput{
			Path:       path(*filterPtr, "./private-data/filter.xlsx"),
			HasHeaders: true,
			RequiredHeaders: []string{
				"claim ref",
				"seemis ID",
			},
		},
	}
}
