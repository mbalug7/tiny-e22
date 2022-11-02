package main

import (
	"fmt"
	"log"
	"machine"
	"sync"
	"time"
	"tiny-e22/e22"
	"tiny-e22/hal"
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
	time.Sleep(time.Second * 10)
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.High()

	fmt.Println("creating HW handler")
	hw, err := NewHWHandler(machine.GP12, machine.GP13, machine.GP11, machine.UART1)
	if err != nil {
		println("could not configure HWHandler:", err)

	}
	led.Low()

	led.High()
	fmt.Println("creating a new module")
	module, err := e22.NewModule(hw, messageEvent)
	if err != nil {
		println("could not configure Module:", err)
	}
	fmt.Println("all good, setting mode")
	hw.SetMode(hal.ModeNormal)
	fmt.Println(module.GetModuleConfiguration())
	led.Low()

	cb := e22.NewConfigBuilder(module).Address(0, 1).Channel(23).AirDataRate(e22.ADR_2400).TransmissionMethod(e22.TRANSMISSION_FIXED)
	err = cb.WritePermanentConfig() // update registers on the module with the new data
	if err != nil {
		// log write error
		log.Printf("config write error: %s", err)
	} else {
		log.Println(module.GetModuleConfiguration())
	}

	machine.GP12.Low()
	machine.GP13.Low()
	hw.SetMode(hal.ModeNormal)
	// err = module.SendMessage("MESSAGE")
	// if err != nil {
	// 	fmt.Printf("faield to send: %s", err)
	// }

	for {
		led.Low()
		time.Sleep(time.Second * 60)
		//fmt.Println("msg") // UART0

		err = module.SendFixedMessage(0, 2, 23, "PING")
		if err != nil {
			fmt.Printf("failed to send: %s", err)
		}

		led.High()
		time.Sleep(time.Millisecond * 1000)
	}
}
