package main

/**
 * Copyright (c) 2024 Andres Sabini
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * Joystick - for Arduino Nano
 */

import (
	"encoding/json"
	"joystick/internal/hardware/joystick"
	"machine"
	"time"
)

var addr uint16 = 1
var jL *joystick.Joystick
var jR *joystick.Joystick

func main() {
	println("initializing...")

	machine.InitADC()

	jL = joystick.NewJoystick(
		"L",
		joystick.NewkHardware(
			machine.ADC{Pin: machine.ADC0},
			machine.ADC{Pin: machine.ADC1},
			machine.D12,
		))
	jL.Init()

	jR = joystick.NewJoystick(
		"R",
		joystick.NewkHardware(
			machine.ADC{Pin: machine.ADC2},
			machine.ADC{Pin: machine.ADC3},
			machine.D11,
		))
	jR.Init()

	i2c := machine.I2C0
	i2c.Configure(machine.I2CConfig{
		Frequency: 100000,
	})

	time.Sleep(time.Second)
	for {
		xl, yl, sl := jL.Read()
		xr, yr, sr := jR.Read()
		d := Data{
			xl: xl,
			yl: yl,
			sl: sl,
			xr: xr, yr: yr, sr: sr,
		}
		data, err := json.Marshal(d)
		if err != nil {
			println(err)
			continue
		}
		r := make([]byte, 1)
		if err := i2c.Tx(addr, data, r); err != nil {
			println(err)
		}
	}
}

type Data struct {
	xl uint16
	yl uint16
	sl bool
	xr uint16
	yr uint16
	sr bool
}
