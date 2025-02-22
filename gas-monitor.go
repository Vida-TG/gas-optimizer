package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	infuraKey := os.Getenv("INFURA_API_KEY")
	if infuraKey == "" {
		log.Fatal("INFURA_API_KEY not found in .env")
	}

	client, err := ethclient.Dial(fmt.Sprintf("https://mainnet.infura.io/v3/%s", infuraKey))
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	fmt.Println("Starting Ethereum gas price monitor...")
	fmt.Println("Checking prices every second...")

	var prices []float64
	const movingAverageWindow = 10

	for {
		select {
		case <-ticker.C:
			gasPrice, err := client.SuggestGasPrice(context.Background())
			if err != nil {
				log.Printf("Error getting gas price: %v", err)
				continue
			}
			gasPriceGwei := float64(gasPrice.Int64()) / 1e9
			
			prices = append(prices, gasPriceGwei)
			if len(prices) > movingAverageWindow {
				prices = prices[1:]
			}
			
			var sum float64
			for _, p := range prices {
				sum += p
			}
			avg := sum / float64(len(prices))

			fmt.Printf("\rCurrent gas price: %.2f Gwei | 10s Average: %.2f Gwei", 
				gasPriceGwei, avg)
		}
	}
} 