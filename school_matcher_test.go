package main

import "testing"

func TestCompareStrings(t *testing.T) {
	t.Run("identical names = 1", func(t *testing.T) {
		nameA := "chris"
		nameB := "chris"

		score := CompareStrings(nameA, nameB)

		if score < 1 {
			t.Fatalf("Expected score 1 but got %f", score)
		}
	})

	t.Run("similar names have a high score", func(t *testing.T) {
		nameA := "chris"
		nameB := "chriss"

		score := CompareStrings(nameA, nameB)

		if score < 0.9 {
			t.Fatalf("Expected score > 0.9 but got %f", score)
		}
	})

	t.Run("different names have a low score", func(t *testing.T) {
		nameA := "chris"
		nameB := "bob"

		score := CompareStrings(nameA, nameB)

		if score > 0.1 {
			t.Fatalf("Expected score < 0.1 but got %f", score)
		}
	})
}
