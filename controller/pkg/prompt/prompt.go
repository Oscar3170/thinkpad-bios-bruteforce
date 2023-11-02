package prompt

import (
	"GoCrackMyBIOS/pkg/keyboard"
	"GoCrackMyBIOS/pkg/laptop"
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	go2c "github.com/d2r2/go-i2c"
)

type LaptopPrompt struct {
	prompt *prompt.Prompt
	lapt   *laptop.Laptop
	i2c    *go2c.I2C
}

func Setup(lapt *laptop.Laptop, i2c *go2c.I2C) *LaptopPrompt {
	laptopPrompt := &LaptopPrompt{
		lapt: lapt,
		i2c:  i2c,
	}
	pt := prompt.New(laptopPrompt.executor, completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("laptop-prompt"),
	)
	laptopPrompt.prompt = pt
	return laptopPrompt
}

func (p *LaptopPrompt) Run() {
	p.prompt.Run()
}

func (p *LaptopPrompt) executor(text string) {
	text = strings.TrimSpace(text)

	fields := strings.Fields(text)
	if len(fields) > 0 {
		switch fields[0] {
		case "!help":
			fmt.Println("Available commands: help, exit, turn, tab, enter")
		case "!exit":
			os.Exit(0)
		case "!turn":
			// Handle the turn command with arguments
			if len(fields) < 1 {
				fmt.Println("turn command requires arguments")
				return
			}
			desiredState, err := laptop.NewPowerState(fields[1])
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
			if p.lapt.GetState() != desiredState {
				if err := p.lapt.FlipPower(); err != nil {
					fmt.Printf("Error: %s\n", err)
				}
			}
			return
		case "!tab":
			keyboard.Write(p.i2c, "\t")
			if err := keyboard.Write(p.i2c, "\t"); err != nil {
				fmt.Printf("Error: %s\n", err)
			}
			return
		case "!enter":
			keyboard.Write(p.i2c, "\n")
			if err := keyboard.Write(p.i2c, "\n"); err != nil {
				fmt.Printf("Error: %s\n", err)
			}
			return
		case "!esc":
			if err := keyboard.Write(p.i2c, keyboard.KeyEsc); err != nil {
				fmt.Printf("Error: %s\n", err)
			}
			return
		case "!space":
			if err := keyboard.Write(p.i2c, " "); err != nil {
				fmt.Printf("Error: %s\n", err)
			}
			return
		default:
			if err := keyboard.Write(p.i2c, text+"\n"); err != nil {
				fmt.Printf("Error: %s\n", err)
			}
			return
		}
	}
}

func completer(doc prompt.Document) []prompt.Suggest {
	word := doc.GetWordBeforeCursor()
	if word == "" {
		return []prompt.Suggest{}
	}

	s := []prompt.Suggest{
		{Text: "!help", Description: "Describes the other commands and how to interact with the prompt"},
		{Text: "!exit", Description: "Terminates the program"},
		{Text: "!turn", Description: "Has arguments which it sends to a custom function"},
		{Text: "!tab", Description: "Send a tab character"},
		{Text: "!enter", Description: "Executes the function write with the argument \n"},
	}
	return prompt.FilterHasPrefix(s, word, true)
}
