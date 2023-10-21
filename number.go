package layout

import (
	"strconv"
	"strings"
)

func ToInt(v string) int {
	v = strings.TrimSuffix(v, "px")
	i, _ := strconv.Atoi(v)
	return i
}

func ToFloat64(v string) float64 {
	v = strings.TrimSuffix(v, "px")
	f, _ := strconv.ParseFloat(v, 64)
	return f
}

func gcd(a, b int) int {
	if a < 0 || b < 0 {
		return 0
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func GetGcd(n []int) int {
	g := n[0]
	for i := 1; i < len(n); i++ {
		g = gcd(g, n[i])
	}
	return g
}

func GcdInts(in []int) []int {
	n := GetGcd(in)
	out := make([]int, len(in))
	for i, num := range in {
		out[i] = num / n
	}
	return out
}

func SumFloat64(collection []float64) float64 {
	var sum float64 = 0
	for _, val := range collection {
		sum += val
	}
	return sum
}

func SumInt(collection []int) int {
	var sum int = 0
	for _, val := range collection {
		sum += val
	}
	return sum
}
