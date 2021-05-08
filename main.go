package main

import "time"

func main() {
	initialize()
	sets := setTargets()
	for {
		rdata := getReadings()
		displayReadings(rdata)
		ctrl := setControls(sets, rdata)
		displayTargets(sets)
		displayControls(ctrl)
		time.Sleep(2 * time.Second)
	}

}
