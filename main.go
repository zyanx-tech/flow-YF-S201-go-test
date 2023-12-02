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
	flowSensorPinName   = "GPIO13"
	valveControlPinName = "GPIO23"
	flowRate            = 4.5 // em mL por segundo, ajuste conforme as especificações do seu sensor
)

func main() {
	fmt.Println("Start...")

	// Inicializa periph.io
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	flowSensorPin := gpioreg.ByName(flowSensorPinName)
	if flowSensorPin == nil {
		log.Fatalf("Failed to find %s", flowSensorPinName)
	}

	// Configura o pino do sensor de fluxo para entrada
	if err := flowSensorPin.In(gpio.PullDown, gpio.BothEdges); err != nil {
		log.Fatal(err)
	}

	valveControlPin := gpioreg.ByName(valveControlPinName)
	if valveControlPin == nil {
		log.Fatalf("Failed to find %s", valveControlPinName)
	}

	// Configura a GPIO da válvula como saída e abre a válvula
	if err := valveControlPin.Out(gpio.High); err != nil {
		log.Fatal(err)
	}

	// Configura a GPIO da válvula como saída e abre a válvula
	if err := valveControlPin.Out(gpio.Low); err != nil {
		log.Fatal(err)
	}

	// Monitora o sensor de fluxo
	go monitorFlowSensor(flowSensorPin, valveControlPin)

	// Mantém o programa rodando
	select {}
}

func monitorFlowSensor(flowSensorPin, valveControlPin gpio.PinIO) {
	pulseCount := 0
	startTime := time.Now()
	fmt.Println("FlowSensor")
	for {
		if flowSensorPin.WaitForEdge(-1) {
			pulseCount++
			elapsedTime := time.Since(startTime).Seconds()
			volume := float64(pulseCount) / flowRate

			fmt.Printf("Volume medido: %.2f mL\n", volume)

			if volume >= 200 {
				// Fecha a válvula e para o monitoramento
				valveControlPin.Out(gpio.High)
				fmt.Println("Limite de 200 mL atingido, válvula fechada")
				return
			}

			// Impede a verificação constante do sensor
			time.Sleep(time.Millisecond * 100)

			if elapsedTime > 30 { // Tempo máximo de operação (1/2 minuto, ajuste conforme necessário)
				valveControlPin.Out(gpio.Low)
				fmt.Println("Tempo máximo de operação atingido, válvula fechada")
				return
			}
		}
	}
}
