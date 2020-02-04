package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/addjam/fsm-processor/llog"
	"github.com/addjam/fsm-processor/spreadsheet"
)

// FillExistingGrants iterates over the existing FSM and CG grants
// and adds the data to appropriate dependents
func FillExistingGrants(inputData InputData, dependents []Dependent) []Dependent {
	ninoIndex, err := spreadsheet.CreateIndex(inputData.fsmCgAwards, "NI Number", func(nino string) string {
		return CleanString(nino)
	})

	if err != nil {
		return []Dependent{}
	}

	matches := 0

	var writer *csv.Writer

	if inputData.devMode {
		file, err := os.Create("report_existing_awards_matches.csv")
		if err != nil {
			llog.Println("Error creating output")
		}
		defer file.Close()

		writer = csv.NewWriter(file)
		defer writer.Flush()

		writer.Write([]string{"seemis", "claim", "pupil forename", "award forename", "full forename score", "truncated pupil forename", "truncated award forename", "truncated forename score", "pupil surname", "award surname", "combined score", "truncated combined score"})
	}

	for index, dependent := range dependents {
		nino := CleanString(dependent.Person.Nino)
		awardRows := ninoIndex[nino]

		var bestMatch spreadsheet.Row
		bestMatchScore := 0.0
		bestMatchTruncatedScore := 0.0
		bestMatchForenameScore := 0.0
		bestMatchTruncatedForenameScore := 0.0
		for _, r := range awardRows {
			pupilForename := spreadsheet.ColByName(r, "Pupil Forename")
			pupilSurname := spreadsheet.ColByName(r, "Pupil Surname")
			forenameScore := CompareCleanedStrings(dependent.SeemisForename, pupilForename)
			truncatedForenameScore := CompareCleanedStrings(truncateName(dependent.SeemisForename), truncateName(pupilForename))
			surnameScore := CompareCleanedStrings(dependent.SeemisSurname, pupilSurname)
			combinedScore := (forenameScore + surnameScore) / 2
			truncatedCombinedScore := (truncatedForenameScore + surnameScore) / 2

			if combinedScore > bestMatchScore {
				bestMatch = r
				bestMatchTruncatedScore = truncatedCombinedScore
				bestMatchScore = combinedScore
				bestMatchForenameScore = forenameScore
				bestMatchTruncatedForenameScore = truncatedForenameScore
			}
		}

		isMatch := bestMatchScore >= 0.95

		if isMatch {
			matches++
			fsmGranted := spreadsheet.ColByName(bestMatch, "FSM Approved") != ""
			cgGranted := spreadsheet.ColByName(bestMatch, "Payrun Date") != ""

			d := dependent
			d.ExistingFSM = fsmGranted
			d.ExistingCG = cgGranted
			d.AwardsFsmApproved = spreadsheet.ColByName(bestMatch, "FSM Approved")
			d.AwardsNINumber = spreadsheet.ColByName(bestMatch, "NI Number")
			d.AwardsPayrunDate = spreadsheet.ColByName(bestMatch, "Payrun Date")

			dependents[index] = d
		}

		// Log the best match
		if bestMatchScore > 0 && inputData.devMode {
			writer.Write([]string{
				dependent.Seemis, fmt.Sprintf("%d", dependent.Person.ClaimNumber),
				dependent.SeemisForename, spreadsheet.ColByName(bestMatch, "Pupil Forename"), fmt.Sprintf("%f", bestMatchForenameScore),
				truncateName(dependent.SeemisForename), truncateName(spreadsheet.ColByName(bestMatch, "Pupil Forename")), fmt.Sprintf("%f", bestMatchTruncatedForenameScore),
				dependent.SeemisSurname, spreadsheet.ColByName(bestMatch, "Pupil Surname"),
				fmt.Sprintf("%f", bestMatchScore), fmt.Sprintf("%f", bestMatchTruncatedScore),
			})
		}
	}

	llog.Printf("matched %d out of %d dependents in fsm/cg awards\n", matches, len(dependents))
	return dependents
}

// TODO need to consider combination of forename + surname + _DATE OF BIRTH_,
// but it's not currently in the existing awards spreadsheet. We'll need to add this later.
func findDependentIndex(dependents []Dependent, forename, surname, nino string) (int, error) {
	for i, dep := range dependents {
		forenameScore := CompareCleanedStrings(dep.Forename, forename)
		surnameScore := CompareCleanedStrings(dep.Surname, surname)
		hasSimilarName := ((forenameScore + surnameScore) / 2) > 0.95
		if hasSimilarName && CleanString(dep.Person.Nino) == CleanString(nino) {
			return i, nil
		}
	}
	return -1, ErrPersonNotFound
}

func truncateName(name string) string {
	punctuationToSpaces := CleanRegex.ReplaceAllString(name, " ")
	parts := strings.Split(strings.TrimSpace(punctuationToSpaces), " ")
	return parts[0]
}
