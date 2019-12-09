package main

import "testing"

func TestAddDependent(t *testing.T) {
	t.Run("adds the dependent", func(t *testing.T) {
		user := Person{Forename: "bob", Surname: "smith"}
		dependent := Dependent{Forename: "d", Surname: "smith"}

		if len(user.Dependents) != 0 {
			t.Errorf("Expected user to be initialised with 0 dependents but had %d", len(user.Dependents))
		}

		user.AddDependent(dependent)

		if len(user.Dependents) != 1 {
			t.Errorf("Expected user to have 1 dependent but had %d", len(user.Dependents))
		}

		addedDependent := user.Dependents[0]
		if addedDependent.Forename != dependent.Forename || addedDependent.Surname != dependent.Surname {
			t.Errorf("Expected added dependent to be %#v but was %#v", dependent, addedDependent)
		}
	})
}
