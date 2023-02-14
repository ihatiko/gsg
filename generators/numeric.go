package generators

import (
	"github.com/brianvoe/gofakeit/v6"
	"math"
)

func NumericGenerator() any {
	//TODO precision and range
	return toFixed(gofakeit.Float64Range(0, 1000), 2)
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
