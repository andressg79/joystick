package main

import (
	"machine"
	"time"
)

func main() {
	var led = machine.GPIO13
	led.Configure(machine.PinConfig{
		Mode: machine.PinOutput,
	})
	for {
		led.High()
		time.Sleep(time.Millisecond * 500)
		led.Low()
		time.Sleep(time.Millisecond * 500)
		print(".")
	}
}
