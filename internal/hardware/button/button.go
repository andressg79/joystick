package button

import (
	"joystick/internal/pkg/boolean"
	"machine"
)

type Button struct {
	machine.Pin
}

func NewButton(pin machine.Pin, mode machine.PinMode) Button {
	b := pin
	b.Configure(machine.PinConfig{Mode: mode})
	return Button{b}
}

func (b Button) IsON() uint8 {
	return boolean.ToByte(b.Get())
}
