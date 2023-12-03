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

	// Configura GPIO23 como saída e inicializa em LOW
	gpio23 := gpioreg.ByName("GPIO23")
	if gpio23 == nil {
		log.Fatal("Falha ao encontrar GPIO23")
	}
	if err := gpio23.Out(gpio.Low); err != nil {
		log.Fatal(err)
	}

	// Configura GPIO13 como entrada
	gpio13 := gpioreg.ByName("GPIO13")
	if gpio13 == nil {
		log.Fatal("Falha ao encontrar GPIO13")
	}
	if err := gpio13.In(gpio.PullDown, gpio.BothEdges); err != nil {
		log.Fatal(err)
	}

	// Escuta sinais de interrupção para encerrar o programa
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nSinal de interrupção recebido. Definindo GPIO23 para HIGH e encerrando.")
		gpio23.Out(gpio.High)
		os.Exit(0)
	}()

	// Loop para contar os pulsos
	fmt.Println("Lendo dados do sensor. Pressione Ctrl+C para sair.")
	count := 0
	startTime := time.Now()
	for {
		if gpio13.WaitForEdge(-1) { // Espera indefinidamente por uma mudança de borda
			count++
		}

		// Exemplo de cálculo de frequência (a cada segundo)
		if time.Since(startTime) >= time.Second {
			frequency := count
			count = 0
			startTime = time.Now()
			fmt.Printf("Frequência: %d Hz\n", frequency)
		}
	}
}
