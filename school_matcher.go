package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/addjam/fsm-processor/spreadsheet"
)

const definiteMatchThreshold = 0.95

var numComparisons = 0

// PeopleWithChildrenAtNlcSchool returns just the people from the store
// that are likely matches for people in the school roll
func PeopleWithChildrenAtNlcSchool(inputData InputData, store PeopleStore) (matched []Dependent, unmatched []Dependent, err error) {
	schoolRollRows, postcodeIndex, surnameIndex, err := cacheSchoolRoll(inputData.schoolRoll, store)
	if err != nil {
		return nil, nil, err
	}

	var wg sync.WaitGroup
	matchChannel := make(chan dependentMatch)

	comparablePeople := cleanPeople(store.People)
	allDependents := []Dependent{}
	for _, person := range comparablePeople {
		rowsInPostcode := postcodeIndex[person.Postcode]

		for _, dependent := range person.Dependents {
			wg.Add(1)
			rowsWithSurname := surnameIndex[dependent.Surname]
			allDependents = append(allDependents, dependent.Dependent)
			go checkSchoolRoll(&wg, matchChannel, dependent, [][]SchoolRollRow{rowsInPostcode, rowsWithSurname, schoolRollRows})
		}
	}

	go func() {
		wg.Wait()
		close(matchChannel)
	}()

	file, err := os.Create("report_fuzzy_matches.csv")
	if err != nil {
		fmt.Println("Error creating output")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"claim number", "seemis ref", "forename SHBE", "forename SEEMIS", "forename score", "surname SHBE", "surname SEEMIS", "surname score", "postcode SHBE", "postcode SEEMIS", "postcode score", "address SHBE", "address SEEMIS", "address score", "dob SHBE", "dob SEEMIS", "dob score", "weighted score"})

	dobString := func(t time.Time) string {
		return t.Format("02-01-06")
	}

	matchedDependents := []Dependent{}
	unmatchedDependents := []Dependent{}
	for match := range matchChannel {
		isMatch := match.Score >= definiteMatchThreshold
		dependent := match.ComparableDependent.Dependent
		if isMatch {
			dependent.SeemisForename = spreadsheet.ColByName(match.Row.OriginalRow, "Forename")
			dependent.SeemisSurname = spreadsheet.ColByName(match.Row.OriginalRow, "Surname")
			dependent.YearGroup = spreadsheet.ColByName(match.Row.OriginalRow, "Year/Stage")
			dependent.SchoolRollRow = match.Row.OriginalRow
			dependent.NameMatchScore = match.NameScore
			dependent.AddressMatchScore = match.AddressScore
			matchedDependents = append(matchedDependents, dependent)
		} else {
			unmatchedDependents = append(unmatchedDependents, dependent)
		}

		if isMatch {
			err := writer.Write([]string{
				fmt.Sprintf("%d", match.ComparableDependent.Dependent.Person.ClaimNumber),
				match.Row.Seemis,
				match.ComparableDependent.Dependent.SeemisForename, match.Row.Forename, fmt.Sprintf("%f", match.ForenameScore),
				match.ComparableDependent.Dependent.SeemisSurname, match.Row.Surname, fmt.Sprintf("%f", match.SurnameScore),
				match.ComparableDependent.Dependent.Person.Postcode, match.Row.Postcode, fmt.Sprintf("%f", match.PostcodeScore),
				match.ComparableDependent.Dependent.Person.AddressStreet, match.Row.AddressStreet, fmt.Sprintf("%f", match.StreetScore),
				dobString(match.ComparableDependent.Dob), dobString(match.Row.Dob), fmt.Sprintf("%f", match.DobScore),
				fmt.Sprintf("%f", match.Score),
			})
			if err != nil {
				fmt.Println("Error Writing line")
			}
		}
	}

	fmt.Printf("%d dependents in NLC, %d unmatched, out of %d total\n", len(matchedDependents), len(unmatchedDependents), len(allDependents))
	fmt.Printf("%d comparisons\n", numComparisons)

	return matchedDependents, unmatchedDependents, nil
}

type dependentMatch struct {
	ComparableDependent comparableDependent
	Score               float64
	Row                 SchoolRollRow
	ForenameScore       float64
	SurnameScore        float64
	NameScore           float64 // Combined forename and surname score (average)
	PostcodeScore       float64
	StreetScore         float64
	AddressScore        float64 // highest of postcode or street score
	DobScore            float64
}

// comparablePerson is a Person with cleaned/normalized fields
type comparablePerson struct {
	Forename      string
	Surname       string
	Postcode      string
	AddressStreet string
	Dependents    []comparableDependent

	Person Person
}

// comparableDependent is a Dependent with cleaned/normalized fields
type comparableDependent struct {
	Forename string
	Surname  string

	// Split Dob into parts on create for quicker comparisons
	Dob      time.Time
	DobYear  int
	DobMonth int
	DobDay   int

	ComparablePerson comparablePerson
	Dependent        Dependent
}

type personBySurname []comparablePerson

func (v personBySurname) Len() int           { return len(v) }
func (v personBySurname) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v personBySurname) Less(i, j int) bool { return v[i].Surname < v[j].Surname }

type schoolRowBySurname []SchoolRollRow

func (v schoolRowBySurname) Len() int           { return len(v) }
func (v schoolRowBySurname) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v schoolRowBySurname) Less(i, j int) bool { return v[i].Surname < v[j].Surname }

func cacheSchoolRoll(input spreadsheet.ParserInput, store PeopleStore) (allRows []SchoolRollRow, postcodeIndex map[string][]SchoolRollRow, surnameIndex map[string][]SchoolRollRow, err error) {
	postcodeIndex = make(map[string][]SchoolRollRow)
	surnameIndex = make(map[string][]SchoolRollRow)
	schoolRollRows := []SchoolRollRow{}
	err = spreadsheet.EachRow(input, func(r spreadsheet.Row) {
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
		return nil, nil, nil, err
	}
	fmt.Printf("Generated postcode index with %d items\n", len(postcodeIndex))
	fmt.Printf("Generated surname index with %d items\n", len(surnameIndex))
	fmt.Printf("Loaded %d items into memory\n", len(schoolRollRows))
	sort.Sort(schoolRowBySurname(schoolRollRows))

	return schoolRollRows, postcodeIndex, surnameIndex, err
}

func checkSchoolRoll(wg *sync.WaitGroup, matchesChan chan dependentMatch, d comparableDependent, rowsToSearch [][]SchoolRollRow) {
	defer wg.Done()

	bestMatch := dependentMatch{
		ComparableDependent: d,
	}
	for _, rows := range rowsToSearch {
		matched, match := isInSchoolRollRows(d, rows)

		if match.Score > bestMatch.Score {
			bestMatch = match
		}

		if matched {
			matchesChan <- match
			return
		}
	}

	matchesChan <- bestMatch
}

func isInSchoolRollRows(d comparableDependent, rows []SchoolRollRow) (bool, dependentMatch) {
	for _, row := range rows {
		matched, match := row.isFuzzyMatch(d.ComparablePerson, d)
		if matched {
			return true, match
		}
	}

	return false, dependentMatch{}
}

// SchoolRollRow represents the columns we care about from the school roll
// it can be used for fuzzy matching
type SchoolRollRow struct {
	Forename      string
	Surname       string
	Postcode      string
	AddressStreet string
	Seemis        string

	// Split DOB into parts on create for dd/mm swap comparisons
	Dob      time.Time
	DobYear  int
	DobMonth int
	DobDay   int

	OriginalRow spreadsheet.Row
}

func cleanedColByName(r spreadsheet.Row, colName string) string {
	rowValue := spreadsheet.ColByName(r, colName)
	return CleanString(rowValue)
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
		Seemis:        spreadsheet.ColByName(r, "SEEMIS reference"),
		OriginalRow:   r,
	}, nil
}

// isFuzzyMatch determins if the dependent/person pair are a match for
// a school roll row
func (r SchoolRollRow) isFuzzyMatch(person comparablePerson, d comparableDependent) (bool, dependentMatch) {
	numComparisons++
	forenameScore := CompareStrings(d.Forename, r.Forename)
	surnameScore := CompareStrings(d.Surname, r.Surname)

	combinedNameScore := (forenameScore + surnameScore) / 2
	if combinedNameScore < 0.7 {
		return false, dependentMatch{}
	}

	dobScore := compareDob(d, r)
	if dobScore == 0 {
		return false, dependentMatch{}
	}

	postcodeScore := CompareStrings(person.Postcode, r.Postcode)

	// We compare only the first 30 characters, as limit in one sheet is 32 and the other is 30
	streetScore := CompareStrings(takeString(person.AddressStreet, 30), takeString(r.AddressStreet, 30))

	// Address score is whichever is highest out of postcode, street
	addressScore := math.Max(postcodeScore, streetScore)

	aggregateScore := calculateWeightedScore(forenameScore, surnameScore, dobScore, addressScore)
	match := aggregateScore >= definiteMatchThreshold

	if match {
		d.Dependent.Seemis = r.Seemis
	}

	return match, dependentMatch{
		ComparableDependent: d,
		Score:               aggregateScore,
		Row:                 r,
		ForenameScore:       forenameScore,
		SurnameScore:        surnameScore,
		PostcodeScore:       postcodeScore,
		StreetScore:         streetScore,
		DobScore:            dobScore,
		NameScore:           combinedNameScore,
		AddressScore:        addressScore,
	}
}

// Score is weighted twice for dob and postcode
func calculateWeightedScore(forenameScore, surnameScore, dobScore, addressScore float64) float64 {
	return (forenameScore + surnameScore + (2 * dobScore) + (2 * addressScore)) / 6
}

// compareDob returns a score of how likely the dob are to be the same
// by allowing for the day/month columns to be switched
func compareDob(d comparableDependent, r SchoolRollRow) float64 {
	// Years must match
	if d.DobYear != r.DobYear {
		return 0
	}

	if d.Dob == r.Dob {
		return 1
	}

	// Year is the same and the month/day is switched
	if d.DobMonth == r.DobDay && d.DobDay == r.DobMonth {
		return 0.9
	}

	return 0
}

func cleanPeople(people []Person) []comparablePerson {
	comparablePeople := []comparablePerson{}
	for _, person := range people {
		p := comparablePerson{
			Surname:       CleanString(person.Surname),
			Forename:      CleanString(person.Forename),
			Postcode:      CleanString(person.Postcode),
			AddressStreet: CleanString(person.AddressStreet),
		}

		for _, dependent := range person.Dependents {
			d := comparableDependent{
				Surname:          CleanString(dependent.Surname),
				Forename:         CleanString(dependent.Forename),
				Dob:              dependent.Dob,
				DobYear:          dependent.Dob.Year(),
				DobMonth:         int(dependent.Dob.Month()),
				DobDay:           dependent.Dob.Day(),
				ComparablePerson: p,
				Dependent:        dependent,
			}

			p.Dependents = append(p.Dependents, d)
		}

		comparablePeople = append(comparablePeople, p)
	}
	sort.Sort(personBySurname(comparablePeople))

	return comparablePeople
}

func takeString(str string, length int) string {
	if len(str) < length {
		return str
	}

	return str[0:length]
}
