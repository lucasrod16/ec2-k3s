package utils

import (
	"reflect"
	"testing"
)

// TestDerefString tests that we can dereference string pointers to strings correctly
func TestDerefString(t *testing.T) {
	var input *string
	output := DerefString(input)
	got := reflect.TypeOf(output).Kind()
	expected := reflect.String

	if expected != got {
		t.Errorf("error DerefString(): expected: %s | got: %s", expected, got)
	}

}
