/* This package manages the NeoPixel LEDs (and maintains a local data representation) via ws2812 */
package npxl

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers/ws2812"
)

type ErrOutOfRange struct{}

func (ErrOutOfRange) Error() string {
	return "led index must be 0 < x < LED count-1"
}

type led = color.RGBA

type NeoPixels struct {
	ledCount  uint
	colors    []led
	conn      ws2812.Device
	AutoWrite bool // Write changes directly to the LEDs. If false, you must call .Flush()
}

// Returns a new NeoPixel struct set to read off the given pin.
// LED indices are zero-indexed, obviously.
// Initializes all LEDs to OFF.
func New(pin machine.Pin, count uint) *NeoPixels {
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	np := &NeoPixels{
		ledCount: count,
		colors:   make([]led, count),
		conn:     ws2812.NewWS2812(pin),
	}

	np.Off()
	np.conn.WriteColors(np.colors)

	return np
}

// Sets all LEDs to off.
func (np *NeoPixels) Off() {
	for i := range np.colors {
		np.colors[i].R = 0
		np.colors[i].G = 0
		np.colors[i].B = 0
		np.colors[i].A = 0
	}

	np.aw()
}

// Set the color of a single LED
func (np *NeoPixels) SetLED(index uint, color led) error {
	if index >= np.ledCount {
		return ErrOutOfRange{}
	}

	np.colors[index] = color

	np.aw()

	return nil
}

// Sets all LEDs to the colors associated with their index.
// Only updates to the lower of the two lengths (LED count and len(colors)).
func (np *NeoPixels) SetAllLEDs(colors []led) {
	maxIndex := max(uint(len(colors)), np.ledCount)
	for i := range maxIndex {
		np.colors[i] = colors[i]
	}

	np.aw()
}

// internal autowrite helper function to automatically write the changes iff the bool is set
func (np *NeoPixels) aw() {
	if np.AutoWrite {
		np.conn.WriteColors(np.colors)
	}

}
