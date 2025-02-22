package gas

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Monitor struct {
	client    *ethclient.Client
	analyzer  *Analyzer
	stopChan  chan struct{}
	prices    []PricePoint
	isRunning bool
}

type PricePoint struct {
	Price float64
	Time  time.Time
}

func NewMonitor(infuraKey string) (*Monitor, error) {
	client, err := ethclient.Dial(fmt.Sprintf("https://mainnet.infura.io/v3/%s", infuraKey))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Infura: %w", err)
	}

	return &Monitor{
		client:   client,
		analyzer: NewAnalyzer(),
		stopChan: make(chan struct{}),
	}, nil
}

func (m *Monitor) Start() {
	if m.isRunning {
		return
	}
	m.isRunning = true

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	fmt.Println("Starting Ethereum gas price monitor...")
	fmt.Println("Collecting data and analyzing patterns...")

	for {
		select {
		case <-m.stopChan:
			return
		case <-ticker.C:
			gasPrice, err := m.client.SuggestGasPrice(context.Background())
			if err != nil {
				fmt.Printf("\rError getting gas price: %v", err)
				continue
			}

			priceGwei := float64(gasPrice.Int64()) / 1e9
			point := PricePoint{
				Price: priceGwei,
				Time:  time.Now(),
			}

			m.prices = append(m.prices, point)
			if len(m.prices) > 24*60*60 {
				m.prices = m.prices[1:]
			}

			recommendation := m.analyzer.GetRecommendation(m.prices)
			m.printStatus(priceGwei, recommendation)
		}
	}
}

func (m *Monitor) Stop() {
	if !m.isRunning {
		return
	}
	m.stopChan <- struct{}{}
	m.isRunning = false
	m.client.Close()
}

func (m *Monitor) printStatus(currentPrice float64, recommendation string) {
	avg := m.analyzer.CalculateMovingAverage(m.prices, 10)
	fmt.Printf("\rCurrent: %.2f Gwei | 10s Avg: %.2f Gwei | %s    \n", 
		currentPrice, avg, recommendation)
} 