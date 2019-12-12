package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"

	"github.com/addjam/fsm-processor/spreadsheet"
)

// GenerateAwardList looks at the AwardDependents and generates an award list spreadsheet
func GenerateAwardList(inputData InputData, store PeopleStore, name string) {
	fmt.Printf("Writing awards list for %d dependents\n", len(store.AwardDependents))

	fileName := fmt.Sprintf("report_awards_%s.csv", name)
	filePath := path.Join(inputData.outputFolder, fileName)
	fmt.Printf("Outputting award list to %s\n", filePath)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating output")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{
		"Record no", "SEEMIS reference",

		// Benefit Extract Columns
		"Claim Number", "NINO", "Clmt First Forename", "Clmt Surname",
		"Ptnr NINO", "Ptnr First Forename", "Ptnr Surname",
		"Address1", "PostCode",

		// Consent360
		"Consent",

		// School Roll
		"Forename", "Surname", "Date of Birth", "Pupil's property",
		"Pupil's street", "Pupil's town", "School Name", "School Name 2",
		"Year/Stage",

		// New
		"Name match", "Address match",

		// FSM&CG Awards
		"NI Number", "Payrun Date",

		// New
		"CG Qualifier",

		// FSM&CG Awards
		"FSM Approved",

		// New
		"FSM Qualifier", "check attendance",
	})

	for _, d := range store.AwardDependents {
		writer.Write(buildLine(inputData, d))
	}
}

var identifier = 0

func buildLine(inputData InputData, d Dependent) []string {
	identifier++
	line := []string{
		fmt.Sprintf("%d", identifier),

		d.Seemis,

		// Benefit Extract Columns
		fmt.Sprintf("%d", d.Person.ClaimNumber),
	}

	benefitColumns := []string{"NINO", "Clmt First Forename", "Clmt Surname",
		"Ptnr NINO", "Ptnr First Forename", "Ptnr Surname",
		"Address1", "PostCode"}
	benefitRow := d.Person.BenefitExtractRow
	for _, colName := range benefitColumns {
		value := spreadsheet.ColByName(benefitRow, colName)
		line = append(line, value)
	}

	// Consent360
	if d.Person.ConsentDesc == "" {
		line = append(line, "Absent")
	} else if d.Person.ConsentDesc == "FSM&CG Consent Removed" {
		line = append(line, "Refused")
	} else {
		line = append(line, "Given")
	}

	// School Roll

	schoolRollRow := d.SchoolRollRow
	line = append(line, spreadsheet.ColByName(schoolRollRow, "Forename"))
	line = append(line, spreadsheet.ColByName(schoolRollRow, "Surname"))
	line = append(line, spreadsheet.ColByName(schoolRollRow, "Date of Birth"))
	line = append(line, spreadsheet.ColByName(schoolRollRow, "Pupil's property"))
	line = append(line, spreadsheet.ColByName(schoolRollRow, "Pupil's street"))
	line = append(line, spreadsheet.ColByName(schoolRollRow, "Pupil's town"))
	line = append(line, spreadsheet.ColByName(schoolRollRow, "School Name"))
	line = append(line, "TODO School name 2")
	line = append(line, spreadsheet.ColByName(schoolRollRow, "Year/Stage"))

	// Name match and address match
	line = append(line, fmt.Sprintf("%f", d.NameMatchScore))
	line = append(line, fmt.Sprintf("%f", d.AddressMatchScore))

	// FSM&CG Awards
	line = append(line, d.AwardsNINumber)
	line = append(line, d.AwardsPayrunDate)

	// New
	if d.NewCG {
		line = append(line, "HB-LCTR IN PAYMENT")
	} else {
		line = append(line, "")
	}

	// FSM&CG Awards
	line = append(line, d.AwardsFsmApproved)

	// New
	line = append(line, d.Person.QualiferType)

	if d.IsAtLeast16(inputData.rolloverMode) {
		line = append(line, "Yes")
	} else {
		line = append(line, "No")
	}

	return line
}
