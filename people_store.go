package main

import (
	"errors"
)

// ErrPersonNotFound is returned when a person isn't stored
var ErrPersonNotFound = errors.New("person not found")

// PersonStorer can store people and perform CRUD operations
type PersonStorer interface {
	FindExisting(Person) (Person, error)
	Add(Person)
	Update(Person)
	FindByClaimNumber(int) (Person, error)
}

// PeopleStore is an in-memory PersonStorer
type PeopleStore struct {
	People                       []Person
	ReportForEducationDependents []Dependent
	AwardDependents              []Dependent
}

// Add a Person to the PersonStore
func (p *PeopleStore) Add(person Person) {
	p.People = append(p.People, person)
}

// FindExisting returns a person already in the store with matching details
// This can be useful for finding by name
func (p *PeopleStore) FindExisting(person Person) (Person, error) {
	for _, existingPerson := range p.People {
		if existingPerson.IsSameAs(person) {
			return existingPerson, nil
		}
	}

	return Person{}, ErrPersonNotFound
}

// FindByClaimNumber finds an existing person by the provided claim number
// Returns ErrPersonNotFound if there are no matches
func (p *PeopleStore) FindByClaimNumber(claimNumber int) (Person, error) {
	for _, existingPerson := range p.People {
		if existingPerson.ClaimNumber == claimNumber {
			return existingPerson, nil
		}
	}

	return Person{}, ErrPersonNotFound
}

// Update finds an existing Person by claim number and replaces the entire struct
func (p *PeopleStore) Update(newDetails Person) error {
	for i, existingPerson := range p.People {
		if existingPerson.ClaimNumber == newDetails.ClaimNumber {
			p.People[i] = newDetails
			return nil
		}
	}

	return ErrPersonNotFound
}

// Delete removes the person from the store. Doesn't preserve order.
func (p *PeopleStore) Delete(person Person) error {
	for i, existingPerson := range p.People {
		if existingPerson.ClaimNumber == person.ClaimNumber {
			// Move element to end and truncate
			p.People[i] = p.People[len(p.People)-1]
			p.People = p.People[:len(p.People)-1]
			return nil
		}
	}

	return errors.New("Person doesn't exist")
}
