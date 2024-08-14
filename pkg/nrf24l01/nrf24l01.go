package nrf24l01

// Datasheet:
// https://www.mouser.com/datasheet/2/297/nRF24L01_Product_Specification_v2_0-9199.pdf
// https://github.com/tinygo-org/drivers/issues/245

import (
	"machine"
	"time"
)

var err error

type Device struct {
	spi *machine.SPI // Digital Input	bus

	ce  *machine.Pin // Digital Input	Chip Enable Activates RX or TX mode
	csn *machine.Pin // Digital Input	SPI Chip Select
}

// New returns a new NRF device.
// spi - SPI bus with LSBFirst = false (defualt).
// ce - use for switch radio to RX or TX.
// csn - begin-end of transmitting to NRF.
func New(spi *machine.SPI, ce, csn *machine.Pin) *Device {
	return &Device{
		spi: spi,
		ce:  ce,
		csn: csn,
	}
}

// Set start values of parameters.
func (d *Device) Configure() error {
	d.ce.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d.csn.Configure(machine.PinConfig{Mode: machine.PinOutput})

	d.Enable()   // start position
	d.csn.High() // start position

	// 1500uS timeouts - minimum for 32B payload in ESB@250KBPS.
	// WARNING: If this is ever lowered, either 250KBS mode with AA is broken or maximum packet
	// sizes must never be used. See documentation for a more complete explanation.
	// Source: https://github.com/AlexGyver/nRF24L01/blob/master/RF24-master/RF24.cpp
	// SETUP_RETR = 0x04
	// ARD = 0b01000000  +  ARC = 0b00001111
	err = d.SetAutoRetransmission(15, 15)
	if err != nil {
		return err
	}

	err = d.Set8CRC()
	if err != nil {
		return err
	}

	err = d.SetRFChannel(0)
	if err != nil {
		return err
	}

	// err = d.SetRF1MBPS() // 0b00 – 1Mbps
	// if err != nil {
	// 	return err
	// }

	// err = d.SetRFPower(0b11) // 0b11 – 0dBm
	// if err != nil {
	// 	return err
	// }
	err = d.SetRegisterState(RF_SETUP, 0b00000111) // 0dBm + 1Mbps
	if err != nil {
		return err
	}

	err = d.SetRegisterState(FEATURE, 0)
	if err != nil {
		return err
	}

	d.FlushRX()
	d.FlushTX()

	d.Disable() // start position

	return nil
}

// Enable chip
func (d *Device) Enable() {
	d.ce.High()
}

// Disable chip
func (d *Device) Disable() {
	d.ce.Low()
}

// Get STATUS-register state.
func (d *Device) GetStatus() (byte, error) {
	d.csn.Low()
	defer d.csn.High()
	return d.spi.Transfer(NOP)
}

// Get register state.
// r - register.
func (d *Device) GetRegisterState(r byte) (byte, error) {
	d.csn.Low()
	defer d.csn.High()
	_, err = d.spi.Transfer(r)
	if err != nil {
		// return 0, errors.New("Error with request d.spi.Transfer(r) in nrf24l01.GetRegisterState(r): " + err.Error())
		return 0, err
	}
	return d.spi.Transfer(NOP)
}

// Set register state.
// r - register.
// s - state which have to be set into register.
func (d *Device) SetRegisterState(r byte, s byte) error {
	d.csn.Low()
	defer d.csn.High()
	_, err = d.spi.Transfer(0x20 + r)
	if err != nil {
		return err
	}
	_, err = d.spi.Transfer(s)
	if err != nil {
		return err
	}
	return nil
}

// Set nrf24l01 to RX-mode.
func (d *Device) SetRXMode() error {
	// nRF24L01P_Product_Specification_1_0.
	// 6.1.4 RX mode, p. 23.

	// Set PRIM_RX = 1, PWR_UP = 1.
	state, err := d.GetRegisterState(CONFIG)
	if err != nil {
		return err
	}
	err = d.SetRegisterState(CONFIG, state|0b00000011) // PWR_UP, PRIM_RX
	if err != nil {
		return err
	}
	// d.powerUp = true

	// Wait while mode is switching. PLL settling delay (130µs). Specs p.22.
	time.Sleep(130 * time.Microsecond)

	// Switch to work mode.
	d.Enable()

	// d.rxMode = true

	return nil
}

// Set nrf24l01 to TX-mode.
func (d *Device) SetTXMode() error {

	// Set PRIM_RX = 0, PWR_UP = 1.
	state, err := d.GetRegisterState(CONFIG)
	if err != nil {
		return err
	}
	err = d.SetRegisterState(CONFIG, state|0b00000011-0b00000001) // PWR_UP, PRIM_RX
	if err != nil {
		return err
	}
	// d.powerUp = true

	// Switch to work mode.
	d.Enable()
	// Wait while mode is switching. PLL settling delay (130µs). Specs p.22.
	time.Sleep(130 * time.Microsecond)

	// d.rxMode = false

	return nil
}

// Set RF channel.
// c - channel. Have range only from 0 to 127.
func (d *Device) SetRFChannel(c byte) error {

	err = d.SetRegisterState(RF_CH, c)
	if err != nil {
		return err
	}

	// d.channel = c

	return nil
}

// Set nrf24l01 to power down mode.
func (d *Device) SetPowerDownMode() error {
	c, err := d.GetRegisterState(CONFIG)
	if err != nil {
		return err
	}
	err = d.SetRegisterState(CONFIG, c|0b00000010-0b00000010)
	if err != nil {
		return err
	}

	// d.powerUp = false

	return nil
}

// Set nrf24l01 to power up mode.
func (d *Device) SetPowerUpMode() error {
	c, err := d.GetRegisterState(CONFIG)
	if err != nil {
		return err
	}
	err = d.SetRegisterState(CONFIG, c|0b00000010)
	if err != nil {
		return err
	}

	return nil
}

//  Flush TX FIFO, used in TX mode.
// TODO: Need add check that set TX mode.
func (d *Device) FlushTX() error {
	d.csn.Low()
	defer d.csn.High()
	_, err = d.spi.Transfer(FLUSH_TX)
	if err != nil {
		return err
	}

	return nil
}

// Flush RX FIFO, used in RX mode
// Should not be executed during transmission of
// acknowledge, that is, acknowledge package will
// not be completed.
// TODO: Need add check that set RX mode.
func (d *Device) FlushRX() error {
	d.csn.Low()
	defer d.csn.High()
	_, err = d.spi.Transfer(FLUSH_RX)
	if err != nil {
		return err
	}

	return nil
}

// Set 16-bit CRC.
func (d *Device) Set16CRC() error {
	state, err := d.GetRegisterState(CONFIG)
	if err != nil {
		return err
	}
	err = d.SetRegisterState(CONFIG, state|0b00000100)
	if err != nil {
		return err
	}
	// d.crc = 16
	return nil
}

// Set 8-bit CRC.
func (d *Device) Set8CRC() error {
	state, err := d.GetRegisterState(CONFIG)
	if err != nil {
		return err
	}
	err = d.SetRegisterState(CONFIG, state|0b00000100-0b00000100)
	if err != nil {
		return err
	}
	// d.crc = 8
	return nil
}

// Setup of Automatic Retransmission
//
// delay - Auto Retransmit Delay
//
// ‘0000’ – Wait 250µS;
// ‘0001’ – Wait 500µS;
// ‘0010’ – Wait 750µS;
// ……..;
// ‘1111’ – Wait 4000µS
//
// (Delay defined from end of transmission to start of
// next transmission)
//
// count - Auto Retransmit Count
//
// ‘0000’ – Re-Transmit disabled
// ‘0001’ – Up to 1 Re-Transmit on fail of AA
// ……
// ‘1111’ – Up to 15 Re-Transmit on fail of AA
func (d *Device) SetAutoRetransmission(delay byte, count byte) error {
	err = d.SetRegisterState(SETUP_RETR, delay<<4+count)
	if err != nil {
		return err
	}
	// d.autoRetransmitDelay = delay
	// d.autoRetransmitCount = count
	return nil
}

// // Get data rate
// func (d *Device) GetDataRate() byte {
// 	return d.dataRate
// }

// // Set data rate.
// // Allowed only:
// // 0b00 – 1Mbps, 0b01 – 2Mbps, 0b10 – 250kbps
// func (d *Device) SetDataRate(r byte) error {
// 	state, err := d.GetRegisterState(RF_SETUP)
// 	if err != nil {
// 		return err
// 	}

// 	switch r {
// 	case 0b10: // ‘10’ – 250kbps
// 		err = d.SetRegisterState(RF_SETUP, state|0b00101000-0b00001000) // [RF_DR_LOW, RF_DR_HIGH]
// 	case 0b01: // ‘01’ – 2Mbps
// 		err = d.SetRegisterState(RF_SETUP, state|0b00101000-0b00100000) // [RF_DR_LOW, RF_DR_HIGH]
// 	case 0b00: // ‘00’ – 1Mbps
// 		err = d.SetRegisterState(RF_SETUP, state|0b00101000-0b00101000) // [RF_DR_LOW, RF_DR_HIGH]
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	// d.dataRate = r
// 	return nil
// }

// // Get power
// func (d *Device) GetPower() byte {
// 	return d.rfPower
// }

// Set power.
// Allowed only:
// 0b00 – -18dBm, 0b01 – -12dBm, 0b10 – -6dBm, 0b11 – 0dBm
func (d *Device) SetRFPower(p byte) error {
	state, err := d.GetRegisterState(RF_SETUP)
	if err != nil {
		return err
	}

	switch p {
	case 0b00: // '00' – -18dBm
		err = d.SetRegisterState(RF_SETUP, state|0b00000110-0b00000110) // RF_PWR 2:1
	case 0b01: // '01' – -12dBm
		err = d.SetRegisterState(RF_SETUP, state|0b00000110-0b00000100) // RF_PWR 2:1
	case 0b10: // '10' – -6dBm
		err = d.SetRegisterState(RF_SETUP, state|0b00000110-0b00000010) // RF_PWR 2:1
	case 0b11: // '11' – 0dBm
		err = d.SetRegisterState(RF_SETUP, state|0b00000110) // RF_PWR 2:1
	}

	if err != nil {
		return err
	}

	// d.rfPower = p
	return nil
}

// Set radio frequency 250kbps
func (d *Device) SetRF250KBPS() error {
	state, err := d.GetRegisterState(RF_SETUP)
	if err != nil {
		return err
	}

	err = d.SetRegisterState(RF_SETUP, state|0b00101000-0b00001000)
	if err != nil {
		return err
	}

	return nil
}

// Set radio frequency 1Mbps
func (d *Device) SetRF1MBPS() error {
	state, err := d.GetRegisterState(RF_SETUP)
	if err != nil {
		return err
	}

	err = d.SetRegisterState(RF_SETUP, state|0b00101000-0b00101000)
	if err != nil {
		return err
	}

	return nil
}

// Set radio frequency 2Mbps
func (d *Device) SetRF2MBPS() error {
	state, err := d.GetRegisterState(RF_SETUP)
	if err != nil {
		return err
	}

	err = d.SetRegisterState(RF_SETUP, state|0b00101000-0b00100000)
	if err != nil {
		return err
	}

	return nil
}

// Get RX payload width
func (d *Device) GetRXPayloadWidth() (byte, error) {
	width, err := d.GetRegisterState(R_RX_PL_WID)
	if err != nil {
		return 0, err
	}
	return width, nil
}

// Enable ‘Auto Acknowledgment’ for all data pipes.
func (d *Device) EnableAutoAck(enable bool) error {
	if enable {
		err = d.SetRegisterState(EN_AA, 0b00111111)
		if err != nil {
			return err
		}
		// TODO: setting pipes EN_AA
	} else {
		err = d.SetRegisterState(EN_AA, 0b00000000)
		if err != nil {
			return err
		}
		// TODO: setting pipes EN_AA
	}
	return nil
}

// Enable ‘Auto Acknowledgment’ for separate data pipe.
// 0b00111111 - one bit - one pipe. From 0 to 5.
// Example: pipe = 0b00000010 - it's pipe 1.
// Can set multiple pipes.
func (d *Device) EnablePipesAutoAck(pipes byte, enable bool) error {
	state, err := d.GetRegisterState(EN_AA)
	if err != nil {
		return err
	}

	if enable {
		err = d.SetRegisterState(EN_AA, state|pipes)
		if err != nil {
			return err
		}
	} else {
		err = d.SetRegisterState(EN_AA, state|pipes-pipes)
		if err != nil {
			return err
		}
	}

	return nil
}

// Enable dynamic payloads for all data pipes.
func (d *Device) EnableDynamicPayloads(enable bool) error {
	if enable {
		err = d.SetRegisterState(DYNPD, 0b00111111)
		if err != nil {
			return err
		}
	} else {
		err = d.SetRegisterState(DYNPD, 0b00000000)
		if err != nil {
			return err
		}
	}
	return nil
}

// Enable dynamic payloads for separate data pipe.
// 0b00111111 - one bit - one pipe. From 0 to 5.
// Example: pipe = 0b00000010 - it's pipe 1.
// Can set multiple pipes.
func (d *Device) EnablePipesDynamicPayloads(pipes byte, enable bool) error {
	state, err := d.GetRegisterState(EN_AA)
	if err != nil {
		return err
	}

	if enable {
		err = d.SetRegisterState(EN_AA, state|pipes)
		if err != nil {
			return err
		}
	} else {
		err = d.SetRegisterState(EN_AA, state|pipes-pipes)
		if err != nil {
			return err
		}
	}

	return nil
}

// Enable RX address for all data pipes.
func (d *Device) EnableRXAddresses(enable bool) error {
	if enable {
		err = d.SetRegisterState(EN_RXADDR, 0b00111111)
		if err != nil {
			return err
		}
	} else {
		err = d.SetRegisterState(EN_RXADDR, 0b00000000)
		if err != nil {
			return err
		}
	}
	return nil
}

// Enable RX address for separate data pipe.
// 0b00111111 - one bit - one pipe. From 0 to 5.
// Example: pipe = 0b00000010 - it's pipe 1.
// Can set multiple pipes.
func (d *Device) EnablePipesRXAddresses(pipes byte, enable bool) error {
	state, err := d.GetRegisterState(EN_RXADDR)
	if err != nil {
		return err
	}

	if enable {
		err = d.SetRegisterState(EN_RXADDR, state|pipes)
		if err != nil {
			return err
		}
	} else {
		err = d.SetRegisterState(EN_RXADDR, state|pipes-pipes)
		if err != nil {
			return err
		}
	}

	return nil
}

// Get pipe RX address.
//
// pipe - number of pipe. From 0 to 5.
func (d *Device) GetPipeRXAddress(pipe byte) (byte, error) {
	return d.GetRegisterState(GetPipesRXAddressRegisters()[pipe])
}

// Get full pipe RX address.
//
// pipe - number of pipe. From 0 to 5.
//
// address - slice with address. Length must be meeting to register SETUP_AW.
func (d *Device) GetFullPipeRXAddress(pipe byte, address []byte) error {
	d.csn.Low()
	defer d.csn.High()

	_, err = d.spi.Transfer(GetPipesRXAddressRegisters()[pipe])
	if err != nil {
		return err
	}
	err = d.spi.Tx(nil, address)
	if err != nil {
		return err
	}

	return nil
}

// Set pipe RX address.
//
// pipe - number of pipe. From 0 to 5.
//
// address - address of pipe.
func (d *Device) SetPipeRXAddress(pipe, address byte) error {
	return d.SetRegisterState(GetPipesRXAddressRegisters()[pipe], address)
}

// Set full pipe RX address.
//
// pipe - number of pipe. From 0 to 5.
//
// address - slice with address. Length must be meeting to register SETUP_AW.
func (d *Device) SetFullPipeRXAddress(pipe byte, address []byte) error {
	d.csn.Low()
	defer d.csn.High()

	_, err = d.spi.Transfer(0x20 + GetPipesRXAddressRegisters()[pipe])
	if err != nil {
		return err
	}
	err = d.spi.Tx(address, nil)
	if err != nil {
		return err
	}

	return nil
}

// Get pipe RX payload width.
//
// pipe - number of pipe. From 0 to 5.
func (d *Device) GetPipeRXPayloadWidth(pipe byte) (byte, error) {
	return d.GetRegisterState(GetPipesRXPayloadWidthRegisters()[pipe])
}

// Set pipe RX payload width.
//
// pipe - number of pipe. From 0 to 5.
//
// width - width of pipe (1 to 32 bytes).
//
// 0 Pipe not used
//
// 1 = 1 byte
//
// …
//
// 32 = 32 bytes
func (d *Device) SetPipeRXPayloadWidth(pipe, width byte) error {
	return d.SetRegisterState(GetPipesRXPayloadWidthRegisters()[pipe], width)
}

// Transmit data whith AUTOACK on this specific packet to TX pipe address.
//
// w - writer. It's width of pipe (1 to 32 bytes).
// Length of w must be the same as width of pipe on receiver.
func (d *Device) TransmitDataWithAck(w []byte) error {
	d.csn.Low()
	defer d.csn.High()

	_, err = d.spi.Transfer(W_TX_PAYLOAD)
	if err != nil {
		return err
	}
	err = d.spi.Tx(w, nil)
	if err != nil {
		return err
	}
	return nil
}

// Transmit data whithout AUTOACK on this specific packet to TX pipe address.
//
// w - writer. It's width of pipe (1 to 32 bytes).
// Length of w must be the same as width of pipe on receiver.
func (d *Device) TransmitDataWithoutAck(w []byte) error {
	d.csn.Low()
	defer d.csn.High()

	_, err = d.spi.Transfer(W_TX_PAYLOAD_NO_ACK)
	if err != nil {
		return err
	}
	err = d.spi.Tx(w, nil)
	if err != nil {
		return err
	}
	return nil
}

// Receive data to TX pipe address.
//
// r - Reader. Into this parameter return data. It's width of pipe (1 to 32 bytes).
// Length of r must be the same as width of receiver pipe.
func (d *Device) ReceiveData(r []byte) error {
	d.csn.Low()
	defer d.csn.High()

	_, err = d.spi.Transfer(R_RX_PAYLOAD)
	if err != nil {
		return err
	}
	err = d.spi.Tx(nil, r)
	if err != nil {
		return err
	}

	return nil
}

// Hear carrier on channels
//
// Only for test RX mode.
func (d *Device) HearChannels() error {
	// Header channels numbers
	for i := 0; i < 128; i++ {
		if i < 100 {
			print(0)
		} else {
			print(1)
		}
	}
	println("")
	for i := 0; i < 128; i++ {
		if i/100 > 0 {
			print((i % 100) / 10)
		} else {
			print(i / 10)
		}
	}
	println("")
	for i := 0; i < 128; i++ {
		if i/10 > 0 {
			print(i % 10)
		} else {
			print(i)
		}
	}
	println("")

	// Try read PAYLOAD
	for {
		for i := 0; i < 128; i++ {
			err = d.SetRFChannel(byte(i))
			if err != nil {
				return err
			}

			time.Sleep(1000 * time.Microsecond)
			res, err := d.GetRegisterState(RPD)
			if err != nil {
				return err
			}

			d.FlushRX()
			d.FlushTX()

			print(res)
		}
		println("")
		time.Sleep(time.Second / 2)
	}
}
