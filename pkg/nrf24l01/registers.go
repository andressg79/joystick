package nrf24l01

const (
	// ------------------- CONFIG -------------------
	// Configuration Register
	CONFIG = 0x00

	// // Mask interrupt caused by RX_DR
	// // 1: Interrupt not reflected on the IRQ pin
	// // 0: Reflect RX_DR as active low interrupt on the
	// // IRQ pin
	// CONFIG_MASK_RX_DR = 0b01000000 // 6 R/W

	// // Mask interrupt caused by TX_DS
	// // 1: Interrupt not reflected on the IRQ pin
	// // 0: Reflect TX_DS as active low interrupt on the IRQ
	// // pin
	// CONFIG_MASK_TX_DS = 0b00100000 // 5 R/W

	// // Mask interrupt caused by MAX_RT
	// // 1: Interrupt not reflected on the IRQ pin
	// // 0: Reflect MAX_RT as active low interrupt on the
	// // IRQ pin
	// CONFIG_MASK_MAX_RT = 0b00010000 // 4 R/W

	// // Enable CRC. Forced high if one of the bits in the
	// // EN_AA is high
	// CONFIG_EN_CRC = 0b00001000 // 3  R/W

	// // CRC encoding scheme
	// // '0' - 1 byte
	// // '1' – 2 bytes
	// CONFIG_CRCO = 0b00000100 // 2 R/W

	// // 1: POWER UP, 0:POWER DOWN
	// CONFIG_PWR_UP = 0b00000010 // 1 R/W

	// // RX/TX control
	// // 1: PRX, 0: PTX
	// CONFIG_PRIM_RX = 0b00000001 // 0 R/W

	// --------------------- EN_AA ---------------------
	// Enhanced ShockBurst™
	// Enable ‘Auto Acknowledgment’ Function Disable
	// this functionality to be compatible with nRF2401
	EN_AA = 0x01

	// Enable auto acknowledgement data pipe 5
	ENAA_P5 = 0b00100000

	// Enable auto acknowledgement data pipe 4
	ENAA_P4 = 0b00010000

	// Enable auto acknowledgement data pipe 3
	ENAA_P3 = 0b00001000

	// Enable auto acknowledgement data pipe 2
	ENAA_P2 = 0b00000100

	// Enable auto acknowledgement data pipe 1
	ENAA_P1 = 0b00000010

	// Enable auto acknowledgement data pipe 0
	ENAA_P0 = 0b00000001

	// ------------------- EN_RXADDR -------------------
	// Enabled RX Addresses
	EN_RXADDR = 0x02

	// Enable data pipe 5
	ERX_P5 = 0b00100000

	// Enable data pipe 4
	ERX_P4 = 0b00010000

	// Enable data pipe 3
	ERX_P3 = 0b00001000

	// Enable data pipe 2
	ERX_P2 = 0b00000100

	// Enable data pipe 1
	ERX_P1 = 0b00000010

	// Enable data pipe 0
	ERX_P0 = 0b00000001

	// ------------------- SETUP_AW --------------------
	// Setup of Address Widths
	// (common for all data pipes)
	SETUP_AW = 0x03

	// ------------------ SETUP_RETR -------------------
	// Setup of Automatic Retransmission
	SETUP_RETR = 0x04

	// --------------------- RF_CH ---------------------
	// RF Channel
	RF_CH = 0x05

	// ------------------- RF_SETUP --------------------
	// RF Setup Register
	RF_SETUP = 0x06

	// -------------------- STATUS ---------------------
	// Status Register (In parallel to the SPI command
	// word applied on the MOSI pin, the STATUS register
	// is shifted serially out on the MISO pin)
	STATUS = 0x07

	// ------------------ OBSERVE_TX -------------------
	// Transmit observe register
	OBSERVE_TX = 0x08

	// ---------------------- RPD ----------------------
	// Received Power Detector. This register is called
	// CD (Carrier Detect) in the nRF24L01. The name is
	// different in nRF24L01+ due to the different input
	// power level threshold for this bit.
	RPD = 0x09

	// ----------------- RX_ADDR_P0 --------------------
	// Receive address data pipe 0. 5 Bytes maximum
	// length. (LSByte is written first. Write the number of
	// bytes defined by SETUP_AW)
	RX_ADDR_P0 = 0x0A

	// ----------------- RX_ADDR_P1 --------------------
	// Receive address data pipe 1. 5 Bytes maximum
	// length. (LSByte is written first. Write the number of
	// bytes defined by SETUP_AW)
	RX_ADDR_P1 = 0x0B

	// ----------------- RX_ADDR_P2 --------------------
	// Receive address data pipe 2. Only LSB. MSBytes
	// are equal to RX_ADDR_P1
	RX_ADDR_P2 = 0x0C

	// ----------------- RX_ADDR_P3 --------------------
	// Receive address data pipe 3. Only LSB. MSBytes
	// are equal to RX_ADDR_P1
	RX_ADDR_P3 = 0x0D

	// ----------------- RX_ADDR_P4 --------------------
	// Receive address data pipe 4. Only LSB. MSBytes
	// are equal to RX_ADDR_P1
	RX_ADDR_P4 = 0x0E

	// ----------------- RX_ADDR_P5 --------------------
	// Receive address data pipe 5. Only LSB. MSBytes
	// are equal to RX_ADDR_P1
	RX_ADDR_P5 = 0x0F

	// ------------------ TX_ADDR ----------------------
	// Transmit address. Used for a PTX device only.
	// (LSByte is written first)
	// Set RX_ADDR_P0 equal to this address to handle
	// automatic acknowledge if this is a PTX device with
	// Enhanced ShockBurst™ enabled.
	TX_ADDR = 0x10

	// ----------------- RX_PW_P0 ----------------------
	// Number of bytes in RX payload in data pipe 0 (1 to
	// 	32 bytes).
	// 	0 Pipe not used
	// 	1 = 1 byte
	// 	…
	// 	32 = 32 bytes
	RX_PW_P0 = 0x11

	// ----------------- RX_PW_P1 ----------------------
	// Number of bytes in RX payload in data pipe 1 (1 to
	// 	32 bytes).
	// 	0 Pipe not used
	// 	1 = 1 byte
	// 	…
	// 	32 = 32 bytes
	RX_PW_P1 = 0x12

	// ----------------- RX_PW_P2 ----------------------
	// Number of bytes in RX payload in data pipe 2 (1 to
	// 	32 bytes).
	// 	0 Pipe not used
	// 	1 = 1 byte
	// 	…
	// 	32 = 32 bytes
	RX_PW_P2 = 0x13

	// ----------------- RX_PW_P3 ----------------------
	// Number of bytes in RX payload in data pipe 3 (1 to
	// 	32 bytes).
	// 	0 Pipe not used
	// 	1 = 1 byte
	// 	…
	// 	32 = 32 bytes
	RX_PW_P3 = 0x14

	// ----------------- RX_PW_P4 ----------------------
	// Number of bytes in RX payload in data pipe 4 (1 to
	// 	32 bytes).
	// 	0 Pipe not used
	// 	1 = 1 byte
	// 	…
	// 	32 = 32 bytes
	RX_PW_P4 = 0x15

	// ----------------- RX_PW_P5 ----------------------
	// Number of bytes in RX payload in data pipe 5 (1 to
	// 	32 bytes).
	// 	0 Pipe not used
	// 	1 = 1 byte
	// 	…
	// 	32 = 32 bytes
	RX_PW_P5 = 0x16

	// --------------- FIFO_STATUS ---------------------
	// FIFO Status Register
	FIFO_STATUS = 0x17

	// ------------------ DYNPD ------------------------
	// Enable dynamic payload length
	DYNPD = 0x1C

	// ----------------- FEATURE -----------------------
	// Feature Register
	FEATURE = 0x1D
)
