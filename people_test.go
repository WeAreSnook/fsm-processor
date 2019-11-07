package main

import "testing"

func TestFindExisting(t *testing.T) {
	chris := Person{forename: "Chris", surname: "Sloey", age: 29}
	michael := Person{forename: "Michael", surname: "Hayes", age: 31}

	store := PeopleStore{
		people: []Person{
			chris,
			michael,
		},
	}

	t.Run("finds the correct person", func(t *testing.T) {
		want := chris
		got, err := store.FindExisting(Person{forename: chris.forename, surname: chris.surname})

		if err != nil {
			t.Errorf("error searching for %#v", want)
		}

		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("error when person doesn't exist", func(t *testing.T) {
		unknownPerson := Person{forename: "Doesn't", surname: "Exist"}
		_, err := store.FindExisting(unknownPerson)

		if err == nil {
			t.Errorf("No error for missing person")
		}
	})
}

func TestAdd(t *testing.T) {
	chris := Person{forename: "Chris", surname: "Sloey", age: 29}
	michael := Person{forename: "Michael", surname: "Hayes", age: 31}
	newPerson := Person{forename: "Bob", surname: "WOW", age: 92}

	store := PeopleStore{
		people: []Person{
			chris,
			michael,
		},
	}

	t.Run("adds person successfully", func(t *testing.T) {
		store.Add(newPerson)
		addedUser, err := store.FindExisting(newPerson)

		if err != nil {
			t.Errorf("error searching for added user %#v", err)
		}

		if !addedUser.IsSameAs(newPerson) {
			t.Errorf("user wasn't added, expected %#v to be same as %#v", addedUser, newPerson)
		}
	})
}

func TestUpdate(t *testing.T) {
	chris := Person{forename: "Chris", surname: "Sloey", age: 29}
	michael := Person{forename: "Michael", surname: "Hayes", age: 31}

	store := PeopleStore{
		people: []Person{
			chris,
			michael,
		},
	}

	t.Run("updates correct person", func(t *testing.T) {
		updatedChrisDetails := Person{forename: "Christopher", surname: "Sloey", age: 29}
		err := store.Update(chris, updatedChrisDetails)

		if err != nil {
			t.Errorf("error updating user %#v", err)
		}

		updatedUser, err := store.FindExisting(updatedChrisDetails)

		if err != nil {
			t.Errorf("error finding updated user user %#v", err)
		}

		if !updatedUser.IsSameAs(updatedChrisDetails) {
			t.Errorf("user wasn't added, expected %#v to be same as %#v", updatedUser, updatedChrisDetails)
		}
	})
}
