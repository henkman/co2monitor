package co2monitor

import (
	"crypto/rand"
	"errors"

	"github.com/karalabe/hid"
)

type Reading struct {
	TemperatureKelvin float64
	CO2PPM            int
}

func (r Reading) TemperatureCelcius() float64 {
	return r.TemperatureKelvin - 273.15
}

func (r Reading) TemperatureFahrenheit() float64 {
	return r.TemperatureKelvin*(9.0/5.0) - 459.67
}

type CO2Monitor struct {
	key    [8]byte
	device *hid.Device
}

var (
	ErrorCouldNotGenerateKey        = errors.New("could not generate key")
	ErrorDidNotFindDevice           = errors.New("did not find device")
	ErrorCouldNotOpenDevice         = errors.New("could not open device")
	ErrorCouldNotSendFeatureRequest = errors.New("could not send feature request")
)

// Opens the CO2Monitor for reading, be sure to Close it
// Returns an error in case
// - it could not generate the decryption key
// - it did not find the device
// - it could not open the device for reading
// - it could not send the feature report to the device
func (cm *CO2Monitor) Open() error {
	const (
		vendor_id  = 0x04d9
		product_id = 0xa052
	)
	if _, err := rand.Read(cm.key[:]); err != nil {
		return ErrorCouldNotGenerateKey
	}
	devinfos := hid.Enumerate(vendor_id, product_id)
	if len(devinfos) == 0 {
		return ErrorDidNotFindDevice
	}
	dev, err := devinfos[0].Open()
	if err != nil {
		return ErrorCouldNotOpenDevice
	}
	var report [9]byte
	// report[0] = 0x00
	copy(report[1:], cm.key[:])
	if _, err := dev.SendFeatureReport(report[:]); err != nil {
		return ErrorCouldNotSendFeatureRequest
	}
	cm.device = dev
	return nil
}

var (
	CouldNotReadFromDevice = errors.New("could not read from device")
)

// Reads a fresh reading into r
// Returns an error in case of a read error
func (cm *CO2Monitor) Read(r *Reading, buf *[8]byte) error {
	var readtemp, readco2 bool
	for {
		_, err := cm.device.Read(buf[:])
		if err != nil {
			return CouldNotReadFromDevice
		}
		first := buf[2] ^ cm.key[0]
		last := buf[3] ^ cm.key[7]
		unit := ((first >> 3) | (last << 5)) - 0x84
		switch unit {
		case 0x50:
			second := buf[4] ^ cm.key[1]
			third := buf[0] ^ cm.key[2]
			high := ((second >> 3) | (first << 5)) - 0x47
			low := ((third >> 3) | (second << 5)) - 0x56
			value := uint(high)<<8 | uint(low)
			r.CO2PPM = int(value)
			if readco2 {
				return nil
			}
			readtemp = true
		case 0x42:
			second := buf[4] ^ cm.key[1]
			third := buf[0] ^ cm.key[2]
			high := ((second >> 3) | (first << 5)) - 0x47
			low := ((third >> 3) | (second << 5)) - 0x56
			value := uint(high)<<8 | uint(low)
			r.TemperatureKelvin = float64(value) / 16.0
			if readtemp {
				return nil
			}
			readco2 = true
		}
	}
}

func (cm *CO2Monitor) Close() error {
	return cm.device.Close()
}
