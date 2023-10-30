package pkg

import (
	"strings"
)

func PwdCharsCombinations(config *Config, passwords [][]string) [][]string {
	var result [][]string

	for _, pwd := range passwords {
		// Create one slice per original password
		var pwdCombinations [][]string
		// Add original password to slice
		pwdCombinations = append(pwdCombinations, make([]string, len(pwd)))
		copy(pwdCombinations[0], pwd)

		for iword, word := range pwd {
			iters := len(pwdCombinations)
			// fmt.Printf("Word: %s\n", word)
			// fmt.Printf("Passwords length: %d\n", iters)

			combinations := WordCharsCombinations(config, word)
			// fmt.Printf("Combinations %v\n", combinations)
			for _, c := range combinations[1:] {
				for i := 0; i < iters; i++ {
					newPwd := make([]string, len(pwdCombinations[i]))
					copy(newPwd, pwdCombinations[i])
					newPwd[iword] = c
					pwdCombinations = append(pwdCombinations, newPwd)
					// fmt.Printf("%v\n", generatedPwds)
				}
			}
		}
		result = append(result, pwdCombinations...)
	}
	return result
}

// Get all the characters combinations using the config character replacements
func WordCharsCombinations(config *Config, word string) []string {
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
		// fmt.Printf("Character: %c\n", char)
		// fmt.Printf("Words length: %d\n", iters)
		for i := 0; i < iters; i++ {
			words[i][ichar] = char
		}

		replacements, ok := config.CharReplacements[string(char)]
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
					newPerm = append(newPerm[:ichar], append([]rune(r), newPerm[ichar+1:]...)...)
				}
				words = append(words, newPerm)
				// fmt.Printf("%v\n", words)
			}
		}
	}
	// fmt.Printf("Final words length: %d\n", len(words))

	joinedWords := make([]string, len(words))
	for iw, w := range words {
		joinedWords[iw] = string(w)
	}

	return joinedWords

}

func ReplaceWordsCombinations(config *Config, passwords [][]string) [][]string {
	var result [][]string

	for _, pwd := range passwords {
		// Create one slice per original password
		var pwdCombinations [][]string
		// Add original password to slice
		pwdCombinations = append(pwdCombinations, make([]string, len(pwd)))
		copy(pwdCombinations[0], pwd)

		for iword, word := range pwd {
			iters := len(pwdCombinations)
			// fmt.Printf("Word: %s\n", word)
			// fmt.Printf("Passwords length: %d\n", iters)
			for i := 0; i < iters; i++ {
				pwdCombinations[i][iword] = word
			}

			replacements, ok := config.WordReplacements[word]
			if !ok {
				continue
			}

			for _, r := range replacements {
				for i := 0; i < iters; i++ {
					newPermutation := make([]string, len(pwdCombinations[i]))
					copy(newPermutation, pwdCombinations[i])
					newPermutation[iword] = r
					pwdCombinations = append(pwdCombinations, newPermutation)
				}
			}
		}
		result = append(result, pwdCombinations...)
	}
	return result
}

// Generate all passwords with all the configured separators
func SeparatorsCombinations(config *Config, passwords [][]string) []string {
	var joinedPwds []string

	for _, separator := range config.Separators {
		for _, pwd := range passwords {
			joinedPwds = append(joinedPwds, strings.Join(pwd, separator))
		}
	}

	return joinedPwds
}

func EndingsCombinations(config *Config, passwords []string) []string {
	var results []string

	for _, ending := range config.Endings {
		for _, word := range passwords {
			result := word + ending
			results = append(results, result)
		}
	}

	return results
}

// TODO: replaceWords
// TODO: calculate the time to test all the passwords
// TODO: config from file
