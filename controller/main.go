package main

import (
	"GoCrackMyBIOS/pkg/keyboard"
	lt "GoCrackMyBIOS/pkg/laptop"
	"errors"
	"flag"
	go2c "github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
	"github.com/stianeikeland/go-rpio/v4"
	"os"
	"strings"
	"time"
)

var log = logger.NewPackageLogger("main", logger.DebugLevel)

func main() {
	_ = logger.ChangePackageLogLevel("i2c", logger.InfoLevel)

	passwordFile := flag.String("file", "passwords.txt", "file with 1 password per line")

	content, err := os.ReadFile(*passwordFile)
	if err != nil {
		log.Fatal(err)
	}
	passwords := strings.Split(string(content), "\n")

	if err = run(passwords); err != nil {
		log.Fatal(err)
	}
}

func run(passwords []string) error {
	if err := rpio.Open(); err != nil {
		return err
	}
	i2c, err := go2c.NewI2C(0x8, 1)
	if err != nil {
		return err
	}
	defer i2c.Close()
	lapt, err := lt.NewLaptop(27, 4)
	if err != nil {
		return err
	}

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
		log.Infof("Writing passwords: %d %d %d", i+1, i+2, i+3)
		done, err = runSession(i2c, lapt, pwdGroup)
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
		if err := laptop.FlipPower(); err != nil {
			return done, err
		}
	}
	time.Sleep(4 * time.Second)
	// Turn on
	if err := laptop.FlipPower(); err != nil {
		return done, err
	}
	time.Sleep(12 * time.Second)

	if err := keyboard.Write(i2c, keyboard.KeyEsc); err != nil {
		return done, err
	}
	time.Sleep(1 * time.Second)

	for _, pwd := range passwords {
		log.Infof("Sending password: %s", pwd)
		if err = enterPassword(i2c, pwd); err != nil {
			log.Warningf("Failed to send password %s", pwd)
			return done, err
		}
	}

	if err := laptop.WaitForState(lt.OFF, 60*time.Second); err != nil {
		done = true
	}

	return done, nil
}

func enterPassword(i2c *go2c.I2C, pwd string) error {
	if err := keyboard.Write(i2c, pwd+"\n"); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)
	if err := keyboard.Write(i2c, "\n"); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)
	return nil
}
