package main

import "errors"

// ErrPersonNotFound is returned when a person isn't stored
var ErrPersonNotFound = errors.New("person not found")

// PeopleStorer can store people and perform CRUD operations
type PeopleStorer interface {
	FindExisting(person Person) (Person, error)
	Add(person Person)
	Update(person Person)
}

// PeopleStore is an in-memory PersonStorer
type PeopleStore struct {
	people []Person
}

// Add a Person to the PersonStore
func (p *PeopleStore) Add(person Person) {
	p.people = append(p.people, person)
}

// FindExisting returns a person already in the store with matching details
// This can be useful for finding by name
func (p *PeopleStore) FindExisting(person Person) (Person, error) {
	for _, existingPerson := range p.people {
		if existingPerson.IsSameAs(person) {
			return existingPerson, nil
		}
	}

	return Person{}, ErrPersonNotFound
}

// Update finds an existing matching Person and replaces the entire struct
func (p *PeopleStore) Update(existingPerson, newDetails Person) error {
	for i, existingPerson := range p.people {
		if existingPerson.IsSameAs(existingPerson) {
			p.people[i] = newDetails
			return nil
		}
	}

	return ErrPersonNotFound
}
