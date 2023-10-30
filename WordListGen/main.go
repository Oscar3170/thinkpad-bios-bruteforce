package main

import (
	"WordListGen/pkg"
	"errors"
	"io"
	"log"
	"os"
	"time"

	flag "github.com/spf13/pflag"
)

type cmdFlags struct {
	inputPath  string
	outputPath string
	configPath string
}

var (
	cmdArgs cmdFlags
	config  *pkg.Config
)

func setup() error {
	flag.StringVarP(&cmdArgs.outputPath, "out", "o", "", "path to file containing the base passwords. Defaults to stdin")
	flag.StringVarP(&cmdArgs.inputPath, "in", "i", "", "path to output file. Defaults to stdout")
	flag.StringVarP(&cmdArgs.configPath, "config", "c", "", "config file path")

	flag.Parse()

	if cmdArgs.configPath != "" {
		var err error
		config, err = pkg.ReadConfig(cmdArgs.configPath)
		if err != nil {
			return err
		}
	} else {
		config = pkg.DefaultConfig
	}

	if cmdArgs.inputPath == "" {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return errors.New("stdin is not pipe")
		}
	}

	return nil
}

func readPasswords() (passwords [][]string, err error) {
	var file io.Reader = os.Stdin

	if cmdArgs.inputPath != "" {
		osFile, err := os.Open(cmdArgs.inputPath)
		if err != nil {
			return nil, err
		}
		defer osFile.Close()
		file = osFile
	}

	return pkg.ReadPasswords(file)
}

func writePasswords(passwords []string) error {
	var file io.Writer = os.Stdout

	if cmdArgs.outputPath != "" {
		osFile, err := os.Create(cmdArgs.outputPath)
		if err != nil {
			return err
		}
		defer osFile.Close()

		err = os.Chmod(cmdArgs.outputPath, 0600)
		if err != nil {
			return err
		}
		file = osFile
	}

	return pkg.DumpPasswords(file, passwords)
}

func main() {
	err := setup()
	if err != nil {
		log.Fatal(err)
	}

	passwords, err := readPasswords()
	if err != nil {
		log.Fatal(err)
	}

	var generatedPasswords []string
	for _, pwd := range passwords {
		splitPwds := pkg.ReplaceWordsCombinations(config, [][]string{pwd})

		splitPwds = pkg.PwdCharsCombinations(config, splitPwds)

		pwds := pkg.SeparatorsCombinations(config, splitPwds)

		pwds = pkg.EndingsCombinations(config, pwds)

		generatedPasswords = append(generatedPasswords, pwds...)
	}

	err = writePasswords(generatedPasswords)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Generated %d passwords\n", len(generatedPasswords))
	log.Printf("Would take %v\n", time.Duration(len(generatedPasswords)*3)*time.Second)
}
