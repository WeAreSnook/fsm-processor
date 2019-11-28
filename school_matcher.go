package main

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/addjam/fsm-processor/spreadsheet"
	jellyfish "github.com/jamesturk/go-jellyfish"
)

const definiteMatchThreshold = 0.95

var numComparisons = 0
var totalDependents = 0

type personBySurname []Person

func (v personBySurname) Len() int           { return len(v) }
func (v personBySurname) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v personBySurname) Less(i, j int) bool { return v[i].Surname < v[j].Surname }

type schoolRowBySurname []SchoolRollRow

func (v schoolRowBySurname) Len() int           { return len(v) }
func (v schoolRowBySurname) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v schoolRowBySurname) Less(i, j int) bool { return v[i].Surname < v[j].Surname }

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

	// clean and sort people
	for i, person := range store.People {
		person.Surname = cleanString(person.Surname)
		person.Forename = cleanString(person.Forename)
		person.Postcode = cleanString(person.Postcode)
		person.AddressStreet = cleanString(person.AddressStreet)

		for di, d := range person.Dependents {
			d.Surname = cleanString(d.Surname)
			d.Forename = cleanString(d.Forename)
			d.Person = person
			person.Dependents[di] = d
		}

		store.People[i] = person
	}
	sort.Sort(personBySurname(store.People))

	// Sort the rows
	sort.Sort(schoolRowBySurname(schoolRollRows))

	fmt.Println("Sorted!")
	fmt.Printf("First person %s, row %s\n", store.People[0].Surname, schoolRollRows[0].Surname)
	fmt.Printf("Last person %s, row %s\n",
		store.People[len(store.People)-1].Surname,
		schoolRollRows[len(schoolRollRows)-1].Surname)
	fmt.Printf("%d people\n", len(store.People))

	// Find matches
	var wg sync.WaitGroup
	dependentChannel := make(chan Dependent)

	for _, person := range store.People {
		rowsInPostcode := postcodeIndex[person.Postcode]

		for _, dependent := range person.Dependents {
			wg.Add(1)
			rowsWithSurname := surnameIndex[dependent.Surname]
			totalDependents++
			go checkSchoolRoll(&wg, dependentChannel, dependent, [][]SchoolRollRow{rowsInPostcode, rowsWithSurname, schoolRollRows})
		}
	}

	go func() {
		wg.Wait()
		close(dependentChannel)
	}()

	totalMatched := 0
	for range dependentChannel {
		// fmt.Printf("Match: %#v\n\n", d)
		totalMatched++
	}

	fmt.Printf("%d dependents in NLC out of %d\n", totalMatched, totalDependents)
	fmt.Printf("%d comparisons\n", numComparisons)

	return people, nil
}

func checkSchoolRoll(wg *sync.WaitGroup, ch chan Dependent, d Dependent, rowsToSearch [][]SchoolRollRow) {
	defer wg.Done()

	for _, rows := range rowsToSearch {
		matched := isInSchoolRollRows(d, rows)
		if matched {
			ch <- d
			return
		}
	}
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
	streetScore := compareStrings(takeString(person.AddressStreet, 30), takeString(r.AddressStreet, 30))

	addressScore := math.Max(postcodeScore, streetScore)

	// TODO when do we use street vs postcode
	aggregateScore := calculateWeightedScore(forenameScore, surnameScore, dobScore, addressScore)

	return aggregateScore >= definiteMatchThreshold
}

// Score is weighted twice for dob and postcode
func calculateWeightedScore(forenameScore, surnameScore, dobScore, addressScore float64) float64 {
	return (forenameScore + surnameScore + (2 * dobScore) + (2 * addressScore)) / 6
}

func compareStrings(nameA, nameB string) float64 {
	return jellyfish.JaroWinkler(nameA, nameB)
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
