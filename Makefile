tty=/dev/ttyACM0
target=pico

ll:
	ls -la /dev/ttyU* /dev/ttyA*

monitor: 
	tinygo monitor -baudrate 9600 -port $(tty)

blink:
	tinygo flash -port $(tty) -target $(target) cmd/blink/main.go

flash:
	tinygo flash -port $(tty) -target $(target) $(gofile)

joystick:
	tinygo flash -port $(tty) -target $(target) cmd/joystick/main.go && tinygo monitor -baudrate 9600 -port $(tty)

e-joystick:
	tinygo flash -port $(tty) -target $(target) cmd/experiments/joystick/main.go && tinygo monitor -baudrate 9600 -port $(tty)

e-joystick-d:
	tinygo flash -port $(tty) -target $(target) cmd/experiments/joystick-double/main.go && tinygo monitor -baudrate 9600 -port $(tty)

display-simple:
	tinygo flash -port $(tty) -target $(target) cmd/display-simple/main.go

tricycle:
	tinygo flash -port $(tty) -target $(target) cmd/tricycle/main.go

ne-joystick:
	tinygo flash -port $(tty) -target $(target) cmd/nano/joystick/main.go && tinygo monitor -baudrate 9600 -port $(tty)
