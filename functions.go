package main

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"time"
)

const ON = 1
const OFF = 0

func getSerial() string {
	serial, err := exec.Command(
		"bash",
		"-c",
		"grep Serial /proc/cpuinfo | grep -Po '[\\d]+'",
	).Output()

	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(serial))

}

func displayReadings(rdata readings) {
	currentTime := time.Now()

	fmt.Printf(
		"Unit: %s %s \nReadings\tT: %5.1fC\tH: %5.1f%%\tP: %6.1fmb\n",
		getSerial(),
		currentTime.Format(time.RFC1123),
		rdata.temperature,
		rdata.humidity,
		rdata.pressure,
	)
}

func initialize() {
	rand.Seed(time.Now().UnixNano())
}

func getReadings() readings {
	return readings{
		timeStamp:   time.Now(),
		temperature: getTemperature(),
		humidity:    getHumidity(),
		pressure:    getPressure(),
	}
}

func setTargets() setPoints {
	return setPoints{
		temperature: 25.0,
		humidity:    55.0,
	}
}
func setControls(target setPoints, rdata readings) controls {
	state := controls{}

	if rdata.temperature < target.temperature {
		state.heater = ON
	} else {
		state.heater = OFF
	}

	if rdata.humidity < target.humidity {
		state.humidifier = ON
	} else {
		state.humidifier = OFF
	}

	return state
}

func displayTargets(spts setPoints) {
	fmt.Printf("Targets\t\tT: %5.1fC\tH: %5.1f%%\t\n", spts.temperature, spts.humidity)
}

func displayControls(ctrl controls) {
	fmt.Printf("Controls\tHeater: %d\tHumidifier: %d\n\n", ctrl.heater, ctrl.humidifier)
}
