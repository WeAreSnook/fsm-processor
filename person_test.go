package main

import "testing"

func TestIsSameAs(t *testing.T) {
	t.Run("returns true for users with the same name", func(t *testing.T) {
		userA := Person{forename: "bob", surname: "smith"}
		userB := Person{forename: "bob", surname: "smith"}

		if !userA.IsSameAs(userB) {
			t.Errorf("Expected %#v to be same as %#v but got false", userA, userB)
		}
	})

	t.Run("returns true for users with the same name", func(t *testing.T) {
		userA := Person{forename: "bob", surname: "smith"}
		userB := Person{forename: "jane", surname: "smith"}

		if userA.IsSameAs(userB) {
			t.Errorf("Expected %#v NOT to be same as %#v but got false", userA, userB)
		}
	})
}
