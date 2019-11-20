package main

import (
	"log"
	"strconv"

	"github.com/addjam/fsm-processor/spreadsheet"
)

// PeopleInHouseholdsWithChildren returns only the people in the store that belong to households
// which have children, with those children added as dependants.
// Data Source: SHBE
func PeopleInHouseholdsWithChildren(inputData InputData, store PeopleStore) []Person {
	householdPeopleStore := PeopleStore{}

	spreadsheet.EachRow(inputData.dependentsSHBE, func(row spreadsheet.Row) {
		claimNumStr := row.Col(0)
		if claimNumStr == "" {
			return
		}

		claimNumber, err := strconv.Atoi(claimNumStr)

		if err != nil {
			log.Fatalf(`Error parsing claim number "%s" in shbe`, row.Col(0))
		}

		// Check our local store, fall back to the overall store
		person, err := householdPeopleStore.FindByClaimNumber(claimNumber)
		alreadyAdded := err == nil
		if err == ErrPersonNotFound {
			person, err = store.FindByClaimNumber(claimNumber)

			if err == ErrPersonNotFound {
				return
			}
		}

		age, err := strconv.Atoi(row.Col(5))
		if err != nil {
			log.Fatalf("Unable to parse age %d\n", age)
		}

		dependent := Dependent{
			Forename: row.Col(3),
			Surname:  row.Col(2),
			AgeYears: age,
			Dob:      row.Col(4),
		}
		person.AddDependent(dependent)

		if alreadyAdded {
			householdPeopleStore.Update(person)
		} else {
			householdPeopleStore.Add(person)
		}
	})

	return householdPeopleStore.People
}
