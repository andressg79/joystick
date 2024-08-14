package led

import "machine"

type Led struct {
	machine.Pin
} 

func NewLed(pin machine.Pin, mode machine.PinMode) Led {
	l := pin
	l.Configure(machine.PinConfig{Mode: machine.PinOutput})
	l.Low()
	return Led{l}
}

