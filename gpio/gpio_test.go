package gpio

import "testing"

func TestGpioPath(t *testing.T) {
	p := gpioPin(3)
	if p.path() != "/sys/class/gpio/gpio3" {
		t.Fail()
	}

	p = gpioPin(4)
	if p.path() != "/sys/class/gpio/gpio4" {
		t.Fail()
	}
}
