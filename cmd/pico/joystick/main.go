package main

/**
 * Copyright (c) 2024 Andres Sabini
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * Joystick master
 * Escribe por un canal RF24L01
 */

import (
	"errors"
	"joystick/internal/hardware"
	"joystick/pkg/nrf24l01"
	"machine"
	"runtime/debug"
	"time"
)

const (
	BUFF_LENGTH     = 12
	TX_IDENTIFIER   = 0b00000001 // Idetifier for TX
	PRINT_RF_STATUS = false      // print rf status
	PRINT_MSG       = true       // print message
)

// Dispositivo de destino
var destControlled uint8 = 0x00000000

var nrf *nrf24l01.Device

func main() {

	defer func() {
		if r := recover(); r != nil {
			println("stacktrace from panic: \n", string(debug.Stack()))
		}
	}()

	println("initializing...")

	machine.InitADC()

	time.Sleep(time.Second)

	println("init RF24L01")

	nrf = hardware.NewTX(BUFF_LENGTH)

	time.Sleep(time.Second)

	if nrf != nil {
		println("initialized RF24L01")
	}

	for {
		send, err := rfSend(nrf, prepareRFMessage())
		time.Sleep(time.Millisecond * 150)
	}

}

// rfSend sends data to the given nrf24l01 device.
func rfSend(nrf *nrf24l01.Device, w []byte) (bool, error) {

	if nrf == nil {
		return false, errors.New("no nrf24l01 device")
	}

	if err := nrf.TransmitDataWithoutAck(w); err != nil {
		return false, err
	}

	if PRINT_RF_STATUS {
		// Show status
		config, err := nrf.GetRegisterState(nrf24l01.CONFIG)
		println("TX: CONFIG:              ", config)

		res, _ := nrf.GetRegisterState(nrf24l01.RF_CH)
		println("TX: RF_CH:               ", res)

		res, _ = nrf.GetRegisterState(nrf24l01.OBSERVE_TX)
		println("TX: OBSERVE_TX:          ", res)

		res, _ = nrf.GetRegisterState(nrf24l01.FIFO_STATUS)
		println("TX: FIFO_STATUS:         ", res)
		println("TX: FIFO_STATUS.TX_REUSE:", res&0b01000000 > 0)
		println("TX: FIFO_STATUS.TX_FULL: ", res&0b00100000 > 0)

		status, _ := nrf.GetRegisterState(nrf24l01.STATUS)
		println("TX: STATUS:              ", status)
		println("TX: STATUS.TX_DS:        ", status&0b00100000>>5)

		res, _ = nrf.GetRegisterState(nrf24l01.FEATURE)
		println("TX: FEATURE:             ", res)
		println("------------------------------------------------")
		println("")
		if err != nil {
			return false, err
		}
	}
	return true, nil
}
