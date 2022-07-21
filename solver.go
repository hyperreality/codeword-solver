package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
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
		if len(line) <= 30 {
			words[len(line)-1] = append(words[len(line)-1], line)
		}
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

func singlePatternSearch(words [][]string, inp string) []string {
	output := []string{}

	matchedLetters := lettersCorrect(inp, words)
	if len(matchedLetters) == 0 {
		output = append(output, "> No exact results found, printing closest matches:\n")
		return output
	}

	uniqueWildcardsWords := wildcardsUnique(inp, matchedLetters)
	if len(uniqueWildcardsWords) == 0 {
		output = append(output, "> No exact results found, printing closest matches:\n")
		for _, word := range matchedLetters {
			output = append(output, word)
		}
		return output
	}

	matchedNumsWords := numberPatterns(inp, uniqueWildcardsWords)
	if len(matchedNumsWords) == 0 {
		output = append(output, "> No exact results found, printing closest matches:\n")
		for _, word := range uniqueWildcardsWords {
			output = append(output, word)
		}
		return output
	}

	for _, word := range matchedNumsWords {
		output = append(output, word)
	}
	return output
}

// Search for two overlapping words in the codeword
// Specify the zero-indexed positions where they overlap
//
// This function is truly horrific and needs a cleanup
func doublePatternSearch(words [][]string, pattern1 string, pattern2 string) []string {
	output := []string{}

	results1 := singlePatternSearch(words, pattern1)
	if strings.Contains(results1[0], ">") {
		return []string{}
	}
	results2 := singlePatternSearch(words, pattern2)
	if strings.Contains(results2[0], ">") {
		return []string{}
	}

	for _, word1 := range results1 {
		if len(word1) == 0 {
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

		for _, word2 := range results2 {
			if len(word2) == 0 {
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
						// or a letter used in a number in first word
						if _, ok := letterNumber[word2[i]]; ok {
							match = false
							break
						}
					}
				}
			}

			if match == false {
				continue
			}

			correct_words := fmt.Sprintf("%s:%s", word1, word2)
			output = append(output, correct_words)
		}
	}

	return output
}

// Canonicalise special chars to dots for wildcard input positions
func canonicaliseInput(inp string) string {
	inp = strings.TrimSpace(strings.ToLower(inp))

	reg, _ := regexp.Compile("[ *?]")
	return reg.ReplaceAllString(inp, ".")
}

func singlePatternWrapper(words [][]string, pattern string) []string {
	pattern = canonicaliseInput(pattern)
	return singlePatternSearch(words, pattern)
}

func doublePatternWrapper(words [][]string, pattern1 string, pattern2 string) []string {
	pattern1 = canonicaliseInput(pattern1)
	pattern2 = canonicaliseInput(pattern2)
	output1 := doublePatternSearch(words, pattern1, pattern2)
	output2 := doublePatternSearch(words, pattern2, pattern1)

	output2_flipped := []string{}
	for _, word := range output2 {
		splat := strings.Split(word, ":")
		correct_words := fmt.Sprintf("%s:%s", splat[1], splat[0])
		output2_flipped = append(output2_flipped, correct_words)
	}

	return intersection(output1, output2_flipped)
}

func printResults(results []string) {
	var sb strings.Builder
	if len(results) == 0 {
		sb.WriteString("> No results found.\n")
	} else {
		for _, words := range results {
			sb.WriteString(words + "\n")
		}
	}
	fmt.Print(sb.String())
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
					printResults(singlePatternWrapper(words, text))
				} else if len(splat) == 2 {
					printResults(doublePatternWrapper(words, splat[0], splat[1]))
				} else {
					fmt.Println("Invalid number of command line arguments specified, please read the README")
				}
			}
		}
	} else if len(flag.Args()) == 1 {
		printResults(singlePatternWrapper(words, flag.Arg(0)))
	} else if len(flag.Args()) == 2 {
		printResults(doublePatternWrapper(words, flag.Arg(0), flag.Arg(1)))
	} else {
		panic("Invalid number of command line arguments specified, please read the README")
	}
}
