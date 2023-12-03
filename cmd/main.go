package main

import (
	"fmt"
	"log"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

const (
	flowSensorPinName   = "GPIO13" // Pino do sensor de fluxo
	valveControlPinName = "GPIO23" // Pino de controle da válvula solenoide
	runTime             = 60 * time.Second
	sampleTime          = 1 * time.Second
)

func main() {
	fmt.Println("Iniciando...")

	// Inicializa periph.io
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Configura o pino do sensor de fluxo
	flowSensorPin := gpioreg.ByName(flowSensorPinName)
	if flowSensorPin == nil {
		log.Fatalf("Falha ao encontrar %s", flowSensorPinName)
	}
	if err := flowSensorPin.In(gpio.PullUp, gpio.FallingEdge); err != nil {
		log.Fatal(err)
	}

	// Configura o pino da válvula solenoide
	valveControlPin := gpioreg.ByName(valveControlPinName)
	if valveControlPin == nil {
		log.Fatalf("Falha ao encontrar %s", valveControlPinName)
	}
	if err := valveControlPin.Out(gpio.Low); err != nil {
		log.Fatal(err)
	}

	// Monitora o sensor de fluxo
	go monitorFlowSensor(flowSensorPin, valveControlPin)

	// Executa por um tempo definido
	time.Sleep(runTime)
	fmt.Println("Execução finalizada.")

	// Desliga a válvula solenoide ao finalizar
	valveControlPin.Out(gpio.High)
}

func monitorFlowSensor(flowSensorPin, valveControlPin gpio.PinIO) {
	pulseCount := 0
	startTime := time.Now()

	for {
		if flowSensorPin.WaitForEdge(-1) {
			pulseCount++
		}

		if time.Since(startTime) > sampleTime {
			fmt.Printf("Contagem de pulsos: %d\n", pulseCount)
			pulseCount = 0
			startTime = time.Now()
		}
	}
}
