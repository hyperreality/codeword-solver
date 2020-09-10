package main

import (
    "testing"
    "strings"
)

func TestSinglePatternSearch(t *testing.T) {
    var words = importWords("dictionaries/merged.txt")

    var pattern = "11.e"
    var results = singlePatternSearch(words, pattern)
    var expected = "ooze\n"

    if results != expected {
        t.Error("got: {}, want: {}.", results, expected)
    }
}

func TestDoublePatternSearch(t *testing.T) {
    var words = importWords("dictionaries/merged.txt")

    var pattern1 = ".osmos"
    var pattern2 = "m....o..sm"
    var results = doublePatternSearch(words, pattern1, pattern2)
    var expected = "cosmos:metabolism\n"

    if results != expected {
        t.Error("got: {}, want: {}.", results, expected)
    }
}

func TestDoublePatternIssue(t *testing.T) {
    var words = importWords("dictionaries/merged.txt")

    var pattern1 = "12.33"
    var pattern2 = "31...2"
    var results = doublePatternSearch(words, pattern1, pattern2)

    var expected = "aglee:earwig\n"
    if !strings.Contains(results, expected) {
        t.Error("got: {}, want: {}.", results, expected)
    }

    // two r's in third position should not qualify
    var not_expected = "agree:earwig\n"
    if strings.Contains(results, not_expected) {
        t.Error("got: {}, did not want: {}.", results, not_expected)
    }

}
