package laptop

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/d2r2/go-logger"
	"github.com/labstack/gommon/log"
	"github.com/stianeikeland/go-rpio/v4"
)

var lg = logger.NewPackageLogger("laptop", logger.DebugLevel)

type PowerState uint8

const (
	OFF PowerState = iota
	ON
)

func (s PowerState) String() string {
	if s == 0 {
		return "OFF"
	}
	return "ON"
}

func NewPowerState(state string) (PowerState, error) {
	switch strings.ToUpper(state) {
	case ON.String():
		return ON, nil
	case OFF.String():
		return OFF, nil
	default:
		return OFF, errors.New("invalid PowerState, must be either \"ON\" or \"OFF\"")
	}
}

type Laptop struct {
	powerLed    rpio.Pin
	powerButton rpio.Pin
}

func NewLaptop(powerButtonPin uint8, powerLedPin uint8) (*Laptop, error) {
	err := rpio.Open()
	if err != nil {
		return nil, err
	}
	laptop := &Laptop{
		powerButton: rpio.Pin(powerButtonPin),
		powerLed:    rpio.Pin(powerLedPin),
	}
	laptop.powerButton.Output()
	laptop.powerLed.Input()
	laptop.powerLed.Low()
	return laptop, nil
}

func (l *Laptop) FlipPower() error {
	desiredState := 1 - PowerState(l.powerLed.Read())
	lg.Infof("Turning laptop %s", desiredState)

	l.powerButton.Low()
	l.powerButton.High()
	time.Sleep(200 * time.Millisecond)
	l.powerButton.Low()

	err := l.WaitForState(desiredState, 2*time.Second)
	if err != nil {
		return err
	}

	return nil
}

func (l *Laptop) FlipPowerWithRetry(retries int) error {
	var err error
	for i := 0; i <= retries; i++ {
		err = l.FlipPower()
		if err == nil {
			return nil
		}
		log.Infof("Retrying after failing to flip laptop power. Retries Left=%d Error=%s", retries-i, err)
	}
	return fmt.Errorf("failed to flip power after %d tries: %w", retries+1, err)
}

func (l *Laptop) WaitForState(desiredState PowerState, timeout time.Duration) error {
	lg.Infof("Waiting for laptop to turn %s", desiredState)

	var timedOut bool
	ch := make(chan bool, 1)
	go func() {
		defer close(ch)
		var currentState PowerState
		for {
			if timedOut {
				return
			}
			currentState = l.GetState()
			lg.Debugf("Laptop is currently %s", currentState)
			if currentState == desiredState {
				lg.Debugf("Laptop state reached")
				ch <- true
				return
			}
			time.Sleep(200 * time.Millisecond)
		}
	}()

	select {
	case <-ch:
		return nil
	case <-time.After(timeout):
		timedOut = true
		return errors.New(fmt.Sprintf("Timed out waiting for laptop to turn %s", desiredState))
	}
}

func (l *Laptop) GetState() PowerState {
	return PowerState(l.powerLed.Read())
}
