package util

import "math"

// F = C * 1.8 + 32
func CtoF(c float64) float64 { return c*1.8 + 32 }

// K = C + 273  (conforme requisito fornecido)
func CtoK_Custom(c float64) float64 { return c + 273 }

// Arredonda com uma casa (ex.: 28.49 -> 28.5)
func Round1(v float64) float64 {
	return math.Round(v*10) / 10
}
