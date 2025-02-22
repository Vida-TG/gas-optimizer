package gas

import (
	"fmt"
	"math"
	"time"
)

type Analyzer struct {
	lowThreshold  float64
	highThreshold float64
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		lowThreshold:  20,
		highThreshold: 100,
	}
}

func (a *Analyzer) GetRecommendation(prices []PricePoint) string {
	if len(prices) < 60 {
		return "Collecting data..."
	}

	current := prices[len(prices)-1].Price
	avg := a.CalculateMovingAverage(prices, 60)

	if current < a.lowThreshold {
		return "🟢 Excellent time to send transaction!"
	} else if current < avg*0.9 {
		return "🟡 Good time to send transaction"
	} else if current > a.highThreshold {
		return "🔴 High gas prices, consider waiting"
	}

	bestTime := a.analyzeDailyPattern(prices)
	if bestTime != "" {
		return fmt.Sprintf("⏰ Consider waiting until %s", bestTime)
	}

	return "⚪ Average gas prices"
}

func (a *Analyzer) CalculateMovingAverage(prices []PricePoint, window int) float64 {
	if len(prices) == 0 {
		return 0
	}

	count := math.Min(float64(len(prices)), float64(window))
	sum := 0.0
	for i := len(prices) - int(count); i < len(prices); i++ {
		sum += prices[i].Price
	}
	return sum / count
}

func (a *Analyzer) analyzeDailyPattern(prices []PricePoint) string {
	if len(prices) < 24*60*60 {
		return ""
	}

	hourlyAverages := make(map[int]float64)
	hourlyCount := make(map[int]int)

	for _, p := range prices {
		hour := p.Time.Hour()
		hourlyAverages[hour] += p.Price
		hourlyCount[hour]++
	}

	lowestHour := 0
	lowestPrice := math.MaxFloat64

	for hour, total := range hourlyAverages {
		avg := total / float64(hourlyCount[hour])
		if avg < lowestPrice {
			lowestPrice = avg
			lowestHour = hour
		}
	}

	now := time.Now()
	if now.Hour() == lowestHour {
		return ""
	}

	return fmt.Sprintf("%02d:00", lowestHour)
} 