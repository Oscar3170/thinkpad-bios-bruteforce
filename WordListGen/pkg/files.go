package pkg

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func ReadPasswords(file io.Reader) ([][]string, error) {
	scanner := bufio.NewScanner(file)

	var passwords [][]string
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.Trim(text, " ")
		if len(text) == 0 || text[0] == '#' {
			continue
		}
		pwd := strings.Split(text, " ")
		passwords = append(passwords, pwd)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return passwords, nil
}

func DumpPasswords(file io.Writer, passwords []string) error {
	for _, p := range passwords {
		_, err := fmt.Fprintln(file, p)
		if err != nil {
			return err
		}
	}

	return nil
}
