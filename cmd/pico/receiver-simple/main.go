package main

/**
 * Copyright (c) 2024 Andres Sabini
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * Receptor sencillo, con display SSD1306
 * Escucha por un canal RF24L01 y luego muestra en display SSD1306
 */

import (
	"bytes"
	"image/color"
	"joystick/internal/hardware"
	"joystick/pkg/nrf24l01"
	"machine"
	"strconv"
	"strings"
	"time"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/proggy"
)

const (
	BUFF_LENGTH   = 12
	RX_IDENTIFIER = 0b00000000
	TX_CONTROLLER = 0b00000001
)

func main() {
	time.Sleep(time.Second)
	println("init...")
	//dev := hardware.NewDisplay(hardware.NewIC2(machine.I2C1, machine.I2C1_SCL_PIN, machine.I2C1_SDA_PIN))
	dev := hardware.NewDisplay(hardware.NewIC2(machine.I2C1, machine.GPIO3, machine.GPIO2))

	println("init RF24L01")
	nrf := hardware.NewRX(BUFF_LENGTH)

	rfMessage := make([]byte, BUFF_LENGTH)

	for {
		// Esperar mensajes ...
		newMessage := rfReceive(nrf)

		if newMessage == nil {
			continue
		}

		if bytes.Equal(newMessage[:], make([]byte, BUFF_LENGTH)) {
			continue
		}

		if !bytes.Equal(newMessage[:], rfMessage[:]) {
			rfMessage = newMessage

			ints := make([]int, BUFF_LENGTH)
			strs := make([]string, BUFF_LENGTH)
			for i := 0; i < BUFF_LENGTH; i++ {
				ints[i] = int(rfMessage[i])
				strs[i] = strconv.FormatInt(int64(ints[i]), 10)
			}
			println("RX: ", strings.Join(strs, "\t"))

			dev.ClearDisplay()
			tinyfont.WriteLine(
				&dev,
				&proggy.TinySZ8pt7b,
				0,
				0x09,
				strings.Join(strs, ", "),
				color.RGBA{255, 255, 255, 255},
			)
			dev.Display()

			time.Sleep(time.Millisecond * 200)
		}
	}

}

// rfReceive receives data from the given nrf24l01 device.
//
// nrf - a pointer to nrf24l01.Device
// []byte - returns a byte slice
func rfReceive(nrf *nrf24l01.Device) []byte {
	status, err := nrf.GetStatus()
	fifoStatus, err := nrf.GetRegisterState(nrf24l01.FIFO_STATUS)
	if err != nil {
		println(err)
	}
	r := make([]byte, BUFF_LENGTH)
	if fifoStatus&TX_CONTROLLER == 0 {
		// RX
		err = nrf.ReceiveData(r)
		if err != nil {
			println(err)
		}

		// STATUS.RX_DR
		err = nrf.SetRegisterState(nrf24l01.STATUS, status&0b01000000)
		if err != nil {
			println(err)
		}

		// FIFO_STATUS.RX_EMPTY
		fifoStatus, err = nrf.GetRegisterState(nrf24l01.FIFO_STATUS)
		if err != nil {
			println(err)
		}
	}
	return r
}
