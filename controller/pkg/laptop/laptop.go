package laptop

import (
	"errors"
	"fmt"
	"github.com/d2r2/go-logger"
	"github.com/stianeikeland/go-rpio/v4"
	"time"
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
