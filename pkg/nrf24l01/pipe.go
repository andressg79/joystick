package nrf24l01

// var addressLengthArray = [3]byte{3, 4, 5}

type Pipe struct {
	EnableAutoAck        byte // 1 - yes / 0 - no
	EnableRXAddress      byte // 1 - yes / 0 - no
	RXAddress            byte
	FullRXAddress        []byte
	RXPayloadWidth       byte
	DynamicPayloadLength byte // 1 - yes / 0 - no
}

// Get array of registers of pipes RX address
func GetPipesRXAddressRegisters() [6]byte {
	return [...]byte{RX_ADDR_P0, RX_ADDR_P1, RX_ADDR_P2, RX_ADDR_P3, RX_ADDR_P4, RX_ADDR_P5}
}

// Get array of registers of pipes RX payload width
func GetPipesRXPayloadWidthRegisters() [6]byte {
	return [...]byte{RX_PW_P0, RX_PW_P1, RX_PW_P2, RX_PW_P3, RX_PW_P4, RX_PW_P5}
}

// Read pipe config from nrf24l01 to "pipe".
//
// n - pipe number. Allowed only 0...5.
//
// pipe - strusture for loading data.
func (d *Device) GetPipeConfig(n byte, pipe *Pipe) error {
	pipe.EnableAutoAck, err = d.GetRegisterState(EN_AA)
	if err != nil {
		return err
	}
	pipe.EnableAutoAck = pipe.EnableAutoAck >> n & 1

	pipe.EnableRXAddress, err = d.GetRegisterState(EN_RXADDR)
	if err != nil {
		return err
	}
	pipe.EnableRXAddress = pipe.EnableRXAddress >> n & 1

	pipe.RXAddress, err = d.GetRegisterState(GetPipesRXAddressRegisters()[n])
	if err != nil {
		return err
	}

	err = d.GetFullPipeRXAddress(n, pipe.FullRXAddress)
	if err != nil {
		return err
	}

	pipe.RXPayloadWidth, err = d.GetRegisterState(GetPipesRXPayloadWidthRegisters()[n])
	if err != nil {
		return err
	}

	pipe.DynamicPayloadLength, err = d.GetRegisterState(DYNPD)
	if err != nil {
		return err
	}
	pipe.EnableRXAddress = pipe.EnableRXAddress >> n & 1

	return nil
}

// Write pipe config from "pipe" to nrf24l01.
//
// n - pipe number. Allowed only 0...5.
//
// pipe - strusture with data.
func (d *Device) SetPipeConfig(n byte, pipe *Pipe) error {
	res, err := d.GetRegisterState(EN_AA)
	if err != nil {
		return err
	}
	err = d.SetRegisterState(EN_AA, pipe.EnableAutoAck<<n|res)
	if err != nil {
		return err
	}

	res, err = d.GetRegisterState(EN_RXADDR)
	if err != nil {
		return err
	}
	err = d.SetRegisterState(EN_AA, pipe.EnableRXAddress<<n|res)
	if err != nil {
		return err
	}

	err = d.SetRegisterState(GetPipesRXAddressRegisters()[n], pipe.RXAddress)
	if err != nil {
		return err
	}

	err = d.SetFullPipeRXAddress(n, pipe.FullRXAddress)
	if err != nil {
		return err
	}

	err = d.SetRegisterState(GetPipesRXPayloadWidthRegisters()[n], pipe.RXPayloadWidth)
	if err != nil {
		return err
	}

	res, err = d.GetRegisterState(DYNPD)
	if err != nil {
		return err
	}
	err = d.SetRegisterState(DYNPD, pipe.DynamicPayloadLength<<n|res)

	return nil
}
