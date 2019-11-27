package main

import (
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"

	"github.com/addjam/fsm-processor/spreadsheet"
	jellyfish "github.com/jamesturk/go-jellyfish"
)

const definiteMatchThreshold = 0.95

// PeopleWithChildrenAtNlcSchool returns just the people from the store
// that are likely matches for people in the school roll
func PeopleWithChildrenAtNlcSchool(inputData InputData, store PeopleStore) ([]Person, error) {
	people := []Person{}

	var wg sync.WaitGroup
	dependentChannel := make(chan Dependent)

	err := spreadsheet.EachRow(inputData.schoolRoll, func(r spreadsheet.Row) {
		wg.Add(1)
		go findMatchingPerson(&wg, dependentChannel, r, store)
	})

	if err != nil {
		return people, err
	}

	go func() {
		wg.Wait()
		close(dependentChannel)
	}()

	fmt.Printf("%d dependents in NLC", len(dependentChannel))
	// for d := range dependentChannel {
	// 	fmt.Printf("Match: %#v\n\n", d)
	// }

	return people, nil
}

func findMatchingPerson(wg *sync.WaitGroup, ch chan Dependent, r spreadsheet.Row, store PeopleStore) {
	defer wg.Done()

	for _, person := range store.People[0:100] {
		for _, dependent := range person.Dependents {
			if isFuzzyMatch(person, dependent, r) {
				ch <- dependent
				return
			}
		}
	}
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
	forenameScore := compareStrings(dependent.Forename, forename)
	surnameScore := compareStrings(dependent.Surname, surname)
	combinedNameScore := (forenameScore + surnameScore) / 2
	if combinedNameScore < 0.7 {
		// Can never be a definite match, even with full marks from
		// the other scores
		return false
	}

	postcode := cleanedColByName("Pupil's postcode")
	// street := cleanedColByName("Pupil's street")

	dobStr := spreadsheet.ColByName(schoolRollRow, "Date of Birth")
	dob, err := time.Parse("2-Jan-06", dobStr)
	if err != nil {
		log.Fatalf("Error parsing dob from %s", dobStr)
	}

	dobScore := compareDates(dependent.Dob, dob)
	postcodeScore := compareStrings(person.Postcode, postcode)
	// streetScore := compareStrings(person.AddressStreet, street)

	// TODO use weighted score when we have them all
	aggregateScore := calculateWeightedScore(forenameScore, surnameScore, dobScore, postcodeScore)

	// if aggregateScore >= definiteMatchThreshold {
	// 	fmt.Printf("Match of %f\n", aggregateScore)
	// 	fmt.Printf("%s %s\n", dependent.Forename, forename)
	// 	fmt.Printf("%s %s\n", dependent.Surname, surname)
	// 	fmt.Printf("%s %s\n", person.Postcode, postcode)
	// 	fmt.Println()
	// }

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
