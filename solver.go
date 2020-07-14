package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"strconv"
)

func importWords(filename string) [][]string {
	words := make([][]string, 30)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		words[len(line)-1] = append(words[len(line)-1], line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return words
}

// Filter down to all the words where just the specified letters match
func lettersCorrect(inp string, words [][]string) []string {
	var matched []string

	reg, _ := regexp.Compile("[^A-Za-z -]")
	inpCharsOnly := reg.ReplaceAllString(inp, " ")

	for _, word := range words[len(inp)-1] {
		match := true
		for i, letter := range inpCharsOnly {
			if inpCharsOnly[i] != word[i] && letter != ' ' {
				match = false
				break
			}
		}
		if match {
			matched = append(matched, word)
		}
	}

	return matched
}

// Filter down to words where wildcards positions correspond to the same
// letters
func wildcardsUnique(inp string, words []string) []string {
	var matched []string
	var wildcards []int
	letters := make(map[byte]bool)

	for i, letter := range inp {
		if letter == '.' {
			wildcards = append(wildcards, i)
		} else {
			letters[byte(letter)] = true
		}
	}

	for _, word := range words {
		wordLetters := make(map[byte]bool)
		for k, v := range letters {
			wordLetters[k] = v
		}

		match := true
		for _, wildcard := range wildcards {
			if wordLetters[word[wildcard]] {
				match = false
				break
			}
			wordLetters[word[wildcard]] = true
		}
		if match {
			matched = append(matched, word)
		}
	}

	return matched
}

// Filter down to words where wildcard letters do not overlap with specified or "used" letters
func numberPatterns(inp string, words []string) []string {
	var matched []string
	var inpNumsOnly = inp

	for _, word := range words {
		letterNumber := make(map[byte]int)
		numberLetter := make(map[int]byte)

		match := true
		for i, char := range inpNumsOnly {
			if char >= '1' && char <= '9' {
				var num = int(char)
				if val, ok := numberLetter[num]; ok && val != word[i] {
					match = false
					break
				}
				if val, ok := letterNumber[word[i]]; ok && val != num {
					match = false
					break
				}
				letterNumber[word[i]] = num
				numberLetter[num] = word[i]
			} else {
				if val, ok := letterNumber[word[i]]; ok && val != -1 {
					match = false
					break
				}
				letterNumber[word[i]] = -1
			}
		}
		if match {
			matched = append(matched, word)
		}
	}

	return matched
}

func singlePatternSearch(words [][]string, inp string) string {
    var sb strings.Builder

	matchedLetters := lettersCorrect(inp, words)
	if len(matchedLetters) == 0 {
        sb.WriteString("> No results found.\n")
        return sb.String()
	}

    uniqueWildcardsWords := wildcardsUnique(inp, matchedLetters)
    if len(uniqueWildcardsWords) == 0 {
        sb.WriteString("> No exact results found, showing closest matches:\n")
        for _, word := range matchedLetters {
            sb.WriteString(word + "\n")
        }
        return sb.String()
    } 

    matchedNumsWords := numberPatterns(inp, uniqueWildcardsWords)
    if len(matchedNumsWords) == 0 {
        sb.WriteString("> No exact results found, printing closest matches:\n")
        for _, word := range uniqueWildcardsWords {
            sb.WriteString(word + "\n")
        }
        return sb.String()
    } 

    for _, word := range matchedNumsWords {
        sb.WriteString(word + "\n")
    }
    return sb.String()
}

// Search for two overlapping words in the codeword
// Specify the zero-indexed positions where they overlap
//
// This function is truly horrific and needs a cleanup
func doublePatternSearch(words [][]string, pattern1 string, pattern2 string, position1 int, position2 int) string {
    var sb strings.Builder

    results1 := singlePatternSearch(words, pattern1)
    if strings.Contains(results1, ">") {
        sb.WriteString("> No results found.\n")
        return sb.String()
    }
    results2 := singlePatternSearch(words, pattern2)
    if strings.Contains(results2, ">") {
        sb.WriteString("> No results found.\n")
        return sb.String()
    }

    // TODO: refactor to use arrays instead
    splat1 := strings.Split(results1, "\n")
    splat2 := strings.Split(results2, "\n")

    for _, word1 := range splat1 {
        if position1 >= len(word1) {
            continue
        }

        // Match the pattern back to the letters
        // So we can check word2 uses the same letters for numbers
        letterNumber := make(map[byte]int)
        numberLetter := make(map[int]byte)
        usedLetters := make(map[byte]int)
        for i, char := range pattern1 {
			if char >= '1' && char <= '9' {
				var num = int(char)
				letterNumber[word1[i]] = num
				numberLetter[num] = word1[i]
			} else {
				letterNumber[word1[i]] = -1
                if char != '.' {
                    usedLetters[word1[i]] = 1
                }
			}
        }

        for _, word2 := range splat2 {
            if position2 >= len(word2) {
                continue
            }

            // Didn't intersect
            if word1[position1] != word2[position2] {
                continue
            }

            // Ensure patterns are consistent with first word
            var match = true
            for i, char := range pattern2 {
                if char >= '1' && char <= '9' {
                    var num = int(char)
                    if val, ok := numberLetter[num]; ok && val != word2[i] {
                    	match = false
                    	break
                    }
                    if val, ok := letterNumber[word2[i]]; ok && val != num {
                    	match = false
                    	break
                    }
                } else {
                    if val, ok := letterNumber[word2[i]]; ok && val != -1 {
                    	match = false
                    	break
                    }

                    // If wildcard then it can't be a letter used directly in word1
                    if char == '.' {
                        if val, ok := usedLetters[word2[i]]; ok && val != -1 {
                            match = false
                            break
                        }
                    }
                }
            }

            if match == false {
                continue
            }

            sb.WriteString(word1 + ":" + word2 + "\n")
        }
    }

    if sb.Len() == 0 {
        sb.WriteString("> No results found.\n")
    }

    return sb.String()
}

// Canonicalise special chars to dots for wildcard input positions
func canonicaliseInput(inp string) string {
	inp = strings.TrimSpace(strings.ToLower(inp))

	reg, _ := regexp.Compile("[ *?]")
	return reg.ReplaceAllString(inp, ".")
}

func singlePatternWrapper(words [][]string, pattern string) {
	pattern = canonicaliseInput(pattern)
    output := singlePatternSearch(words, pattern)
    fmt.Print(output)
}

func doublePatternWrapper(words [][]string, pattern1 string, pattern2 string, pos1 string, pos2 string) {
	pattern1 = canonicaliseInput(pattern1)
	pattern2 = canonicaliseInput(pattern2)

    position1, err := strconv.Atoi(pos1)
    if err != nil {
        fmt.Println("Invalid command line arguments specified, please read the README")
        return
    }
    position2, err := strconv.Atoi(pos2)
    if err != nil {
        fmt.Println("Invalid command line arguments specified, please read the README")
        return
    }
    position1 -= 1
    position2 -= 1

    if position1 < 0 || position1 >= len(pattern1) {
        fmt.Println("Intersection position for word 1 is outside the word")
        return
    }
    if position2 < 0 || position2 >= len(pattern2) {
        fmt.Println("Intersection position for word 2 is outside the word")
        return
    }

    output := doublePatternSearch(words, pattern1, pattern2, position1, position2)
    fmt.Print(output)
}

func main() {
	dictPtr := flag.String("dict", "dictionaries/merged.txt", "Dictionary file location.")
	flag.Parse()

	var words = importWords(*dictPtr)

	if len(flag.Args()) == 0 {
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter pattern: ")
			text, _ := reader.ReadString('\n')
			if len(text) > 1 {
				text = text[:len(text)-1]
                splat := strings.Fields(text)
                if len(splat) == 1 {
                    singlePatternWrapper(words, text)
                } else if len(splat) == 2 {
                    doublePatternWrapper(words, splat[0], splat[1], "1", "1")
                } else if len(splat) == 4 {
                    doublePatternWrapper(words, splat[0], splat[1], splat[2], splat[3])
                } else {
                    fmt.Println("Invalid number of command line arguments specified, please read the README")
                }
			}
		}
	} else if len(flag.Args()) == 1 {
        singlePatternWrapper(words, flag.Arg(0))
	} else if len(flag.Args()) == 2 {
		doublePatternWrapper(words, flag.Arg(0), flag.Arg(1), "1", "1")
    } else if len(flag.Args()) == 4 {
		doublePatternWrapper(words, flag.Arg(0), flag.Arg(1), flag.Arg(2), flag.Arg(3))
    } else {
        panic("Invalid number of command line arguments specified, please read the README")
    }
}
