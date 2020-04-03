package co2monitor_test

import (
	"fmt"
	"testing"

	"github.com/henkman/co2monitor"
)

func TestReading(t *testing.T) {
	var cm co2monitor.CO2Monitor
	if err := cm.Open(); err != nil {
		t.Fatal(err)
	}
	defer cm.Close()
	var r co2monitor.Reading
	if err := cm.Read(&r); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("CO2 PPM: %d\nTemp °K: %f\nTemp °C: %f\nTemp °F: %f\n",
		r.CO2PPM, r.TemperatureKelvin,
		r.TemperatureCelcius(), r.TemperatureFahrenheit())
}
