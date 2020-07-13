package main

import (
    "testing"
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
