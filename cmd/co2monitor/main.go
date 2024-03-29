package main

import (
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"

	"github.com/henkman/co2monitor"
)

var (
	ColorHorrible  = walk.RGB(0xFF, 0x00, 0x00)
	ColorLousy     = walk.RGB(0xFD, 0x32, 0x0C)
	ColorBad       = walk.RGB(0xFD, 0xFF, 0x00)
	ColorGood      = walk.RGB(0x56, 0x82, 0x03)
	ColorExcellent = walk.RGB(0x22, 0x8B, 0x22)
	ColorPerfect   = walk.RGB(0x29, 0xAB, 0x87)
)

func main() {
	var mw *walk.MainWindow
	var co2Label *walk.Label
	var tempLabel *walk.Label

	var co2mon co2monitor.CO2Monitor
	if err := co2mon.Open(); err != nil {
		panic(err)
	}

	go func() {
		var (
			reading co2monitor.Reading
			buf     [8]byte
		)
		for {
			if err := co2mon.Read(&reading, &buf); err == nil {
				co2 := fmt.Sprint(reading.CO2PPM, " PPM")
				temp := fmt.Sprintf("%0.2f °C", reading.TemperatureCelcius())
				co2Label.SetText(co2)
				tempLabel.SetText(temp)
				mw.SetTitle(temp + " | " + co2)
				switch {
				case reading.CO2PPM < 400:
					co2Label.SetTextColor(ColorPerfect)
				case reading.CO2PPM < 600:
					co2Label.SetTextColor(ColorExcellent)
				case reading.CO2PPM < 800:
					co2Label.SetTextColor(ColorGood)
				case reading.CO2PPM < 1000:
					co2Label.SetTextColor(ColorBad)
				case reading.CO2PPM < 1200:
					co2Label.SetTextColor(ColorLousy)
				default:
					co2Label.SetTextColor(ColorHorrible)
				}
			}
		}
	}()
	defer co2mon.Close()

	const (
		HEIGHT    = 350
		FONT_SIZE = 80
	)

	var oh = HEIGHT

	var appIcon, _ = walk.NewIconFromResourceId(2)
	if err := (MainWindow{
		Icon:     appIcon,
		AssignTo: &mw,
		Title:    "CO2Monitor",
		Size: Size{
			Width:  560,
			Height: HEIGHT,
		},
		MinSize: Size{
			Width:  560,
			Height: HEIGHT,
		},
		OnSizeChanged: func() {
			s := mw.Size()
			h := s.Height
			if h == oh {
				return
			}
			rh := float32(h) / float32(HEIGHT)
			nfs := float32(FONT_SIZE) * rh
			f, _ := walk.NewFont("Segoe UI", int(nfs), walk.FontStyle(0))
			co2Label.SetFont(f)
			tempLabel.SetFont(f)
			oh = h
		},
		Background: SolidColorBrush{Color: walk.RGB(0xcc, 0xcc, 0xcc)},
		Layout:     VBox{},
		Children: []Widget{
			Label{
				AssignTo:      &tempLabel,
				Text:          "XX.XX °C",
				Font:          Font{Family: "Segoe UI", PointSize: FONT_SIZE},
				EllipsisMode:  EllipsisEnd,
				TextAlignment: AlignCenter,
			},
			Label{
				AssignTo:      &co2Label,
				Text:          "XXX PPM",
				Font:          Font{Family: "Segoe UI", PointSize: FONT_SIZE},
				TextColor:     ColorBad,
				EllipsisMode:  EllipsisEnd,
				TextAlignment: AlignCenter,
			},
		},
	}.Create()); err != nil {
		panic(err)
	}

	r := mw.Bounds()
	scrWidth := int(win.GetSystemMetrics(win.SM_CXSCREEN))
	scrHeight := int(win.GetSystemMetrics(win.SM_CYSCREEN))
	mw.SetBounds(walk.Rectangle{
		X:      int((scrWidth - r.Width) / 2),
		Y:      int((scrHeight - r.Height) / 2),
		Width:  r.Width,
		Height: r.Height,
	})
	mw.Run()
}
