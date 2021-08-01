package restcountries

import (
	"testing"
)

func TestLCFirst(t *testing.T) {
	got := lCFirst("ChrisCross")
	want := "chrisCross"

	if got != want {
		t.Errorf("got %s; wanted %s", got, want)
	}
}

func TestLCFirstEmpty(t *testing.T) {
	got := lCFirst("")
	want := ""

	if got != want {
		t.Errorf("got %s; wanted %s", got, want)
	}
}
