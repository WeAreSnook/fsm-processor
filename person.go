package main

import (
	"fmt"
	"time"

	"github.com/addjam/fsm-processor/spreadsheet"
)

// Person represents all the details associated with someone
type Person struct {
	Forename      string
	Surname       string
	AgeYears      int
	ClaimNumber   int
	Nino          string
	AddressStreet string
	Postcode      string

	BenefitExtractRow spreadsheet.Row
	Dependents        []Dependent
}

// Dependent represents someone who depends on a Person
type Dependent struct {
	Person   Person
	Forename string
	Surname  string
	AgeYears int
	Dob      time.Time
	Seemis   string

	// Entitlements
	ExistingFSM bool
	ExistingCG  bool
	NewFSM      bool
	NewCG       bool
	Award       string // Just the new award that they aren't registered for. "FSM", "CG", or "Both"
}

func (p Person) String() string {
	return fmt.Sprintf("[Person %s %s, nino: %s, claim no: %d]", p.Forename, p.Surname, p.Nino, p.ClaimNumber)
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

func (d Dependent) String() string {
	return fmt.Sprintf("[Dependent %s %s, seemis: %s, nino: %s, claim no: %d]", d.Forename, d.Surname, d.Seemis, d.Person.Nino, d.Person.ClaimNumber)
}

// HasNewEntitlements returns true if either FSM, CG, or both are now entitlements and weren't previously
func (d Dependent) HasNewEntitlements() bool {
	fsmAdded := !d.ExistingFSM && d.NewFSM
	cgAdded := !d.ExistingCG && d.NewCG

	return fsmAdded || cgAdded
}
