package hardware

import (
	"machine"

	"tinygo.org/x/drivers/ssd1306"
)

func NewDisplay(ic2 *machine.I2C) ssd1306.Device {
	println("Display init")
	dev := ssd1306.NewI2C(ic2)
	dev.Configure(ssd1306.Config{Width: 128, Height: 32, Address: ssd1306.Address_128_32, VccState: ssd1306.SWITCHCAPVCC})
	dev.ClearDisplay()
	return dev
}
