package test

import (
	"testing"
)

func TestAlwaysFailing(t *testing.T) {
	// Perform some test setup

	// Run the test
	actual := 1
	expected := 2

	if actual != expected {
		t.Errorf("Test failed. Expected: %d, but got: %d", expected, actual)
	}
}
