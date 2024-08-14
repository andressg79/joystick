package hardware

import (
	"machine"
)

func NewIC2(ic2 *machine.I2C, scl, sda machine.Pin) *machine.I2C {
	println("init I2C")
	ic2.Configure(machine.I2CConfig{
		Frequency: 400000,
		SCL:       scl,
		SDA:       sda,
	})
	return ic2
}
