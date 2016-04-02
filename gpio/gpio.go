package gpio

import (
	"fmt"
	"os"
	"path"
)

const (
	GPIO_DIR = "/sys/class/gpio"

	GPIO_OUT = "out"
	GPIO_IN  = "in"
)

type Direction string
type Pin int

func openWFile(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_WRONLY, 0666)
}

func (p Pin) path() string {
	return path.Join(GPIO_DIR, fmt.Sprintf("gpio%d", p))
}

func (p Pin) SetDirection(direction Direction) error {
	return writef(path.Join(p.path(), "direction"), "%s\n", direction)
}

func (p Pin) SetValue(v bool) error {
	numVal := 0
	if v {
		numVal = 1
	}

	return writef(path.Join(p.path(), "value"), "%d\n", numVal)
}

func writef(path string, format string, a ...interface{}) error {
	f, err := openWFile(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, format, a...); err != nil {
		return err
	}

	return nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func exportPin(pin int) error {
	f, err := openWFile(path.Join(GPIO_DIR, "export"))
	defer f.Close()

	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(f, "%d\n", pin); err != nil {
		return err
	}

	return nil
}

func BuildGpioPin(pin int) (Pin, error) {
	PinDir := path.Join(GPIO_DIR, fmt.Sprintf("gpio%d", pin))
	if ok, err := exists(PinDir); err != nil {
		return 0, err
	} else if !ok {
		if err := exportPin(pin); err != nil {
			return 0, err
		}

		if ok, err := exists(PinDir); !ok {
			return 0, err
		}
	}

	return Pin(pin), nil
}
