package main

/**
 * Copyright (c) 2024 Andres Sabini
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * Joystick
 */

import (
	"joystick/internal/hardware/joystick"
	"machine"
	"time"
)

var joyLeft *joystick.Joystick
var joyRight *joystick.Joystick

func main() {
	println("initializing...")

	machine.InitADC()

	joyLeft = joystick.NewJoystick(
		"Left",
		joystick.NewkHardware(
			machine.ADC{Pin: machine.ADC1},
			machine.ADC{Pin: machine.ADC0},
			machine.GPIO22,
		))
	joyLeft.Init()

	p := machine.GP20
	p.Configure(machine.PinConfig{Mode: machine.PinOutput})
	joyRight = joystick.NewJoystick(
		"Right",
		joystick.NewkHardware(
			machine.ADC{Pin: machine.ADC3},
			machine.ADC{Pin: machine.ADC2},
			machine.GPIO21,
		))
	joyRight.Init()

	time.Sleep(time.Second)

	for {
		lx, ly, ls := joyLeft.ReadUnit8()
		rx, ry, rs := joyRight.ReadUnit8()
		//println("Left:\t", joyLeft.X, "\t", joyLeft.Y, "\t", joyLeft.Sw, "\t\tRight:\t", joyRight.X, "\t", joyRight.Y, "\t", joyRight.Sw)
		println("Left:\t", lx, "\t", ly, "\t", ls, "\t\tRight:\t", rx, "\t", ry, "\t", rs)
		time.Sleep(time.Millisecond * 10)
	}

}
