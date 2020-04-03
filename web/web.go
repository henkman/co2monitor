package web

import (
	"encoding/json"
	"net/http"

	"github.com/henkman/co2monitor"
)

// http.Handler that reads your CO2Monitor on request
// and writes it out as json
//
// can specify temperature unit with url param u
// c -> celcius
// f -> fahrenheit
// k or empty -> kelvin
//
// e.g. http://urltohandler?u=c for celcius
//
type JsonCO2Monitor struct {
	co2monitor.CO2Monitor
	reading co2monitor.Reading
}

// Run this in a go routine
func (jcm *JsonCO2Monitor) Run() {
	for {
		var reading co2monitor.Reading
		if err := jcm.Read(&reading); err == nil {
			jcm.reading = reading
		}
	}
}

func (jcm *JsonCO2Monitor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reading := jcm.reading
	var temp float64
	switch r.URL.Query().Get("u") {
	case "c":
		temp = reading.TemperatureCelcius()
	case "f":
		temp = reading.TemperatureFahrenheit()
	case "k":
		fallthrough
	default:
		temp = reading.TemperatureKelvin
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Temp float64 `json:"t"`
		CO2  int     `json:"c"`
	}{
		Temp: temp,
		CO2:  reading.CO2PPM,
	})
}
