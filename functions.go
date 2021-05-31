package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
)

const (
	on           = 1
	off          = 0
	lowestTemp   = -10.0
	highestTemp  = 50.0
	lowestHumid  = 0
	highestHumid = 100
	lowestPress  = 975
	highestPress = 1016
	tempBar      = 7
	humidBar     = 5
	pressBar     = 3
)

var (
	targetStyle = tcell.StyleDefault.
			Foreground(tcell.ColorBlack).
			Background(tcell.ColorRed)

	onStyle = tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorGreen)

	offStyle = tcell.StyleDefault.
			Foreground(tcell.ColorBlack).
			Background(tcell.ColorBlack)

	defStyle = tcell.StyleDefault.
			Background(tcell.ColorBlack).
			Foreground(tcell.ColorWhite)

	green = pixel{
		red:   0x00,
		green: 0xFF,
		blue:  0x00,
	}

	red = pixel{
		red:   0xFF,
		green: 0x00,
		blue:  0x00,
	}
)

func getId() string {
	var id string

	if _, err := os.Stat("./id"); os.IsNotExist(err) {
		f, err := os.Create("id")
		if err != nil {
			println("Could not create file \"id\" ")
			log.Fatal(err)
		}
		id = uuid.NewString()
		_, err = f.WriteString(id)

		if err != nil {
			println("Could not store id ")
			log.Fatal(err)

		}

	} else {
		f, err := ioutil.ReadFile("./id")
		if err != nil {
			println("Could not read from file \"id\"")
			log.Fatal(err)
		}
		id = string(f)

	}
	return id
}

func initialize() tcell.Screen {
	encoding.Register()
	s, e := tcell.NewScreen()
	if e != nil {
		log.Fatal(e)
	}

	if e := s.Init(); e != nil {
		log.Fatal(e)
	}

	rand.Seed(time.Now().UnixNano())
	return s
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
		state.heater = on
	} else {
		state.heater = off
	}

	if rdata.humidity < target.humidity {
		state.humidifier = on
	} else {
		state.humidifier = off
	}

	return state
}

func displayReadings(s tcell.Screen, x, y int, rdata readings) {
	currentTime := time.Now()

	printmv(s, x, y, defStyle, fmt.Sprintf(
		"Unit: %s %s",
		getId(),
		currentTime.Format(time.RFC1123),
	))

	printmv(s, x, y+1, defStyle, fmt.Sprintf(
		"Readings     T: %5.1fC       H: %5.1f%%       P: %6.1fmb",
		rdata.temperature,
		rdata.humidity,
		rdata.pressure,
	))

}

func displayTargets(s tcell.Screen, x, y int, spts setPoints) {
	printmv(s, x, y, defStyle, fmt.Sprintf(
		"Targets      T: %5.1fC       H: %5.1f%%",
		spts.temperature,
		spts.humidity,
	))
}

func displayControls(s tcell.Screen, x, y int, ctrl controls) {
	printmv(s, x, y, defStyle, fmt.Sprintf(
		"Controls     Heater: %d       Humidifier: %d\n\n",
		ctrl.heater,
		ctrl.humidifier,
	))
}

func checkExit(s tcell.Screen) {
	defer func() {
		s.Fini()
		os.Exit(1)
	}()
	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			//displayHelloWorld(s)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyRune && ev.Rune() == 'q' || ev.Rune() == 'Q' {
				return
			}
		}
	}
}

func displayOnMatrix(s tcell.Screen, values readings, targets setPoints) {
	scaledValue := scale(values.temperature, lowestTemp, highestTemp)
	setVerticalBar(s, tempBar, green, scaledValue)

	scaledValue = scale(values.humidity, lowestHumid, highestHumid)
	setVerticalBar(s, humidBar, green, scaledValue)

	scaledValue = scale(values.pressure, lowestPress, highestPress)
	setVerticalBar(s, pressBar, green, scaledValue)

	scaledValue = scale(targets.temperature, lowestTemp, highestTemp)
	writePixel(s, tempBar, scaledValue, red)

	scaledValue = scale(targets.humidity, lowestHumid, highestHumid)
	writePixel(s, humidBar, scaledValue, red)

}

func setVerticalBar(s tcell.Screen, bar int, px pixel, value int) {
	if value > 7 {
		value = 7
	}
	if bar >= 0 && bar < 8 && value < 8 {
		for i := 0; i <= value; i++ {
			writePixel(s, bar, i, px)
		}
		for i := value + 1; i < 8; i++ {
			writePixel(s, bar, i, pixel{})
		}
	}
}

func writePixel(s tcell.Screen, x, y int, pixelColour pixel) {
	var (
		simulatedXCoordinate = (7 - x) * 4
		simulatedYCoordinate = (7 - y) * 2
	)

	switch {
	case pixelColour.red != 0:
		printmv(s, simulatedXCoordinate, simulatedYCoordinate, targetStyle, "  ")
	case pixelColour.green != 0:
		printmv(s, simulatedXCoordinate, simulatedYCoordinate, onStyle, "  ")
	default:
		printmv(s, simulatedXCoordinate, simulatedYCoordinate, offStyle, "  ")
	}

}

func scale(value, min, max float64) int {
	return int((8.0 * (((value - min) / (max - min)) + 0.05)) - 1.0)
}
