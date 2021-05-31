package main

import (
	"time"
)

func main() {

	s := initialize()
	sets := setTargets()
	s.SetStyle(defStyle)
	go checkExit(s)
	for {
		rdata := getReadings()
		displayOnMatrix(s, rdata, sets)
		displayReadings(s, 35, 2, rdata)
		ctrl := setControls(sets, rdata)
		displayTargets(s, 35, 4, sets)
		displayControls(s, 35, 5, ctrl)
		s.Show()
		time.Sleep(2 * time.Second)
	}

}
