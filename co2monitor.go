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
func (cm *CO2Monitor) Read(r *Reading) error {
	var buf [8]byte
	var readtemp, readco2 bool
	for {
		n, err := cm.device.Read(buf[:])
		if err != nil {
			return CouldNotReadFromDevice
		}
		decrypt(cm.key[:], buf[:n])
		const (
			Temp = 0x42
			CO2  = 0x50
		)
		switch buf[0] {
		case CO2:
			value := uint(buf[1])<<8 | uint(buf[2])
			r.CO2PPM = int(value)
			if readco2 {
				return nil
			}
			readtemp = true
		case Temp:
			value := uint(buf[1])<<8 | uint(buf[2])
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

func decrypt(key, data []byte) {
	var tmp [8]byte
	for i, v := range []int{2, 4, 0, 7, 1, 6, 5, 3} {
		tmp[v] = data[i] ^ key[v]
	}
	second := []byte{0x84, 0x47, 0x56, 0xd6, 0x07, 0x93, 0x93, 0x56}
	data[0] = (tmp[0] >> 3) | (tmp[7] << 5) - second[0] + 0xFF + 1
	for i := 1; i < 8; i++ {
		data[i] = ((tmp[i] >> 3) | (tmp[i-1] << 5)) - second[i] + 0xFF + 1
	}
}
