package main

import (
	"machine"
	"time"

	nrf24l01 "joystick/pkg/nrf24l01"
)

const (
	BUFF_LENGTH = 4
)

var (
	err error
)

func main() {

	time.Sleep(1 * time.Second)

	println("----------------RX MODE----------------\n")
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
		println("failed with nrf.Configure():", err)
	}

	err = nrf.SetRXMode()
	if err != nil {
		println("failed with nrf.SetRXMode():", err)
	}

	// Enable RX address for all data pipes
	err = nrf.EnableRXAddresses(true)
	if err != nil {
		println("failed with nrf.EnableRXAddresses(true):", err)
	}

	// Disable Auto Acknowledgment
	err = nrf.EnableAutoAck(false)
	if err != nil {
		println("failed with nrf.EnableAutoAck(false):", err)
	}

	// Disable dynamic payloads for all data pipes.
	err = nrf.EnableDynamicPayloads(false)
	if err != nil {
		println("failed with nrf.EnableDynamicPayloads(false):", err)
	}

	// err = nrf.SetRFChannel(120)
	err = nrf.SetRFChannel(100)
	if err != nil {
		println("error with nrf.SetRFChannel(120):", err)
	}

	//err = nrf.SetRF1MBPS()
	// err = nrf.SetRF2MBPS()
	err = nrf.SetRF250KBPS()
	if err != nil {
		println("failed with nrf.SetRF1MBPS():", err)
	}

	err = nrf.SetPipeRXPayloadWidth(0, byte(BUFF_LENGTH))
	if err != nil {
		println("failed with nrf.SetPipeRXPayloadWidth(0, byte(BUFF_LENGTH)):", err)
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
	println("TX: FEATURE:    ", feature)
	if err != nil {
		println(err)
	}

	println("--------------------------------------------------")

	r := make([]byte, BUFF_LENGTH)
	for _, v := range r {
		print(" ", v)
	}
	println("")

	nrf.FlushRX()

	// Try read PAYLOAD
	for {
		status, err = nrf.GetStatus()
		fifoStatus, err := nrf.GetRegisterState(nrf24l01.FIFO_STATUS)
		if err != nil {
			println(err)
		}

		if fifoStatus&0b00000001 == 0 {
			println("STATUS.RX_DR:        ", status&0b01000000>>6)
			println("STATUS.RX_P_NO:      ", status&0b00001110>>1)
			println("FIFO_STATUS.RX_FULL: ", fifoStatus&0b00000010>>1)
			println("FIFO_STATUS.RX_EMPTY:", fifoStatus&0b00000001)
			// println("FIFO_STATUS:   ", fifoStatus)
			err = nrf.ReceiveData(r)
			if err != nil {
				println(err)
			}

			for _, v := range r {
				print(" ", v)
			}
			println("")

			// STATUS.RX_DR
			err = nrf.SetRegisterState(nrf24l01.STATUS, status&0b01000000)
			if err != nil {
				println(err)
			}

			fifoStatus, err = nrf.GetRegisterState(nrf24l01.FIFO_STATUS)
			println("FIFO_STATUS:         ", fifoStatus)
			if err != nil {
				println(err)
			}
		}
		//println("----------------------------------------")
		//println("")

	}
}
