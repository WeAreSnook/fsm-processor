package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/addjam/fsm-processor/spreadsheet"
)

// Person represents all the details associated with someone
type Person struct {
	Forename      string
	Surname       string
	ClaimNumber   int
	Nino          string
	AddressStreet string
	Postcode      string

	ConsentDesc  string
	QualiferType string

	BenefitExtractRow spreadsheet.Row
	Dependents        []Dependent
}

// Dependent represents someone who depends on a Person
type Dependent struct {
	Forename  string
	Surname   string
	AgeYears  int
	Dob       time.Time
	Seemis    string
	YearGroup string

	// Entitlements
	ExistingFSM       bool
	ExistingCG        bool
	NewFSM            bool
	NewCG             bool
	AwardsNINumber    string
	AwardsPayrunDate  string
	AwardsFsmApproved string

	// Data from school roll (seemis)
	SchoolRollRow     spreadsheet.Row
	SeemisForename    string
	SeemisSurname     string
	NameMatchScore    float64
	AddressMatchScore float64

	Person Person
}

// NewPersonFromBenefitExtract creates a person based on the provided benefit extract row
func NewPersonFromBenefitExtract(r spreadsheet.Row) (Person, error) {
	claimNumStr := spreadsheet.ColByName(r, "Claim Number")
	claimNumber, err := strconv.Atoi(claimNumStr)

	if err != nil {
		log.Printf("Error parsing claim number from benefits extract %s", claimNumStr)
		return Person{}, err
	}

	return Person{
		Forename:          spreadsheet.ColByName(r, "Clmt First Forename"),
		Surname:           spreadsheet.ColByName(r, "Clmt Surname"),
		ClaimNumber:       claimNumber,
		Postcode:          spreadsheet.ColByName(r, "PostCode"),
		AddressStreet:     spreadsheet.ColByName(r, "Address1"),
		Nino:              spreadsheet.ColByName(r, "NINO"),
		BenefitExtractRow: r,
	}, nil
}

func (p Person) String() string {
	return fmt.Sprintf("[Person %s %s, nino: %s, claim no: %d]", p.Forename, p.Surname, p.Nino, p.ClaimNumber)
}

// ConsentStr returns the consent string based on the given description
func (p Person) ConsentStr() string {
	if p.ConsentDesc == "" {
		return "Absent"
	} else if p.ConsentDesc == "FSM&CG Consent Removed" {
		return "Refused"
	} else {
		return "Given"
	}
}

// HasConsent returns if the person has given consent
func (p Person) HasConsent() bool {
	return p.ConsentStr() == "Given"
}

// AddDependent adds the provided dependent to the Person
func (p *Person) AddDependent(d Dependent) {
	d.Person = *p
	p.Dependents = append(p.Dependents, d)
}

func (d Dependent) String() string {
	return fmt.Sprintf("[Dependent %s %s, seemis: %s, nino: %s, claim no: %d, existing cg: %t, existing fsm: %t, new cg: %t, new fsm: %t]", d.Forename, d.Surname, d.Seemis, d.Person.Nino, d.Person.ClaimNumber, d.ExistingCG, d.ExistingFSM, d.NewCG, d.NewFSM)
}

// HasNewEntitlements returns true if either FSM, CG, or both are now entitlements and weren't previously
func (d Dependent) HasNewEntitlements() bool {
	fsmAdded := !d.ExistingFSM && d.NewFSM
	cgAdded := !d.ExistingCG && d.NewCG

	return fsmAdded || cgAdded
}

// IsAtLeastP1 returns true if the dependent is in a year group P1-S6
// We just check the first character is p or s to allow typos on the number (e.g. S9 is a typo of S6 that has been encountered)
func (d Dependent) IsAtLeastP1() bool {
	firstCharacter := strings.ToLower(string(d.YearGroup[0]))
	return firstCharacter == "s" || firstCharacter == "p"
}

// AgeOn returns the age the person will be on the given date
func (d Dependent) AgeOn(date time.Time) int {
	years := date.Year() - d.Dob.Year()
	if date.YearDay() < d.Dob.YearDay() {
		years--
	}
	return years
}

// IsAtLeast16 determines if the dependent age is >= 16. If rolloverMode is true, their age on the 30th of September is used.
func (d Dependent) IsAtLeast16(rolloverMode bool) bool {
	// If rolloverMode, we consider the age on the 30th of Septemeber. Otherwise, current age.
	now := time.Now()
	ageByDate := now
	if rolloverMode {
		ageByDate = time.Date(now.Year(), 9, 30, 0, 0, 0, 0, time.UTC)
	}

	return d.AgeOn(ageByDate) >= 16
}
