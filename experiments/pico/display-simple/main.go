package main

import (
	//"image/color"
	"image/color"
	"machine"
	"time"

	//"time"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

func main() {

	var led = machine.GPIO13
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	time.Sleep(time.Millisecond * 100) // Please wait some time after turning on the device to properly initialize the display

	println("init I2C0")
	machine.I2C1.Configure(machine.I2CConfig{
		Frequency: 400000,
		//SCL:       machine.I2C1_SCL_PIN,
		//SDA:       machine.I2C1_SDA_PIN,
		SCL: machine.GPIO3,
		SDA: machine.GPIO2,
	})

	// Display
	println("Display init")
	dev := ssd1306.NewI2C(machine.I2C1)
	dev.Configure(ssd1306.Config{Width: 128, Height: 32, Address: ssd1306.Address_128_32, VccState: ssd1306.SWITCHCAPVCC})
	dev.ClearDisplay()

	tinyfont.WriteLine(&dev, &freemono.Regular9pt7b, 1, 0x09, "Hola Noel :)", color.RGBA{255, 255, 255, 255})

	dev.Display()

	/*bf := make([]byte, 128*32/8)
	for i := 0; i < 128*32/8; i++ {
		bf[i] = 0x00
	}*/

	for {
		/*dev.ClearBuffer()

		bf := make([]byte, 128*32/8)
		for i := 0; i < 128*32/8; i++ {
			bf[i] = 0xFF
		}
		dev.Display()
		//Display(dev, bf)
		*/
		led.High()
		time.Sleep(time.Millisecond * 500)
		led.Low()
		time.Sleep(time.Millisecond * 500)
		print(".")
	}
}

func Display(dev ssd1306.Device, bf []byte) {
	dev.SetBuffer(bf)

	dev.Command(ssd1306.COLUMNADDR)
	dev.Command(0)
	dev.Command(uint8(127))
	dev.Command(ssd1306.PAGEADDR)
	dev.Command(0)
	dev.Command(uint8(32/8) - 1)

	err := machine.I2C0.WriteRegister(0x3C, 0x40, bf)
	if err != nil {
		println("Failed to display")
		println(err.Error())
	}
}
