package blync

import (
	"fmt"
	"github.com/boombuler/hid"
	"os"
)

const (
	blyncVendorId  = 0x0E53
	blyncProductId = 0x2517
)

const (
	BlinkOff    = 0x00
	BlinkFast   = 0x46
	BlinkMedium = 0x64
	BlinkSlow   = 0x96
)

var Red = [3]byte{0xFF, 0x00, 0x00}
var Green = [3]byte{0x00, 0xFF, 0x00}
var Blue = [3]byte{0x00, 0x00, 0xFF}

type BlyncLight struct {
	devices []hid.Device
	bytes   []byte
}

func NewBlyncLight() (blync BlyncLight) {
	blync.devices = findDevices()
	blync.bytes = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x02, 0xFF}
	return
}

func findDevices() []hid.Device {
	devices := []hid.Device{}
	deviceInfos := hid.Devices()
	for {
		info, more := <-deviceInfos
		if more {
			device, error := info.Open()
			if error != nil {
				fmt.Println(error)
			}
			if !isBlyncDevice(*info) {
				fmt.Printf("%s %s is not a BlyncLight device.\n", info.Manufacturer, info.Product)
			} else {
				devices = append(devices, device)
				fmt.Printf("%s %s is a BlyncLight device.\n", info.Manufacturer, info.Product)
			}
		} else {
			break
		}
	}
	if len(devices) == 0 {
		fmt.Println("No BlyncLights found.")
		os.Exit(1)
	}
	return devices
}

func isBlyncDevice(deviceInfo hid.DeviceInfo) bool {
	// TODO from forums: "Blync creates 2 HID devices and the only way to find out the right device is the MaxFeatureReportLength = 0"
	if deviceInfo.VendorId == blyncVendorId && deviceInfo.ProductId == blyncProductId && deviceInfo.FeatureReportLength == 0 {
		return true
	}
	return false
}

func (b BlyncLight) sendFeatureReport() {
	for _, device := range b.devices {
		error := device.Write(b.bytes)
		if error != nil {
			fmt.Println(error)
		}
	}
}

func (b BlyncLight) Close() {
	b.Reset()
	for _, device := range b.devices {
		device.Close()
	}
}

func (b BlyncLight) Reset() {
	b.bytes = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x02, 0xFF}
	b.sendFeatureReport()
}

// color[0] = r
// color[1] = g
// color[2] = b
func (b BlyncLight) SetColor(color [3]byte) {
	b.bytes[1] = color[0]
	b.bytes[2] = color[2] // They reverse g and b
	b.bytes[3] = color[1]
	b.sendFeatureReport()
}

func (b BlyncLight) SetBlinkRate(rate byte) {
	b.bytes[4] = rate
	b.sendFeatureReport()
}

//16-30 play a tune single time
//49-59 plays never ending versions of the tunes
func (b BlyncLight) Play(mp3 byte) {
	b.bytes[5] = mp3
	b.sendFeatureReport()
}

func (b BlyncLight) StopPlay() {
	b.bytes[5] = 0x00
	b.sendFeatureReport()
}
