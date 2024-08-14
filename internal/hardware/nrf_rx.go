package hardware

import (
	"joystick/pkg/nrf24l01"
	"machine"
)

func NewRX(bufferLength int) *nrf24l01.Device {
	println("----------------RX MODE----------------\n")
	spi := machine.SPI0
	err := spi.Configure(machine.SPIConfig{
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

	err = nrf.SetPipeRXPayloadWidth(0, byte(bufferLength))
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
	r := make([]byte, bufferLength)
	for _, v := range r {
		print(" ", v)
	}
	println("")
	nrf.FlushRX()
	return nrf
}
