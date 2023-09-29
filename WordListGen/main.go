package main

import (
	"fmt"
	"strings"
)

type Config struct {
	WordReplacements map[string][]string
	CharReplacements map[rune][]string
	Endings          []string
	Separators       []string
}

var DefaultConfig = &Config{
	CharReplacements: map[rune][]string{
		'o': {"0"},
		'a': {"4", "@"},
		'i': {"1", "l", "1"},
	},
	Endings:    []string{"1337", "420", "69", "!"},
	Separators: []string{" ", "-", "_", ""},
}

func endingsCombinations(config *Config, words [][]string) [][]string {
	iters := len(words)

	for _, ending := range config.Endings {
		for i := 0; i < iters; i++ {
			word := append(words[i], ending)
			words = append(words, word)
		}
	}

	return words
}

// Get all the characters combinations using the config character replacements
func charsCombinations(config *Config, word string) []string {
	// words layout for input "foo":
	// [ 'f', 'o', 'o'],
	// [ 'f', '0', 'o'],
	// [ 'f', 'o', '0'],
	// [ 'f', '0', '0'],
	var words [][]rune

	wordRune := []rune(word)
	words = append(words, make([]rune, len(wordRune)))
	for ichar, char := range wordRune {
		iters := len(words)
		// fmt.Printf("Character: %s\n", char)
		// fmt.Printf("Words length: %d\n", iters)
		for i := 0; i < iters; i++ {
			words[i][ichar] = char
		}

		replacements, ok := config.CharReplacements[char]
		if !ok {
			continue
		}

		for _, r := range replacements {
			for i := 0; i < iters; i++ {
				newPerm := make([]rune, len(words[i]))
				copy(newPerm, words[i])
				if ichar+1 == len(wordRune) {
					newPerm = append(newPerm[:ichar], []rune(r)...)
				} else {
					newPerm = append(newPerm[:ichar], append([]rune(r), newPerm[ichar:]...)...)
				}
				words = append(words, newPerm)
				// fmt.Printf("%v\n", words)
			}
		}
	}
	fmt.Printf("Final words length: %d\n", len(words))

	joinedWords := make([]string, len(words))
	for iw, w := range words {
		joinedWords[iw] = string(w)
	}

	return joinedWords

}

// TODO: replaceWords
// TODO: separators
// TODO: calculate the time to test all the passwords
// TODO: write to a file by input password
// TODO: config from file

func main() {

	passwords := [][]string{
		{"foo", "bar"},
	}

	for _, pwd := range passwords {
		var generatedPwds [][]string
		generatedPwds = append(generatedPwds, make([]string, len(pwd)))
		copy(generatedPwds[0], pwd)

		for iword, word := range pwd {
			iters := len(generatedPwds)
			fmt.Printf("Word: %s\n", word)
			fmt.Printf("Passwords length: %d\n", iters)

			combinations := charsCombinations(DefaultConfig, word)
			fmt.Printf("Combinations %v\n", combinations)
			for _, c := range combinations[1:] {
				for i := 0; i < iters; i++ {
					newPwd := make([]string, len(generatedPwds[i]))
					copy(newPwd, generatedPwds[i])
					newPwd[iword] = c
					generatedPwds = append(generatedPwds, newPwd)
					// fmt.Printf("%v\n", generatedPwds)
				}
			}
		}

		generatedPwds = endingsCombinations(DefaultConfig, generatedPwds)

		for _, generatedPwd := range generatedPwds {
			fmt.Println(strings.Join(generatedPwd, "_"))

		}
	}
}
