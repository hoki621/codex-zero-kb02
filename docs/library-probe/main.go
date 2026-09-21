// Compile-only API probe. Do not flash: this is not product firmware.
package main

import (
	"image/color"
	"machine"
	"machine/usb/hid/mouse"
	"time"

	keyboard "github.com/sago35/tinygo-keyboard"
	pio "github.com/tinygo-org/pio/rp2-pio"
	"github.com/tinygo-org/pio/rp2-pio/piolib"
	"tinygo.org/x/drivers/encoders"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
)

func main() {
	matrix := keyboard.New().AddMatrixKeyboard(
		[]machine.Pin{machine.GPIO5, machine.GPIO6, machine.GPIO7, machine.GPIO8},
		[]machine.Pin{machine.GPIO9, machine.GPIO10, machine.GPIO11}, nil)
	encoder := encoders.NewQuadratureViaInterrupt(machine.GPIO3, machine.GPIO4)
	if err := encoder.Configure(encoders.QuadratureConfig{Precision: 4}); err != nil {
		panic(err)
	}
	machine.InitADC()
	adc := machine.ADC{Pin: machine.GPIO29}
	adc.Configure(machine.ADCConfig{})
	if err := machine.I2C0.Configure(machine.I2CConfig{SDA: machine.GPIO12, SCL: machine.GPIO13, Frequency: 400000}); err != nil {
		panic(err)
	}
	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{Width: 128, Height: 64, Address: 0x3c})
	tinyfont.WriteLine(display, &tinyfont.TomThumb, 0, 8, "WIBDUE", color.RGBA{255, 255, 255, 255})
	tinydraw.Rectangle(display, 0, 0, 64, 21, color.RGBA{255, 255, 255, 255})
	if err := display.Display(); err != nil {
		panic(err)
	}
	sm, err := pio.PIO0.ClaimStateMachine()
	if err != nil {
		panic(err)
	}
	leds, err := piolib.NewWS2812B(sm, machine.GPIO1)
	if err != nil {
		panic(err)
	}
	if err := leds.WriteRaw(make([]uint32, 12)); err != nil {
		panic(err)
	}
	for {
		println("PROBE", len(matrix.Get()), encoder.Position(), adc.Get())
		mouse.Port().Move(0, 0)
		time.Sleep(time.Millisecond)
	}
}
