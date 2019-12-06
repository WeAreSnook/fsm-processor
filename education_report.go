package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

// GenerateEducationReport generates a spreadsheet of people who were not found in the school roll
func GenerateEducationReport(inputData InputData, store PeopleStore) {
	fmt.Printf("Writing education report for %d dependents\n", len(store.ReportForEducationDependents))
	dependents := FilterUsingExclusionList(inputData, store.ReportForEducationDependents)
	fmt.Printf("Filtered to %d dependents\n", len(dependents))

	file, err := os.Create("report_education.csv")
	if err != nil {
		fmt.Println("Error creating output")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{
		// "SEEMIS reference",
		// "forename",
		// "surname",
		// "Pupil's postcode",
		// "Pupil's property",
		// "Pupil's street",
		// "Pupil's town",
		// "School Name",
		// "Year/Stage",
		"claim",
	})

	// schoolRollCols := []string{
	// 	"SEEMIS reference",
	// 	"forename",
	// 	"surname",
	// 	"Pupil's postcode",
	// 	"Pupil's property",
	// 	"Pupil's street",
	// 	"Pupil's town",
	// 	"School Name",
	// 	"Year/Stage",
	// }

	for _, d := range dependents {
		line := []string{}
		// fmt.Println("Handling", d)

		// for _, colName := range schoolRollCols {
		// 	line = append(line, spreadsheet.ColByName(d.SchoolRollRow, colName))
		// }

		line = append(line, fmt.Sprintf("%d", d.Person.ClaimNumber))

		writer.Write(line)
	}
}
