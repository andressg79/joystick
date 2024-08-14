package main

/**
 * Copyright (c) 2024 Andres Sabini
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * Joystick - for Arduino Nano
 */

import (
	"joystick/internal/hardware/joystick"
	"machine"
	"time"
)

var jL1 *joystick.Joystick
var jL2 *joystick.Joystick
var jR1 *joystick.Joystick

func main() {
	println("initializing...")

	machine.InitADC()

	jL1 = joystick.NewJoystick(
		"L1",
		joystick.NewkHardware(
			machine.ADC{Pin: machine.ADC0},
			machine.ADC{Pin: machine.ADC1},
			machine.D12,
		))
	jL1.Init()

	jL2 = joystick.NewJoystick(
		"L2",
		joystick.NewkHardware(
			machine.ADC{Pin: machine.ADC2},
			machine.ADC{Pin: machine.ADC3},
			machine.D11,
		))
	jL2.Init()

	jR1 = joystick.NewJoystick(
		"R1",
		joystick.NewkHardware(
			machine.ADC{Pin: machine.ADC4},
			machine.ADC{Pin: machine.ADC5},
			machine.D10,
		))
	jR1.Init()

	time.Sleep(time.Second)

	for {
		xl1, yl1, sl1 := jL1.Read()
		xl2, yl2, sl2 := jL2.Read()
		xr1, yr1, sr1 := jR1.Read()
		println(xl1, "\t", yl1, "\t", sl1, "\t", xl2, "\t", yl2, "\t", sl2, "\t", xr1, "\t", yr1, "\t", sr1)
		time.Sleep(time.Microsecond * 100)
	}

}
