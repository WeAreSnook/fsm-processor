package main

// Person represents all the details associated with someone
type Person struct {
	forename    string
	surname     string
	ageYears    int
	claimNumber string
}

// IsSameAs checks if two Person structs refer to the same person
func (p Person) IsSameAs(person Person) bool {
	// TODO jaro-winkler comparison + dob comparison. See notion note.
	return p.forename == person.forename && p.surname == person.surname
}
