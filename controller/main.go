package main

import (
	"GoCrackMyBIOS/pkg/keyboard"
	lt "GoCrackMyBIOS/pkg/laptop"
	"GoCrackMyBIOS/pkg/prompt"
	"bufio"
	"errors"
	"flag"
	"os"
	"strings"
	"time"

	go2c "github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
	"github.com/stianeikeland/go-rpio/v4"
)

var (
	log    = logger.NewPackageLogger("main", logger.DebugLevel)
	laptop *lt.Laptop
	i2c    *go2c.I2C

	interactive   bool
	passwordFile  string
	attemptedFile string

	attemptedPasswords []string
)

func main() {
	flag.StringVar(&passwordFile, "file", "passwords.txt", "file with 1 password per line")
	flag.StringVar(&attemptedFile, "attempted-file", "", "path to store attempted passwords")
	flag.BoolVar(&interactive, "interactive", false, "interactive mode")
	flag.Parse()

	err := setup()
	if err != nil {
		log.Fatal(err)
	}
	defer i2c.Close()

	if interactive {
		prompt.Setup(laptop, i2c).Run()
	} else {
		err := run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func setup() error {
	var err error
	if err = rpio.Open(); err != nil {
		return err
	}

	_ = logger.ChangePackageLogLevel("i2c", logger.InfoLevel)
	i2c, err = go2c.NewI2C(0x8, 1)
	if err != nil {
		return err
	}

	laptop, err = lt.NewLaptop(22, 27)
	if err != nil {
		return err
	}

	return nil
}

func run() error {
	if attemptedFile != "" {
		err := readAttemptedPasswords(attemptedFile)
		if err != nil {
			return err
		}
		log.Infof("Passwords already attempted: %d", len(attemptedPasswords))
	}

	passwords, err := readPasswords(passwordFile)
	if err != nil {
		return err
	}
	log.Infof("Passwords remaining: %d", len(passwords))

	// Padding
	if len(passwords)%3 != 0 {
		for i := 0; len(passwords)%3 != 0; i++ {
			passwords = append(passwords, "")
		}
	}

	// Split passwords in groups of 3
	var passwordGroups [][3]string
	for i := 0; i < len(passwords); i += 3 {
		passwordGroups = append(passwordGroups, [3]string(passwords[i:i+3]))
	}

	var done bool
	for i, pwdGroup := range passwordGroups {
		log.Infof("Writing passwords: %d %d %d", (i*3)+1, (i*3)+2, (i*3)+3)
		done, err = runSession(i2c, laptop, pwdGroup)
		if err != nil {
			return err
		}
		if done {
			break
		}
	}

	if !done {
		return errors.New("no passwords worked")
	}

	log.Info("Laptop Unlocked")
	return nil
}

func runSession(i2c *go2c.I2C, laptop *lt.Laptop, passwords [3]string) (done bool, err error) {
	// Make sure it's off, so we don't start in the middle of an attempt
	if laptop.GetState() == lt.ON {
		if err := laptop.FlipPowerWithRetry(3); err != nil {
			return done, err
		}
		time.Sleep(1 * time.Second)
	}
	// Turn on
	if err := laptop.FlipPowerWithRetry(3); err != nil {
		return done, err
	}
	time.Sleep(11500 * time.Millisecond)

	if err := keyboard.Write(i2c, keyboard.KeyEsc); err != nil {
		return done, err
	}
	time.Sleep(400 * time.Millisecond)

	for _, pwd := range passwords {
		log.Infof("Sending password: %s", pwd)
		if err = enterPassword(i2c, pwd); err != nil {
			log.Warningf("Failed to send password %s", pwd)
			return done, err
		}
		recordAttemptedPassword(attemptedFile, pwd)
	}

	if err := laptop.WaitForState(lt.OFF, 30*time.Second); err != nil {
		done = true
	}

	return done, nil
}

func enterPassword(i2c *go2c.I2C, pwd string) error {
	const sleepPeriod = 300 * time.Millisecond

	if err := keyboard.Write(i2c, pwd+"\n"); err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)
	if err := keyboard.Write(i2c, "\n"); err != nil {
		return err
	}
	time.Sleep(sleepPeriod)
	return nil
}

func readPasswords(path string) ([]string, error) {
	file, err := os.Open(passwordFile)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)

	var passwords []string
	for scanner.Scan() {
		text := scanner.Text()
		alreadyAttempted := false
		for _, pwd := range attemptedPasswords {
			if text == pwd {
				alreadyAttempted = true
				break
			}
		}
		if !alreadyAttempted {
			passwords = append(passwords, text)
		}
	}

	return passwords, nil
}

func readAttemptedPasswords(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Debugf("Attempted passwords file %q doesn't exist", path)
			return nil
		} else {
			return err
		}
	}
	passwords := string(content)
	passwords = strings.TrimSpace(passwords)

	attemptedPasswords = strings.Split(passwords, "\n")

	return nil
}

func recordAttemptedPassword(path string, pwd string) error {
	if attemptedFile == "" {
		return nil
	}

	attemptedPasswords = append(attemptedPasswords, pwd)

	err := os.WriteFile(path, []byte(strings.Join(attemptedPasswords, "\n")), 0600)
	if err != nil {
		return err
	}

	return nil
}
