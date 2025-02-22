package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"gas-monitor/internal/gas"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	infuraKey := os.Getenv("INFURA_API_KEY")
	if infuraKey == "" {
		log.Fatal("INFURA_API_KEY not found in .env")
	}

	monitor, err := gas.NewMonitor(infuraKey)
	if err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go monitor.Start()
	<-sigChan
	monitor.Stop()
} 