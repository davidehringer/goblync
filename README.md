# GoBlync

goblync is a Go library for interacting with [BlyncLight](http://www.embrava.com/) devices.  This library is compatible with BlyncLight
v3.

## Building

Based on the dependency on Go [HID](https://github.com/boombuler/hid), the following operating systems are supported targets (as used by $GOOS environment variable)

* darwin (uses native IOKit framework)
* linux (uses libusb 1.0+)
* windows (uses native Windows HID library)

```
go get github.com/boombuler/hid
go build
```

For building on Windows, see the potential quirks listed [here](https://github.com/boombuler/hid).

## Usage

```
package main

import (
	"github.com/davidehringer/goblync"
	"time"
)

func main() {

	light := blync.NewBlyncLight()
	time.Sleep(time.Second * 2)
	light.SetColor(blync.Red)
	light.Play(52)
	time.Sleep(time.Second * 5)
	light.SetColor(blync.Blue)
	time.Sleep(time.Second * 5)
	light.StopPlay()
	light.Reset()

	for i := 0; i < 256; i++ {
		light.SetColor([3]byte{byte(i), 255 - byte(i), 0x00})
		time.Sleep(13 * time.Millisecond)
	}
	light.SetBlinkRate(blync.BlinkMedium)
	time.Sleep(time.Second * 5)
	light.Close()
}

```