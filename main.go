package main

import (
	"flag"

	"github.com/addjam/fsm-processor/llog"

	"github.com/addjam/fsm-processor/spreadsheet"
)

// InputData represents all options and files received
type InputData struct {
	// Debug options
	debugClaimNumber int

	// Options
	rolloverMode  bool // when NLC wipes out the data for the previous year and prepares the award for the next school year.
	awardCG       bool // e.g. might not awarded after about 20th March
	benefitAmount float32
	ctcWtcFigure  float32
	ctcFigure     float32
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

	llog.Printf("Rollover? %t\n", inputData.rolloverMode)

	fsmStore := GenerateFsmAwards(inputData)
	ctrStore := GenerateCtrBasedAwards(inputData, fsmStore)

	RespondWith(&fsmStore, &ctrStore, nil)
}

func parseInputData() InputData {
	outputFolderPtr := flag.String("output", "./", "path of the folder outputs should be stored in")
	debugClaimNumberPtr := flag.Int("debugclaim", -1, "claimnumber to output debug logs for")
	benefitExtractPtr := flag.String("benefitextract", "", "filepath for benefit extract spreadsheet")
	dependentsSHBEPtr := flag.String("dependents", "", "filepath for dependents SHBE spreadsheet")
	universalCreditPtr := flag.String("universalcredit", "", "filepath for universal credit spreadsheet")
	fsmCgAwardsPtr := flag.String("awards", "", "filepath for current awards spreadsheet")
	schoolRollPtr := flag.String("schoolroll", "", "filepath for school roll spreadsheet")
	consent360Ptr := flag.String("consent", "", "filepath for consent spreadsheet")
	filterPtr := flag.String("filter", "", "filepath for filter spreadsheet")
	rolloverModePtr := flag.Bool("rollover", false, "rollover mode")
	awardCGPtr := flag.Bool("awardcg", true, "if we should award CG")
	developmentModePtr := flag.Bool("dev", false, "development mode, use private-data")
	logModePtr := flag.Bool("log", false, "log output to stdout (breaks json output)")
	benefitAmountPtr := flag.Float64("benefitamount", 610.0, "benefit amount")           // default £610
	ctcWtcFigure := flag.Float64("ctcwtcfigure", 6420.0, "ctc/wtc annual income figure") // default £6420
	ctcFigure := flag.Float64("ctcfigure", 16105.0, "ctc annual income figure")          // default £16105
	flag.Parse()

	path := func(inputPath, devModePath string) string {
		var outputPath string
		if inputPath != "" {
			outputPath = inputPath
		} else if *developmentModePtr {
			outputPath = devModePath
		}

		if outputPath == "" {
			RespondWith(nil, nil, ErrInvalidInputPath{filePath: inputPath})
		}

		return outputPath
	}

	llog.PrintToStdout = *logModePtr

	return InputData{
		debugClaimNumber: *debugClaimNumberPtr,

		rolloverMode:  *rolloverModePtr,
		awardCG:       *awardCGPtr,
		benefitAmount: float32(*benefitAmountPtr),
		ctcWtcFigure:  float32(*ctcWtcFigure),
		ctcFigure:     float32(*ctcFigure),
		outputFolder:  *outputFolderPtr,
		devMode:       *developmentModePtr,

		benefitExtract: spreadsheet.ParserInput{
			Path:       path(*benefitExtractPtr, "./private-data/Benefit Extract.txt"),
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
			Path:       path(*dependentsSHBEPtr, "./private-data/dependants SHBE.xlsx"),
			HasHeaders: true,
		},
		universalCredit: spreadsheet.ParserInput{
			Path:       path(*universalCreditPtr, "./private-data/hb-uc.d.txt"),
			HasHeaders: false,
			Format:     spreadsheet.Ssv,
		},
		fsmCgAwards: spreadsheet.ParserInput{
			Path:       path(*fsmCgAwardsPtr, "./private-data/Current Year Awards.xlsx"),
			HasHeaders: true,
		},
		schoolRoll: spreadsheet.ParserInput{
			Path:       path(*schoolRollPtr, "./private-data/School Roll.xlsx"),
			HasHeaders: true,
		},
		consent360: spreadsheet.ParserInput{
			Path:       path(*consent360Ptr, "./private-data/Consent Report.xls"),
			HasHeaders: true,
			RequiredHeaders: []string{
				"DocDesc",
				"DocDate",
				"CLAIMREFERENCE",
			},
		},
		filter: spreadsheet.ParserInput{
			Path:       path(*filterPtr, "./private-data/Filter File-Test.xlsx"),
			HasHeaders: true,
			RequiredHeaders: []string{
				"claim ref",
				"seemis ID",
			},
		},
	}
}
