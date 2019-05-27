package rand

import (
	"math/rand"
)

type RandGen struct {
	gen *rand.Rand
	max int
}

func NewRandGen(seed int64, max int) *RandGen {
	randSrc := rand.NewSource(seed)
	return &RandGen{gen: rand.New(randSrc), max: max}
}

func (r *RandGen) NextInt() int {
	return r.gen.Intn(r.max)
}
