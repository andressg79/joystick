# nRF24L01 driver

Datasheet https://www.mouser.com/datasheet/2/297/nRF24L01_Product_Specification_v2_0-9199.pdf

## Strong recommendation

Solder capacitors on the VCC and GND pins.
Two capacitors: ceramic at 0.1uF and electrolytic at 100uF.
Ceramic - filter of high frequency of power.
Electrolytic - filter of low frequency of power.

Solder it in any cases, even if you use NRF24L01-adapter.

## Example

TX MODE 
```go
package main

import (
	"machine"
	"time"

    nrf24l01 "tinygo.org/x/drivers/nrf24l01"
)


const (
	BUFF_LENGTH = 4
)

var err error

func main() {

	println("----------------TX MODE----------------\n")
	spi := machine.SPI0
	err = spi.Configure(machine.SPIConfig{})
	if err != nil {
		println("error with spi.Configure(machine.SPIConfig{}):", err)
	}

	ce := machine.D9   // Digital Input	Chip Enable Activates RX or TX mode
	csn := machine.D10 // Digital Input	SPI Chip Select

	nrf := nrf24l01.New(&spi, &ce, &csn)
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

	err = nrf.SetRF1MBPS()
	// err = nrf.SetRF2MBPS()
	// err = nrf.SetRF250KBPS()
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
```

RX MODE
```go
package main

import (
	"machine"

    nrf24l01 "tinygo.org/x/drivers/nrf24l01"
)

const (
	BUFF_LENGTH = 4
)

var (
	err error
)

func main() {
	println("----------------RX MODE----------------\n")
	spi := machine.SPI0
	err = spi.Configure(machine.SPIConfig{})
	if err != nil {
		println("failed with spi.Configure(machine.SPIConfig{}):", err)
	}

	ce := machine.D9   // Digital Input	Chip Enable Activates RX or TX mode
	csn := machine.D10 // Digital Input	SPI Chip Select

	nrf := nrf24l01.New(&spi, &ce, &csn)
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

	err = nrf.SetRF1MBPS()
	// err = nrf.SetRF2MBPS()
	// err = nrf.SetRF250KBPS()
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
		println("----------------------------------------")
		println("")

	}
}

```

And finally scanning of air
```go
package main

import (
	"machine"

    nrf24l01 "tinygo.org/x/drivers/nrf24l01"
)

const (
	BUFF_LENGTH = 4
)

var (
	err error
)

func main() {
	println("----------------RX MODE scanning----------------\n")
	spi := machine.SPI0
	err = spi.Configure(machine.SPIConfig{})
	if err != nil {
		println("failed with spi.Configure(machine.SPIConfig{}):", err)
	}

	ce := machine.D9   // Digital Input	Chip Enable Activates RX or TX mode
	csn := machine.D10 // Digital Input	SPI Chip Select

	nrf := nrf24l01.New(&spi, &ce, &csn)
	err = nrf.Configure()
	if err != nil {
		println("failed with nrf.Configure():", err)
	}

	err = nrf.SetRXMode()
	if err != nil {
		println("failed with nrf.SetRXMode():", err)
	}

	err = nrf.SetRF1MBPS()
	// err = nrf.SetRF2MBPS()
	// err = nrf.SetRF250KBPS()
	if err != nil {
		println("failed with nrf.SetRF1MBPS():", err)
	}

	res, err := nrf.GetRegisterState(nrf24l01.CONFIG)
	println("CONFIG:     ", res)
	res, err = nrf.GetRegisterState(nrf24l01.SETUP_AW)
	println("SETUP_AW:   ", res)
	res, err = nrf.GetRegisterState(nrf24l01.RF_SETUP)
	println("RF_SETUP:   ", res)
	status, err := nrf.GetRegisterState(nrf24l01.STATUS)
	println("STATUS:     ", status)
	fifoStatus, err := nrf.GetRegisterState(nrf24l01.FIFO_STATUS)
	println("FIFO_STATUS:", fifoStatus)
	res, err = nrf.GetRegisterState(nrf24l01.FEATURE)
	println("FEATURE:    ", res)
	println("--------------------------------------------------")

	err = nrf.HearChannels()
    if err != nil { // Don't do that. Handle all errors separately )
        println(err)
    }

}

```

