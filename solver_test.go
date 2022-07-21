package main

import (
	"reflect"
	"testing"
)

func TestSinglePatternSearch(t *testing.T) {
	var words = importWords("dictionaries/merged.txt")

	var pattern = "11.e"
	var results = singlePatternWrapper(words, pattern)
	expected := "ooze"

	if !contains(results, expected) {
		t.Error("got: {}, want: {}.", results, expected)
	}
}

func TestDoublePatternSearch(t *testing.T) {
	var words = importWords("dictionaries/merged.txt")

	var pattern1 = ".osmos"
	var pattern2 = "m....o..sm"
	var results = doublePatternWrapper(words, pattern1, pattern2)

	expected := []string{"cosmos:metabolism"}

	if !reflect.DeepEqual(results, expected) {
		t.Error("got: {}, want: {}.", results, expected)
	}
}

func TestDoublePatternIssue(t *testing.T) {
	var words = importWords("dictionaries/merged.txt")

	var pattern1 = "12.33"
	var pattern2 = "31...2"
	var results = doublePatternWrapper(words, pattern1, pattern2)

	expected := "aglee:earwig"
	if !contains(results, expected) {
		t.Error("got: {}, want: {}.", results, expected)
	}

	// two r's in third position should not qualify
	not_expected := "agree:earwig"
	if contains(results, not_expected) {
		t.Error("got: {}, did not want: {}.", results, not_expected)
	}

}

func TestDoublePatternIssue2(t *testing.T) {
	var words = importWords("dictionaries/merged.txt")

	var pattern1 = "12..s"
	var pattern2 = "s.e21"
	var results = doublePatternWrapper(words, pattern1, pattern2)

	// e specified in second word, so should not appear in first word wildcard space
	var not_expected = "mazes:steam"
	if contains(results, not_expected) {
		t.Error("got: {}, did not want: {}.", results, not_expected)
	}
}
