package keyboard

import (
	go2c "github.com/d2r2/go-i2c"
	"time"
)

const (
	retries = 10
	KeyEsc  = "ESCAPE"
)

func Write(i2c *go2c.I2C, msg string) error {
	err := SendMessage(i2c, []byte(msg))
	if err != nil {
		return err
	}
	return WaitUntilDone(i2c)
}

func SendMessage(i2c *go2c.I2C, msg []byte) error {
	const startOfTransmission byte = 0x01
	const endOfTransmission byte = 0x04

	msgBytes := msg
	msgBytes = append([]byte{startOfTransmission}, msgBytes...)
	msgBytes = append(msgBytes, endOfTransmission)

	sent := 0
	const maxSize = 32
	for i := 0; i < len(msgBytes); i++ {
		if (i > 0 && i%maxSize == maxSize-1) || i == len(msgBytes)-1 {
			var err error
			for j := 0; j < retries; j++ {
				_, err = i2c.WriteBytes(msgBytes[sent*maxSize : i+1])
				if err == nil {
					break
				}
				time.Sleep(50 * time.Millisecond)
			}
			if err != nil {
				return err
			}
			sent++
		}
	}
	return nil
}

func WaitUntilDone(i2c *go2c.I2C) error {
	for {
		buffer := make([]byte, 1)
		_, err := i2c.ReadBytes(buffer)
		if err != nil {
			return err
		}
		if buffer[0] == 0x1 {
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}
}
