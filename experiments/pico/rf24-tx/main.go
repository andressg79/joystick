package main

import (
	"machine"
	"time"

	nrf24l01 "joystick/pkg/nrf24l01"
)

const (
	BUFF_LENGTH = 4
)

var err error

func main() {

	time.Sleep(1 * time.Second)
	println("----------------TX MODE----------------\n")
	spi := machine.SPI0
	err = spi.Configure(machine.SPIConfig{
		SCK: machine.GPIO6,
		SDO: machine.GPIO7,
		SDI: machine.GPIO4,
	})
	if err != nil {
		println("error with spi.Configure(machine.SPIConfig{}):", err)
	}

	ce := machine.GPIO12 // Digital Input	Chip Enable Activates RX or TX mode
	csn := machine.GPIO5 // Digital Input	SPI Chip Select

	nrf := nrf24l01.New(spi, &ce, &csn)
	err = nrf.Configure()
	if err != nil {
		println("error with nrf.Configure():", err)
	}

	err = nrf.SetTXMode()
	if err != nil {
		println("error with nrf.SetTXMode():", err)
	}

	// Enable RX address for all data pipes
	err = nrf.EnableRXAddresses(true)
	if err != nil {
		println("error with nrf.EnableRXAddresses(true):", err)
	}

	// Disable Auto Acknowledgment
	err = nrf.EnableAutoAck(false)
	if err != nil {
		println("error with nrf.EnableAutoAck(false):", err)
	}

	// Disable dynamic payloads for all data pipes.
	err = nrf.EnableDynamicPayloads(false)
	if err != nil {
		println("error with nrf.EnableDynamicPayloads(false):", err)
	}

	// err = nrf.SetRFChannel(127)
	err = nrf.SetRFChannel(100)
	if err != nil {
		println("error with nrf.SetRFChannel(120):", err)
	}

	//err = nrf.SetRF1MBPS()
	// err = nrf.SetRF2MBPS()
	err = nrf.SetRF250KBPS()
	if err != nil {
		println("error with nrf.SetRF1MBPS():", err)
	}

	err = nrf.SetPipeRXPayloadWidth(0, byte(BUFF_LENGTH))
	if err != nil {
		println("error with nrf.SetPipeRXPayloadWidth(0, byte(BUFF_LENGTH)):", err)
	}

	res, err := nrf.GetRegisterState(nrf24l01.CONFIG)
	println("TX: CONFIG:     ", res)
	res, err = nrf.GetRegisterState(nrf24l01.EN_AA)
	println("TX: EN_AA:      ", res)
	res, err = nrf.GetRegisterState(nrf24l01.EN_RXADDR)
	println("TX: EN_RXADDR:  ", res)
	res, err = nrf.GetRegisterState(nrf24l01.SETUP_AW)
	println("TX: SETUP_AW:   ", res)
	res, err = nrf.GetRegisterState(nrf24l01.RF_CH)
	println("TX: RF_CH:      ", res)
	res, err = nrf.GetRegisterState(nrf24l01.RF_SETUP)
	println("TX: RF_SETUP:   ", res)
	status, err := nrf.GetRegisterState(nrf24l01.STATUS)
	println("TX: STATUS:     ", status)
	res, err = nrf.GetRegisterState(nrf24l01.RX_ADDR_P0)
	println("TX: RX_ADDR_P0: ", res)
	res, err = nrf.GetRegisterState(nrf24l01.RX_PW_P0)
	println("TX: RX_PW_P0:   ", res)
	res, err = nrf.GetRegisterState(nrf24l01.RX_PW_P1)
	println("TX: RX_PW_P1:   ", res)
	res, err = nrf.GetRegisterState(nrf24l01.FIFO_STATUS)
	println("TX: FIFO_STATUS:", res)
	feature, err := nrf.GetRegisterState(nrf24l01.FEATURE)
	err = nrf.SetRegisterState(nrf24l01.FEATURE, feature|0b00000001) // EN_DYN_ACK (W_TX_PAYLOAD_NOACK)
	feature, err = nrf.GetRegisterState(nrf24l01.FEATURE)
	println("TX: FEATURE:    ", feature)
	if err != nil {
		println(err)
	}

	println("--------------------------------------------------")

	w := make([]byte, BUFF_LENGTH)

	err = nrf.FlushTX()
	if err != nil {
		println(err)
	}

	for {

		w[0] = 0b10000000
		w[1] = 0b01000000
		w[2] = 0b00100000
		w[3] = 0b00010000
		err = nrf.TransmitDataWithAck(w) // it`s transmission
		// err = nrf.TransmitDataWithoutAck(w)
		if err != nil {
			println(err)
		}

		config, err := nrf.GetRegisterState(nrf24l01.CONFIG)
		println("TX: CONFIG:              ", config)
		res, _ = nrf.GetRegisterState(nrf24l01.RF_CH)
		println("TX: RF_CH:               ", res)
		res, _ := nrf.GetRegisterState(nrf24l01.OBSERVE_TX)
		println("TX: OBSERVE_TX:          ", res)
		res, _ = nrf.GetRegisterState(nrf24l01.FIFO_STATUS)
		println("TX: FIFO_STATUS:         ", res)
		println("TX: FIFO_STATUS.TX_REUSE:", res&0b01000000 > 0)
		println("TX: FIFO_STATUS.TX_FULL: ", res&0b00100000 > 0)
		res, _ = nrf.GetRegisterState(nrf24l01.CONFIG)
		println("TX: CONFIG:              ", res)
		status, _ = nrf.GetRegisterState(nrf24l01.STATUS)
		println("TX: STATUS:              ", status)
		println("TX: STATUS.TX_DS:        ", status&0b00100000>>5)
		res, _ = nrf.GetRegisterState(nrf24l01.FEATURE)
		println("TX: FEATURE:             ", res)

		println("------------------------------------------------")
		println("")
		time.Sleep(time.Second / 4)
		if err != nil { // Don't do that. Handle all errors separately )
			println(err)
		}
	}
}
