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
		words[len(line)-1] = append(words[len(line)-1], line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return words
}

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

func cleanInput(inp string) string {
	inp = strings.TrimSpace(strings.ToLower(inp))

	reg, _ := regexp.Compile("[ *?]")
	return reg.ReplaceAllString(inp, ".")
}

func printResults(words [][]string, inp string) {
	inp = cleanInput(inp)

	matchedLetters := lettersCorrect(inp, words)
	if len(matchedLetters) == 0 {
		fmt.Println("No results found.")
	} else {
		uniqueWildcards := wildcardsUnique(inp, matchedLetters)

		if len(uniqueWildcards) == 0 {
			fmt.Println("No exact results found, printing closest matches: ")
			for _, word := range matchedLetters {
				fmt.Println(word)
			}
		} else {
			matchedNums := numberPatterns(inp, uniqueWildcards)

			if len(matchedNums) == 0 {
				fmt.Println("No exact results found, printing closest matches: ")
				for _, word := range uniqueWildcards {
					fmt.Println(word)
				}
			} else {
				for _, word := range matchedNums {
					fmt.Println(word)
				}
			}
		}
	}
}

func main() {
	dictPtr := flag.String("dict", "dict.txt", "Dictionary file location.")
	flag.Parse()

	var words = importWords(*dictPtr)

	if len(flag.Args()) == 0 {
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter pattern: ")
			text, _ := reader.ReadString('\n')
			if len(text) > 1 {
				text = text[:len(text)-1]
				printResults(words, text)
				fmt.Println()
			}
		}
	} else {
		printResults(words, flag.Arg(0))
	}
}
