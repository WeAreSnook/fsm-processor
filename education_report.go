package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"

	"github.com/addjam/fsm-processor/llog"
)

// GenerateEducationReport generates a spreadsheet of people who were not found in the school roll
func GenerateEducationReport(inputData InputData, store PeopleStore, name string) {
	llog.Printf("Writing education report for %d dependents\n", len(store.ReportForEducationDependents))
	dependents := FilterUsingExclusionList(inputData, store.ReportForEducationDependents)
	llog.Printf("Filtered to %d dependents\n", len(dependents))

	fileName := fmt.Sprintf("report_education_%s.csv", name)
	filePath := path.Join(inputData.outputFolder, fileName)
	file, err := os.Create(filePath)
	llog.Printf("Outputting education report to %s\n", filePath)
	if err != nil {
		llog.Println("Error creating output")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{
		"claim",
		"first name",
		"last name",
		"date of birth",
	})

	for _, d := range dependents {
		dob := d.Dob.Format("02-01-06")

		line := []string{
			fmt.Sprintf("%d", d.Person.ClaimNumber),
			d.Forename,
			d.Surname,
			dob,
		}

		writer.Write(line)
	}
}
