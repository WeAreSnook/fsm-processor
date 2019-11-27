package main

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/addjam/fsm-processor/spreadsheet"
	jellyfish "github.com/jamesturk/go-jellyfish"
)

const definiteMatchThreshold = 0.95

// PeopleWithChildrenAtNlcSchool returns just the people from the store
// that are likely matches for people in the school roll
func PeopleWithChildrenAtNlcSchool(inputData InputData, store PeopleStore) ([]Person, error) {
	people := []Person{}

	err := spreadsheet.EachRow(inputData.schoolRoll, func(r spreadsheet.Row) {
		for _, person := range store.People {
			for _, dependent := range person.Dependents {
				// TODO concurrently
				// TODO skip already matched?

				inNlc := isFuzzyMatch(person, dependent, r)
				if inNlc {
					dependent.InNlcSchool = true
					break
				}
			}
		}

		fmt.Println("Next row")
	})

	if err != nil {
		return people, err
	}

	return people, nil
}

// isFuzzyMatch determins if the dependent/person pair are a match for
// a school roll row
func isFuzzyMatch(person Person, dependent Dependent, schoolRollRow spreadsheet.Row) bool {
	cleanedColByName := func(colName string) string {
		rowValue := spreadsheet.ColByName(schoolRollRow, colName)
		return cleanString(rowValue)
	}

	forename := cleanedColByName("Forename")
	surname := cleanedColByName("Surname")
	postcode := cleanedColByName("Pupil's postcode")
	// street := cleanedColByName("Pupil's street")

	dobStr := spreadsheet.ColByName(schoolRollRow, "Date of Birth")
	dob, err := time.Parse("2-Jan-06", dobStr)
	if err != nil {
		log.Fatalf("Error parsing dob from %s", dobStr)
	}

	forenameScore := compareStrings(dependent.Forename, forename)
	surnameScore := compareStrings(dependent.Surname, surname)
	dobScore := compareDates(dependent.Dob, dob)
	postcodeScore := compareStrings(person.Postcode, postcode)
	// streetScore := compareStrings(person.AddressStreet, street)

	// TODO use weighted score when we have them all
	aggregateScore := calculateWeightedScore(forenameScore, surnameScore, dobScore, postcodeScore)

	if aggregateScore >= definiteMatchThreshold {
		fmt.Printf("Match of %f\n", aggregateScore)
		fmt.Printf("%s %s\n", dependent.Forename, forename)
		fmt.Printf("%s %s\n", dependent.Surname, surname)
		fmt.Printf("%s %s\n", person.Postcode, postcode)
		fmt.Println()
	}

	return aggregateScore >= definiteMatchThreshold
}

// Score is weighted twice for dob and postcode
func calculateWeightedScore(forenameScore, surnameScore, dobScore, postcodeScore float64) float64 {
	return (forenameScore + surnameScore + (2 * dobScore) + (2 * postcodeScore)) / 6
}

func compareStrings(nameA, nameB string) float64 {
	cleanedA := cleanString(nameA)
	cleanedB := cleanString(nameB)
	return jellyfish.JaroWinkler(cleanedA, cleanedB)
}

func compareDates(dateA, dateB time.Time) float64 {
	return 1.0
}

func cleanString(str string) string {
	re := regexp.MustCompile(`[^a-zA-Z\d+]`)
	return re.ReplaceAllString(str, "")
}
