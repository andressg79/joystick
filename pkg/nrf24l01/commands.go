package nrf24l01

const (
	// Read command and status registers. AAAAA =
	// 5 bit Register Map Address
	R_REGISTER = 0x00 // 000A AAAA

	// Write command and status registers. AAAAA = 5
	// bit Register Map Address
	// Executable in power down or standby modes
	// only.
	W_REGISTER = 0x20 // 001A AAAA

	// Read RX-payload: 1 – 32 bytes. A read operation
	// always starts at byte 0. Payload is deleted from
	// FIFO after it is read. Used in RX mode.
	R_RX_PAYLOAD = 0x61 // 0110 0001

	// Write TX-payload: 1 – 32 bytes. A write operation
	// always starts at byte 0 used in TX payload.
	W_TX_PAYLOAD = 0xA0 // 1010 0000

	// Flush TX FIFO, used in TX mode
	FLUSH_TX = 0xE1 // 1110 0001

	// Flush RX FIFO, used in RX mode
	// Should not be executed during transmission of
	// acknowledge, that is, acknowledge package will
	// not be completed.
	FLUSH_RX = 0xE2 // 1110 0010

	// Used for a PTX device
	// Reuse last transmitted payload.
	// TX payload reuse is active until
	// W_TX_PAYLOAD or FLUSH TX is executed. TX
	// payload reuse must not be activated or deactivated during package transmission.
	REUSE_TX_PL = 0xE3 // 1110 0011

	// Read RX payload width for the top
	// R_RX_PAYLOAD in the RX FIFO.
	// Note: Flush RX FIFO if the read value is larger
	// than 32 bytes.
	R_RX_PL_WID = 0x60 // 0110 0000

	// Used in RX mode.
	// Write Payload to be transmitted together with
	// ACK packet on PIPE PPP. (PPP valid in the
	// range from 000 to 101). Maximum three ACK
	// packet payloads can be pending. Payloads with
	// same PPP are handled using first in - first out
	// principle. Write payload: 1– 32 bytes. A write
	// operation always starts at byte 0.
	W_ACK_PAYLOAD = 0xA8 // 1010 1PPP

	// Used in TX mode. Disables AUTOACK on this
	// specific packet.
	W_TX_PAYLOAD_NO_ACK = 0xB0 // 1011 0000

	// No Operation. Might be used to read the STATUS
	// register
	NOP = 0xFF
)
