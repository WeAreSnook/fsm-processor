package main

import (
	"fmt"

	"github.com/addjam/fsm-processor/spreadsheet"
)

// FillExistingGrants iterates over the existing FSM and CG grants
// and adds the data to appropriate dependents
func FillExistingGrants(inputData InputData, store *PeopleStore) {
	if len(store.AwardDependents) == 0 {
		return
	}

	ninoIndex, err := spreadsheet.CreateIndex(inputData.fsmCgAwards, "NI Number", func(nino string) string {
		return CleanString(nino)
	})

	if err != nil {
		return
	}

	matches := 0

	for index, dependent := range store.AwardDependents {
		nino := CleanString(dependent.Person.Nino)
		awardRows := ninoIndex[nino]

		fmt.Println("===============")
		fmt.Printf("Found %d matching rows for dependent/parent:\n", len(awardRows))
		fmt.Println(dependent)
		fmt.Println(dependent.Person)

		for _, r := range awardRows {
			pupilForename := spreadsheet.ColByName(r, "Pupil Forename")
			pupilSurname := spreadsheet.ColByName(r, "Pupil Surname")
			forenameScore := CompareCleanedStrings(dependent.Forename, pupilForename)
			surnameScore := CompareCleanedStrings(dependent.Surname, pupilSurname)
			combinedScore := (forenameScore + surnameScore) / 2
			isMatch := forenameScore > 0.95 || combinedScore > 0.95

			fmt.Printf("Row for %s %s (%f)\n", pupilForename, pupilSurname, forenameScore)

			if isMatch {
				matches++
				fsmGranted := spreadsheet.ColByName(r, "FSM Approved") != ""
				cgGranted := spreadsheet.ColByName(r, "Payrun Date") != ""

				fmt.Printf("Is match. FSM %t, CG %t\n", fsmGranted, cgGranted)

				dependent.ExistingFSM = fsmGranted
				dependent.ExistingCG = cgGranted
				store.AwardDependents[index] = dependent
				break
			}
		}

		fmt.Println("=============")
	}

	fmt.Printf("matched %d out of %d dependents in fsm/cg awards\n", matches, len(store.AwardDependents))
}

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
