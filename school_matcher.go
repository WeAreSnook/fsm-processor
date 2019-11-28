package main

import (
	"fmt"
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

	// Index of postcode to []Row
	postcodeIndex := make(map[string][]spreadsheet.Row)
	schoolRollRows := []spreadsheet.Row{}
	err := spreadsheet.EachRow(inputData.schoolRoll, func(r spreadsheet.Row) {
		postcode := spreadsheet.ColByName(r, "Pupil's postcode")
		postcode = cleanString(postcode)
		if len(postcode) > 4 {
			postcode = postcode[0:4]
		}
		if postcodeIndex[postcode] == nil {
			postcodeIndex[postcode] = []spreadsheet.Row{}
		}

		postcodeIndex[postcode] = append(postcodeIndex[postcode], r)
		schoolRollRows = append(schoolRollRows, r)
	})
	if err != nil {
		return people, err
	}
	fmt.Printf("Generated index with %d items\n", len(postcodeIndex))
	fmt.Printf("Loaded %d items into memory\n", len(schoolRollRows))

	// Find matches
	var wg sync.WaitGroup
	dependentChannel := make(chan Dependent)

	for _, person := range store.People[0:1000] {
		personPostcode := cleanString(person.Postcode)
		if len(personPostcode) > 4 {
			personPostcode = personPostcode[0:4]
		}
		rowsInPostcode := postcodeIndex[personPostcode]

		for _, dependent := range person.Dependents {
			wg.Add(1)
			go checkSchoolRoll(&wg, dependentChannel, dependent, rowsInPostcode, schoolRollRows)
		}
	}

	// TODO with unmatched dependent -> loop over every row and do a check against them.
	// Only useful if unmatched dependents is a small number

	// err := spreadsheet.EachRow(inputData.schoolRoll, func(r spreadsheet.Row) {
	// 	wg.Add(1)
	// 	go findMatchingPerson(&wg, dependentChannel, r, store)
	// })

	// if err != nil {
	// 	return people, err
	// }

	go func() {
		wg.Wait()
		close(dependentChannel)
	}()

	total := 0
	for range dependentChannel {
		// fmt.Printf("Match: %#v\n\n", d)
		total++
	}

	fmt.Printf("%d dependents in NLC", total)

	return people, nil
}

func checkSchoolRoll(wg *sync.WaitGroup, ch chan Dependent, d Dependent, rowsInPostcode []spreadsheet.Row, allRows []spreadsheet.Row) {
	defer wg.Done()

	matched := isInSchoolRollRows(d, rowsInPostcode)
	if matched {
		ch <- d
		return
	}

	matched = isInSchoolRollRows(d, allRows)
	if matched {
		fmt.Println("From entire roll")
		ch <- d
		return
	}

	fmt.Println("No match")
	fmt.Printf("%#v\n", d)
	fmt.Println("")
}

func isInSchoolRollRows(d Dependent, rows []spreadsheet.Row) bool {
	for _, row := range rows {
		if isFuzzyMatch(d.Person, d, row) {
			return true
		}
	}

	return false
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

	// dobStr := spreadsheet.ColByName(schoolRollRow, "Date of Birth")
	// dob, err := time.Parse("2-Jan-06", dobStr)
	// if err != nil {
	// 	log.Fatalf("Error parsing dob from %s", dobStr)
	// }

	dobScore := 1.0 //compareDates(dependent.Dob, dob)
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
