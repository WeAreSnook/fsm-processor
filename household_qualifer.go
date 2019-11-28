package main

import (
	"log"
	"strconv"
	"time"

	"github.com/addjam/fsm-processor/spreadsheet"
)

// PeopleInHouseholdsWithChildren returns only the people in the store that belong to households
// which have children, with those children added as dependants.
// Data Source: SHBE
func PeopleInHouseholdsWithChildren(inputData InputData, store PeopleStore) ([]Person, error) {
	householdPeopleStore := PeopleStore{}

	err := spreadsheet.EachRow(inputData.dependentsSHBE, func(row spreadsheet.Row) {
		claimNumStr := row.Col(0)
		if claimNumStr == "" {
			return
		}

		claimNumber, err := strconv.Atoi(claimNumStr)

		if err != nil {
			log.Fatalf(`Error parsing claim number "%s" in shbe`, claimNumStr)
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

		dobStr := row.Col(4)
		dob, err := time.Parse("01-02-06", dobStr)
		if err != nil {
			log.Fatalf("Unable to parse dob %s", dobStr)
		}

		dependent := Dependent{
			Forename: row.Col(3),
			Surname:  row.Col(2),
			AgeYears: age,
			Dob:      dob,
			DobYear:  dob.Year(),
			DobMonth: int(dob.Month()),
			DobDay:   dob.Day(),
		}
		person.AddDependent(dependent)

		if alreadyAdded {
			householdPeopleStore.Update(person)
		} else {
			householdPeopleStore.Add(person)
		}
	})

	return householdPeopleStore.People, err
}
