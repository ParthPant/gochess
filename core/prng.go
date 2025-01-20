package core

type PRNG struct {
	state uint64
}

func (p *PRNG) Seed(seed uint64) {
	p.state = seed
}

func (p *PRNG) Rand64() uint64 {
	p.state ^= p.state >> 12
	p.state ^= p.state << 25
	p.state ^= p.state >> 27
	return p.state * 2685821657736338717
}

func (p *PRNG) SparseRand64() uint64 {
	return p.Rand64() & p.Rand64() & p.Rand64()
}
