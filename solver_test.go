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

func TestDoublePatternSearch(t *testing.T) {
    var words = importWords("dictionaries/merged.txt")

    var pattern1 = ".osmos"
    var pos1 = 3
    var pattern2 = "m....o..sm"
    var pos2 = 9
    var results = doublePatternSearch(words, pattern1, pattern2, pos1, pos2)
    var expected = "cosmos:metabolism\n"

    if results != expected {
        t.Error("got: {}, want: {}.", results, expected)
    }
}
