package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/addjam/fsm-processor/spreadsheet"
)

// GenerateAwardList looks at the AwardDependents and generates an award list spreadsheet
func GenerateAwardList(inputData InputData, store PeopleStore) {
	childrenWithNewEntitlements := filterOnlyNewEntitlements(store.AwardDependents)
	fmt.Printf("%d out of %d have new entitlements\n", len(childrenWithNewEntitlements), len(store.AwardDependents))

	childrenInMinimumP1 := filterMinimumP1(childrenWithNewEntitlements)
	fmt.Printf("%d are in at least P1\n", len(childrenInMinimumP1))

	// atLeast16, below16 := splitByMinimumAge(inputData, childrenInMinimumP1)
	// fmt.Printf("%d people at least age 16, %d below\n", len(atLeast16), len(below16))

	// TODO for atLeast16 => waiting for a flag to be added to school roll indicating if they are still in education
	// inEducation := below16

	writeAwardsList(inputData, childrenInMinimumP1)
}

func filterOnlyNewEntitlements(dependents []Dependent) []Dependent {
	withNewEntitlements := []Dependent{}

	for _, d := range dependents {

		if d.HasNewEntitlements() {
			withNewEntitlements = append(withNewEntitlements, d)
		}
	}

	return withNewEntitlements
}

func filterMinimumP1(dependents []Dependent) []Dependent {
	result := []Dependent{}

	for _, d := range dependents {
		if d.IsAtLeastP1() {
			result = append(result, d)
		}
	}

	return result
}

func splitByMinimumAge(inputData InputData, dependents []Dependent) (atThreshold []Dependent, belowThreshold []Dependent) {
	for _, d := range dependents {
		if d.IsAtLeast16(inputData.rolloverMode) {
			atThreshold = append(atThreshold, d)
		} else {
			belowThreshold = append(belowThreshold, d)
		}
	}

	return
}

func writeAwardsList(inputData InputData, dependents []Dependent) {
	fmt.Printf("Writing awards list for %d dependents\n", len(dependents))

	file, err := os.Create("report_awards.csv")
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

	for _, d := range dependents {
		writer.Write(buildLine(inputData, d))
	}
}

func buildLine(inputData InputData, d Dependent) []string {
	line := []string{
		"identifier",
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
		line = append(line, "Absent")
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

	// 	// FSM&CG Awards
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
