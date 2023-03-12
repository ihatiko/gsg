package postgres_types_generators

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"math"
)

func NumericGenerator() (any, string) {
	//TODO precision and range
	data := toFixed(gofakeit.Float64Range(0, 1000), 2)
	return data, fmt.Sprintf("%f", data)
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	data := float64(round(num*output)) / output
	return data
}
