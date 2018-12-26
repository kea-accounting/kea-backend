package util

import (
	"testing"
)

func assert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assert")
	}
}

func TestValidName(t *testing.T) {
	assert(t, IsValidName("abc"))
	assert(t, IsValidName("a_bc-Z"))
	assert(t, !IsValidName("+abc"))
	assert(t, !IsValidName(" "))
}

func TestID(t *testing.T) {

	for i := 0; i < 1000; i++ {
		id := NewID()
		assert(t, IsValidName(id))
	}
}
