package stepper

import "testing"

// TestGenerics test generics.go functions
func TestGenerics(t *testing.T) {
	t.Log("stepper/generics.go tests")
	t.Run("Test reverseSequence", testReverse())
}

// testReverse test the function to reverse the switching sequence
func testReverse() func(t *testing.T) {
	return func(t *testing.T) {
		slice := Sequence28BYJ48
		slice = reverseSequence(slice)
		if slice[0][3] != 1 {
			t.Errorf(" slice[0][3] expected 1, got %v", slice[0][3])
		}
		if slice[1][1] != 0 {
			t.Errorf("slice[1][1] expected 0, got %v", slice[1][1])
		}
		if slice[2][0] != 1 {
			t.Errorf("slice[2][0] expected 1, got %v", slice[2][0])
		}
		if slice[3][2] != 0 {
			t.Errorf("slice[3][2] expected 0, got %v", slice[3][2])
		}
		if slice[4][1] != 1 {
			t.Errorf("slice[4][1] expected 1, got %v", slice[4][1])
		}
		if slice[5][3] != 0 {
			t.Errorf("slice[5][3] expected 0, got %v", slice[5][3])
		}
		if slice[6][2] != 1 {
			t.Errorf("slice[6][2] expected 1, got %v", slice[6][2])
		}
		if slice[7][0] != 0 {
			t.Errorf("slice[7][0] expected 0, got %v", slice[7][0])
		}
	}
}
