package pkg

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	WordReplacements map[string][]string
	CharReplacements map[string][]string
	Separators       []string
	Endings          []string
}

func (c *Config) String() string {
	wordR := fmt.Sprintf("Word Replacements: %q", c.WordReplacements)
	charR := fmt.Sprintf("Character Replacements: %q", c.CharReplacements)
	sep := fmt.Sprintf("Word Separators: %q", c.Separators)
	endings := fmt.Sprintf("Password Endings: %q", c.Endings)

	return strings.Join([]string{wordR, charR, sep, endings}, "\n")
}

var DefaultConfig = &Config{
	CharReplacements: map[string][]string{
		"o": {"0"},
		"a": {"4"},
	},
	Endings:    []string{"", "?"},
	Separators: []string{" "},
}

func ReadConfig(filePath string) (*Config, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var c *Config
	// err = viper.Unmarshal(&c)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed config unmarshal: %w", err)
	// }
	// viper.SetConfigType("yaml")
	// err = viper.ReadConfig(file)
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return nil, fmt.Errorf("fail reading config: %w", err)
	}

	log.Printf("Config: \n%s", c)
	// log.Println(viper.Get("endings"))

	return c, nil
}
