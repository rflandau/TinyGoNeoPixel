# TinyGo NeoPixel LEDs
Simple wrapper for TinyGo's built-in ws2812 functionality, applied to controlling LEDs.

## Installation

Fetch the package with:

```bash
go get github.com/rflandau/TinyGoNeoPixel
```

Import it with:
```go
npxl "github.com/rflandau/TinyGoNeoPixel"
```


## Example Usage
```go
var count uint = 7
var neo *neopixel.NeoPixels = npxl.New(machine.NEOPIXELS, count)
// if you want the subroutines to automtaically enact changes, rather than calling
neo.AutoWrite = true

neo.SetLED(2, color.RGBA{})
neo.Flush() // if !neo.AutoWrite
```

