package rand

import (
	"sync/atomic"
	"time"
)

type rand struct {
	a      uint64
	c      uint64
	m      uint64
	m2     float64
	values []uint64
	cap    int32
	idx    atomic.Int32
	seed   uint64
}

func New(maxnum int) *rand {
	r := rand{
		seed: uint64(time.Now().UnixNano()),
		a:    6364136223846793005, // Множитель
		c:    1442695040888963407, // Приращение
		m:    1 << 63,             // Модуль, используется 2^63 для типа uint64
		m2:   1 << 53,             // Для получения float64 из uint64. У float64 мантисса 53 бита
	}

	buffer := make([]uint64, 0, maxnum)

	for idx := 0; idx < maxnum; idx++ { //nolint:intrange
		buffer = append(buffer, r.next())
	}

	r.values = buffer
	r.cap = int32(maxnum)
	r.idx = atomic.Int32{}

	return &r
}

func (r *rand) next() uint64 {
	// https://cs.opensource.google/go/go/+/refs/tags/go1.22.0:src/math/rand/rng.go;l=18
	r.seed = (r.a*r.seed + r.c) % (r.m - 1)

	return r.seed
}

func (r *rand) Uint64() uint64 {
	i := r.idx.Load()
	value := r.values[i]
	r.idx.Store((i + 1) % r.cap)

	return value
}

func (r *rand) Float64() float64 {
	return float64(r.Uint64()>>11) / r.m2
}

func (r *rand) Intn(n int) int {
	if n <= 0 {
		panic("invalid argument to IntN")
	}

	//	return int(float64(n) * r.Float64())
	if n&(n-1) == 0 { // степень двойки
		return int(r.Uint64() & uint64(n-1))
	}

	max := ^uint64(0) - (^uint64(0) % uint64(n))
	for {
		v := r.Uint64()
		if v < max {
			return int(v % uint64(n))
		}
	}
}
