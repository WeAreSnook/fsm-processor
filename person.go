package main

import (
	"time"

	"github.com/addjam/fsm-processor/spreadsheet"
)

// Person represents all the details associated with someone
type Person struct {
	Forename      string
	Surname       string
	AgeYears      int
	ClaimNumber   int
	AddressStreet string
	Postcode      string
	Dependents    []Dependent

	BenefitExtractRow spreadsheet.Row
}

// Dependent represents someone who depends on a Person
type Dependent struct {
	Person   Person
	Forename string
	Surname  string
	AgeYears int
	Dob      time.Time

	// TODO:
	InNlcSchool bool
}

// IsSameAs checks if two Person structs refer to the same person
func (p Person) IsSameAs(person Person) bool {
	// TODO jaro-winkler comparison + dob comparison. See notion note.
	return p.Forename == person.Forename && p.Surname == person.Surname
}

// AddDependent adds the provided dependent to the Person
func (p *Person) AddDependent(d Dependent) {
	d.Person = *p
	p.Dependents = append(p.Dependents, d)
}
