package jlog

import "testing"

func TestCurrentPackageFileName(t *testing.T) {
	expected := "github.com/jerejones/jlog/namers_test.go"
	actual := CurrentPackageFileName()
	if expected != actual {
		t.Errorf("Expected: %s\nGot: %s", expected, actual)
	}
}
