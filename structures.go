package main

import (
	"time"
)

type readings struct {
	timeStamp   time.Time
	temperature float64
	humidity    float64
	pressure    float64
}

type setPoints struct {
	temperature float64
	humidity    float64
}

type controls struct {
	heater     int
	humidifier int
}
type pixel struct {
	red   uint8
	green uint8
	blue  uint8
}
