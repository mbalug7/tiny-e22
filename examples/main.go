package main

import (
	"fmt"
	"log"
	"machine"
	"sync"
	"time"

	"github.com/mbalug7/go-ebyte-lora/pkg/e22"
	"github.com/mbalug7/go-ebyte-lora/pkg/hal"
	"github.com/mbalug7/tiny-e22/pico"
)

var mu sync.Mutex

func messageEvent(msg e22.Message, err error) {
	if err != nil {
		log.Printf("message event error: %s", err)
		return
	}
	log.Printf("DATA: %s\n", string(msg.Payload))
}

func main() {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.High()

	fmt.Println("creating HW handler")
	hw, err := pico.NewHWHandler(machine.GP12, machine.GP13, machine.GP11, machine.UART1)
	if err != nil {
		println("could not configure HWHandler:", err)

	}
	fmt.Println("creating new module")
	module, err := e22.NewModule(hw, messageEvent)
	if err != nil {
		println("could not configure Module:", err)
	}
	fmt.Println("all good, setting mode")
	hw.SetMode(hal.ModeNormal)

	fmt.Println("Building module conf")
	cb := e22.NewConfigBuilder(module).Address(0, 1).Channel(23).AirDataRate(e22.ADR_2400).TransmissionMethod(e22.TRANSMISSION_FIXED)
	err = cb.WritePermanentConfig() // update registers on the module with the new data
	if err != nil {
		// log write error
		log.Printf("config write error: %s", err)
	}

	fmt.Println(module.GetModuleConfiguration())

	go func() {
		for {
			led.Low()
			time.Sleep(time.Second * 2)

			led.High()
			time.Sleep(time.Second * 2)
		}
	}()
	for {
		time.Sleep(time.Second * 5)
		log.Println("Sending Ping")
		go func() {
			err = module.SendFixedMessage(0, 2, 23, "PING")
			if err != nil {
				fmt.Printf("failed to send: %s", err)
			}
		}()
		time.Sleep(time.Millisecond * 2000)
	}
}
