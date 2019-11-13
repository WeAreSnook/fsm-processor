package main

import (
	"encoding/json"
	"fmt"
	"log"
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

// Output represents the result data
type Output struct {
	Success        bool   `json:"success"`
	OutputFilePath string `json:"output_file_path"`
	DebugData      string `json:"debug,omitempty"`
	Error          string `json:"error,omitempty"`
}

func main() {
	store := PeopleStore{}

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

	// Temporary fake output data for integration with app
	output := Output{
		Success:        true,
		OutputFilePath: "none yet",
		DebugData:      fmt.Sprintf("%d people extracted", len(store.people)),
	}
	json, err := json.Marshal(output)
	if err != nil {
		log.Fatal("Error marshalling json from store")
	}

	fmt.Println(string(json))
}
