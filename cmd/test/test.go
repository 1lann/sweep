package main

import (
	"fmt"
	"log"

	"github.com/1lann/sweep"
)

func main() {
	dev, err := sweep.NewDevice("/dev/cu.usbserial-DO0088ZE")
	if err != nil {
		panic(err)
	}
	defer dev.Close()

	log.Println("Welcome, starting up...")

	log.Println(dev.Restart())
	log.Println("wait:", dev.WaitUntilMotorReady())
	log.Println("set speed:", dev.SetMotorSpeed(3))
	log.Println("wait:", dev.WaitUntilMotorReady())
	log.Println("set rate:", dev.SetSampleRate(sweep.Rate1000))

	log.Println("Almost there...")
	log.Println("wait:", dev.WaitUntilMotorReady())

	log.Println("Done!")

	rd, err := dev.StartScan()
	log.Println(err)

	for scan := range rd {
		min := 100000000
		var str byte
		for _, s := range scan {
			if s.Distance < min {
				min = s.Distance
				str = s.SignalStrength
			}
			// fmt.Println(s.Distance)
		}

		fmt.Println("min:", min, "str:", str)
	}
}
