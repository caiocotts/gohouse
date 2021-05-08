package main

import "math/rand"

func getHumidity() float64 {
	return float64(rand.Intn(100-0) + 0)
}

func getTemperature() float64 {
	return float64(rand.Intn(50 - -10) + -10)
}

func getPressure() float64 {
	return float64(rand.Intn(1016-975) + 975)
}
