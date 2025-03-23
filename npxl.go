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

type NeoPixels struct {
	ledCount  uint
	colors    []color.RGBA
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
		colors:   make([]color.RGBA, count),
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
func (np *NeoPixels) SetLED(index uint, color color.RGBA) error {
	if index >= np.ledCount {
		return ErrOutOfRange{}
	}

	np.colors[index] = color

	np.aw()

	return nil
}

// Sets all LEDs to the colors associated with their index.
// Only updates to the lower of the two lengths (LED count and len(colors)).
func (np *NeoPixels) SetAllLEDs(colors []color.RGBA) {
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

// Apply the current state of the LEDs to the board itself.
// This is called automatically if .AutoWrite is set.
// Use .Flush() if you want to ensure multiple LED changes occur simultaneously, rather than incrementally.
func (np *NeoPixels) Flush() error {
	return np.conn.WriteColors(np.colors)
}

//region pre-defined patterns

// sets LEDs alternatively to red and green
func (np *NeoPixels) StaticChristmas() {
	for i := range np.ledCount {
		switch i % 2 {
		case 0:
			np.colors[i] = Green
		case 1:
			np.colors[i] = Red
		}
	}

	np.aw()
}

// sets all LEDs to green
func (np *NeoPixels) StaticGreen() {
	for i := range np.ledCount {
		np.colors[i] = Green
	}

	np.aw()
}

// sets LEDs to the colors of the rainbow
func (np *NeoPixels) StaticRainbow() {
	for i := range np.ledCount {
		switch i % 7 {
		case 0:
			np.colors[i] = Red
		case 1:
			np.colors[i] = Orange
		case 2:
			np.colors[i] = Yellow
		case 3:
			np.colors[i] = Green
		case 4:
			np.colors[i] = Blue
		case 5:
			np.colors[i] = Navy
		case 6:
			np.colors[i] = Purple
		}
	}

	np.aw()
}

//endregion
