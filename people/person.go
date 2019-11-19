package people

// Person represents all the details associated with someone
type Person struct {
	Forename    string
	Surname     string
	AgeYears    int
	ClaimNumber int
	Dependents  []Dependent
}

// Dependent represents someone who depends on a Person
type Dependent struct {
	Forename string
	Surname  string
	AgeYears int
	Dob      string
}

// IsSameAs checks if two Person structs refer to the same person
func (p Person) IsSameAs(person Person) bool {
	// TODO jaro-winkler comparison + dob comparison. See notion note.
	return p.Forename == person.Forename && p.Surname == person.Surname
}

// AddDependent adds the provided dependent to the Person
func (p *Person) AddDependent(d Dependent) {
	p.Dependents = append(p.Dependents, d)
}
