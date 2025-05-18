package main

import (
	"fmt"
	"os"
	"time"

	"github.com/adamroach/sense-logger/rrd"
	"github.com/adamroach/sense-logger/sense"
)

func main() {
	client := sense.NewClient()
	// TODO -- add configuration for username and password
	err := client.Login(os.Getenv("SENSE_USER"), os.Getenv("SENSE_PASS"))
	if err != nil {
		panic(err)
	}
	// TODO -- add configuration for the output directory
	rrdWriter, err := rrd.NewWriter("out")
	if err != nil {
		panic(err)
	}
	devices, err := client.GetDevices()
	if err != nil {
		panic(err)
	}
	err = rrdWriter.UpdateDeviceNames(devices.Devices)
	if err != nil {
		panic(err)
	}
	for {
		if time.Until(client.TokenExpiry()) < time.Minute*5 {
			_, err := client.Refresh()
			if err != nil {
				panic(err)
			}
			// Close the websocket -- we'll open a new one with the new token
			client.Close()
		}
		realtimeUpdate, err := client.GetRealtimeUpdate()
		if err != nil {
			// Sometimes the websocket disconects -- we'll try refreshing the token
			_, err := client.Refresh()
			if err != nil {
				panic(err)
			}
			continue
		}
		deviceCount := len(realtimeUpdate.Payload.Devices)
		if deviceCount == 0 {
			continue
		}
		fmt.Print("\033[H\033[2J")
		fmt.Printf("%v - token expires in %v\n\n", time.Now(), time.Until(client.TokenExpiry()))
		fmt.Println(realtimeUpdate)
		err = rrdWriter.Write(realtimeUpdate)
		if err != nil {
			fmt.Printf("Error writing to RRD: %v\n", err)
			panic(err)
		}

	}
}
