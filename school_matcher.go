package main

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/addjam/fsm-processor/spreadsheet"
	jellyfish "github.com/jamesturk/go-jellyfish"
)

const definiteMatchThreshold = 0.95

var numComparisons = 0

// PeopleWithChildrenAtNlcSchool returns just the people from the store
// that are likely matches for people in the school roll
func PeopleWithChildrenAtNlcSchool(inputData InputData, store PeopleStore) ([]Person, error) {
	people := []Person{}

	// Index of postcode to []Row
	postcodeIndex := make(map[string][]SchoolRollRow)
	surnameIndex := make(map[string][]SchoolRollRow)
	schoolRollRows := []SchoolRollRow{}
	err := spreadsheet.EachRow(inputData.schoolRoll, func(r spreadsheet.Row) {
		row, err := NewSchoolRollRow(r)
		if err != nil {
			return
		}

		// Full cache
		schoolRollRows = append(schoolRollRows, row)

		// By Postcode
		if postcodeIndex[row.Postcode] == nil {
			postcodeIndex[row.Postcode] = []SchoolRollRow{}
		}

		postcodeIndex[row.Postcode] = append(postcodeIndex[row.Postcode], row)

		// By Surname
		if surnameIndex[row.Surname] == nil {
			surnameIndex[row.Surname] = []SchoolRollRow{}
		}

		surnameIndex[row.Surname] = append(surnameIndex[row.Surname], row)
	})
	if err != nil {
		return people, err
	}
	fmt.Printf("Generated postcode index with %d items\n", len(postcodeIndex))
	fmt.Printf("Loaded %d items into memory\n", len(schoolRollRows))

	// Find matches
	var wg sync.WaitGroup
	dependentChannel := make(chan Dependent)

	for _, person := range store.People[0:1000] {
		personPostcode := cleanString(person.Postcode)
		rowsInPostcode := postcodeIndex[personPostcode]

		personSurname := cleanString(person.Surname)
		rowsWithSurname := surnameIndex[personSurname]

		for _, dependent := range person.Dependents {
			wg.Add(1)
			go checkSchoolRoll(&wg, dependentChannel, dependent, rowsInPostcode, rowsWithSurname, schoolRollRows)
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

	fmt.Printf("%d dependents in NLC\n", total)
	fmt.Printf("%d comparisons\n", numComparisons)

	return people, nil
}

func checkSchoolRoll(wg *sync.WaitGroup, ch chan Dependent, d Dependent, rowsInPostcode []SchoolRollRow, rowsWithSurname []SchoolRollRow, allRows []SchoolRollRow) {
	defer wg.Done()

	matched := isInSchoolRollRows(d, rowsInPostcode)
	if matched {
		ch <- d
		return
	}

	matched = isInSchoolRollRows(d, rowsWithSurname)
	if matched {
		fmt.Println("Surname index")
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
	// fmt.Printf("%#v\n", d)
	// fmt.Println("")
}

func isInSchoolRollRows(d Dependent, rows []SchoolRollRow) bool {
	for _, row := range rows {
		if row.isFuzzyMatch(d.Person, d) {
			return true
		}
	}

	return false
}

func findMatchingPerson(wg *sync.WaitGroup, ch chan Dependent, r SchoolRollRow, store PeopleStore) {
	defer wg.Done()

	for _, person := range store.People[0:100] {
		for _, dependent := range person.Dependents {
			if r.isFuzzyMatch(person, dependent) {
				ch <- dependent
				return
			}
		}
	}
}

// SchoolRollRow represents the columns we care about from the school roll
// it can be used for fuzzy matching
type SchoolRollRow struct {
	Forename      string
	Surname       string
	Postcode      string
	AddressStreet string

	// Split DOB into parts on create for quicker comparisons
	Dob      time.Time
	DobYear  int
	DobMonth int
	DobDay   int
}

func cleanedColByName(r spreadsheet.Row, colName string) string {
	rowValue := spreadsheet.ColByName(r, colName)
	return cleanString(rowValue)
}

// NewSchoolRollRow creates a SchoolRollRow struct from a row in the school roll spreadsheet
func NewSchoolRollRow(r spreadsheet.Row) (SchoolRollRow, error) {
	dobStr := spreadsheet.ColByName(r, "Date of Birth")
	dob, err := time.Parse("2-Jan-06", dobStr)
	if err != nil {
		return SchoolRollRow{}, err
	}

	return SchoolRollRow{
		Forename:      cleanedColByName(r, "Forename"),
		Surname:       cleanedColByName(r, "Surname"),
		Postcode:      cleanedColByName(r, "Pupil's postcode"),
		AddressStreet: cleanedColByName(r, "Pupil's street"),
		Dob:           dob,
		DobYear:       dob.Year(),
		DobMonth:      int(dob.Month()),
		DobDay:        dob.Day(),
	}, nil
}

// isFuzzyMatch determins if the dependent/person pair are a match for
// a school roll row
func (r SchoolRollRow) isFuzzyMatch(person Person, dependent Dependent) bool {
	numComparisons++
	forenameScore := compareStrings(dependent.Forename, r.Forename)
	surnameScore := compareStrings(dependent.Surname, r.Surname)

	combinedNameScore := (forenameScore + surnameScore) / 2
	if combinedNameScore < 0.7 {
		return false
	}

	dobScore := compareDob(dependent, r)
	if dobScore == 0 {
		return false
	}

	postcodeScore := compareStrings(person.Postcode, r.Postcode)

	// We compare only the first 30 characters, as limit in one sheet is 32 and the other is 30
	// streetScore := compareStrings(takeString(person.AddressStreet, 30), takeString(r.AddressStreet, 30))

	// TODO when do we use street vs postcode
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

// compareDob returns a score of how likely the dob are to be the same
// by allowing for the day/month columns to be switched
func compareDob(d Dependent, r SchoolRollRow) float64 {
	// Years must match
	if d.DobYear != r.DobYear {
		return 0
	}

	if d.Dob == r.Dob {
		return 1
	}

	// Year is the same and the month/day is switched
	if d.DobMonth == r.DobDay && d.DobDay == r.DobMonth {
		// TODO should this be lower than 1? consult Anne
		return 1
	}

	return 0
}

var re *regexp.Regexp = regexp.MustCompile(`[^a-zA-Z\d+]`)

func cleanString(str string) string {
	return strings.ToLower(re.ReplaceAllString(str, ""))
}

func takeString(str string, length int) string {
	if len(str) < length {
		return str
	}

	return str[0:length]
}
