package pkg_test

import (
	"WordListGen/pkg"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPwdCharsCombinations(t *testing.T) {
	config := &pkg.Config{
		CharReplacements: map[string][]string{
			"o": {"0"},
			"a": {"4", "@"},
		},
	}
	passwords := [][]string{
		{"foo", "bar"},
	}
	expected := [][]string{
		{"foo", "bar"},
		{"f0o", "bar"},
		{"fo0", "bar"},
		{"f00", "bar"},
		{"foo", "b4r"},
		{"f0o", "b4r"},
		{"fo0", "b4r"},
		{"f00", "b4r"},
		{"foo", "b@r"},
		{"f0o", "b@r"},
		{"fo0", "b@r"},
		{"f00", "b@r"},
	}

	result := pkg.PwdCharsCombinations(config, passwords)

	// for _, r := range result {
	// 	fmt.Println(r)
	// }

	assert.Equal(t, expected, result)

}

func TestSeparatorsCombinations(t *testing.T) {
	config := &pkg.Config{
		Separators: []string{" ", "-"},
	}
	passwords := [][]string{
		{"foo", "bar", "test"},
		{"fizz", "buzz", "fizzbuzz"},
	}
	expected := []string{
		"foo bar test", "fizz buzz fizzbuzz",
		"foo-bar-test", "fizz-buzz-fizzbuzz",
	}

	result := pkg.SeparatorsCombinations(config, passwords)

	assert.Equal(t, expected, result)
}

func TestEndingsCombinations(t *testing.T) {
	config := &pkg.Config{
		Endings: []string{"?", "!", ""},
	}
	passwords := []string{
		"foo-bar-test",
		"fizz-buzz-fizzbuzz",
	}
	expected := []string{
		"foo-bar-test?", "fizz-buzz-fizzbuzz?",
		"foo-bar-test!", "fizz-buzz-fizzbuzz!",
		"foo-bar-test", "fizz-buzz-fizzbuzz",
	}

	result := pkg.EndingsCombinations(config, passwords)

	assert.Equal(t, expected, result)
}

func TestReplaceWordsCombinations(t *testing.T) {
	config := &pkg.Config{
		WordReplacements: map[string][]string{
			"foo":  {"faa", "fii"},
			"test": {"teste"},
			"fizz": {"fiz", "faz"},
		},
	}
	passwords := [][]string{
		{"foo", "bar", "test"},
		{"fizz", "buzz", "fizzbuzz"},
	}
	expected := [][]string{
		{"foo", "bar", "test"},
		{"faa", "bar", "test"},
		{"fii", "bar", "test"},
		{"foo", "bar", "teste"},
		{"faa", "bar", "teste"},
		{"fii", "bar", "teste"},
		{"fizz", "buzz", "fizzbuzz"},
		{"fiz", "buzz", "fizzbuzz"},
		{"faz", "buzz", "fizzbuzz"},
	}

	result := pkg.ReplaceWordsCombinations(config, passwords)

	// for _, r := range result {
	// 	fmt.Println(r)
	// }

	assert.Equal(t, expected, result)

}
