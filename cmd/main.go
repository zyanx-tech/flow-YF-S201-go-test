package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

func main() {
	// Inicializa a biblioteca periph
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Configura o GPIO13 para o sensor
	flowsensor := gpioreg.ByName("GPIO13")
	if flowsensor == nil {
		log.Fatal("Falha ao encontrar o GPIO13")
	}
	if err := flowsensor.In(gpio.PullDown, gpio.BothEdges); err != nil {
		log.Fatal(err)
	}

	// Configura o GPIO23 para o relé
	relay := gpioreg.ByName("GPIO23")
	if relay == nil {
		log.Fatal("Falha ao encontrar o GPIO23")
	}
	if err := relay.Out(gpio.Low); err != nil {
		log.Fatal(err)
	}

	// Preparação para capturar sinal de interrupção
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()

	// Variáveis para cálculo da vazão
	var flowFrequency int
	var lHour float64
	startTime := time.Now()

	// Loop principal
	for {
		if flowsensor.WaitForEdge(-1) { // Espera por uma mudança de borda
			flowFrequency++
		}

		// Cálculo da vazão a cada segundo
		if time.Since(startTime) >= time.Second {
			lHour = float64(flowFrequency) * 60.0 / 5.5 // Conversão para litros/hora
			flowFrequency = 0
			startTime = time.Now()
			fmt.Printf("%.2f L/hour\n", lHour) // Formata a saída para duas casas decimais

			// Lógica para controlar o relé
			if lHour > 100 {
				relay.Out(gpio.High)
			} else {
				relay.Out(gpio.Low)
			}
		}
	}
}
