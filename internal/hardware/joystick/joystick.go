package joystick

import (
	"joystick/internal/pkg/boolean"
	"machine"
	"strconv"
)

type Joystick struct {
	ID string
	HW *Hardware
}

func NewJoystick(id string, hw *Hardware) *Joystick {
	return &Joystick{
		ID: id,
		HW: hw,
	}
}

type Hardware struct {
	ID string
	X  machine.ADC
	Y  machine.ADC
	Sw machine.Pin
}

func NewkHardware(x machine.ADC, y machine.ADC, sw machine.Pin) *Hardware {
	return &Hardware{
		X:  x,
		Y:  y,
		Sw: sw,
	}
}

func (j *Hardware) Init() {
	j.X.Configure(machine.ADCConfig{})
	j.Y.Configure(machine.ADCConfig{})
	j.Sw.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
}

func (j *Hardware) Read() (uint16, uint16, bool) {
	return j.X.Get(), j.Y.Get(), j.Sw.Get()
}

func (j *Joystick) Init() {
	j.HW.Init()
}

func (j *Joystick) Read() (uint16, uint16, bool) {
	return j.HW.Read()
}

func (j *Joystick) ReadUnit8() (uint8, uint8, uint8, uint8, uint8) {
	x, y, s := j.Read()
	return uint8(x / 255), uint8(x % 255), uint8(y / 255), uint8(y % 255), boolean.ToByte(s)
}

func (j *Joystick) String() string {
	x, y, s := j.Read()
	return "[" + j.ID +
		" {X:" + strconv.FormatUint(uint64(x), 10) +
		" Y:" + strconv.FormatUint(uint64(y), 10) +
		" SW:" + strconv.FormatBool(s) + "}]"
}
