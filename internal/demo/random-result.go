package demo

import (
	"math/rand"
)

type RandomResult struct {
	Mean    float64
	StdDev  float64
	ErrPerc int
}

func (rr *RandomResult) float64() float64 {
	x := rand.NormFloat64()*rr.StdDev + rr.Mean
	if x < 0 {
		x = 0
	}
	return x
}

func (rr *RandomResult) int64() int64 {
	return int64(rr.float64())
}

func (rr *RandomResult) success() bool {
	// errPerc = 0 .. 100
	// n       = 1 .. 100

	// if errPerc=  0 -> every n is greater than 0         -> always success
	// if errPerc=100 -> every n is less or equal than 100 -> always error
	n := rand.Intn(100) + 1
	return n > rr.ErrPerc
}
