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

func main() {
	println("initializing...")

	machine.InitADC()

	joyLeft = joystick.NewJoystick(
		joystick.NewkHardware(
			machine.ADC{Pin: machine.ADC1},
			machine.ADC{Pin: machine.ADC0},
			machine.GPIO22,
		))
	joyLeft.Init()

	time.Sleep(time.Second)

	for {
		joyLeft.ReadUnit8()
		time.Sleep(time.Microsecond * 100)
	}

}
