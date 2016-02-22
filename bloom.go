package bloom

import (
	"hash/fnv"
	"math"
)

type Filter struct {
	m uint
	k uint

	set []uint64
}

// Determines what bits map to this string
// returns an array of size k
func (b *Filter) bits(what string) (res []uint) {
	fnv := fnv.New64a()
	fnv.Write([]byte(what))
	h := fnv.Sum64()
	y := uint((h >> 32) & 0xffffffff)
	z := uint(h & 0xffffffff)
	for i := uint(0); i < b.k; i++ {
		res = append(res, (y+z*i)%b.m)
	}

	return res
}

// Set the given bit
func (b *Filter) setbit(n uint) {
	b.set[n/64] |= (1 << (n % 64))
}

func (b *Filter) bitisset(n uint) bool {
	return (b.set[n/64] & (1 << (n % 64))) != 0
}

func popcnt(x uint64) uint {
	var res uint64
	for ; x > 0; x >>= 1 {
		res += x & 1
	}
	return uint(res)
}

func (b *Filter) countbits() (tot uint) {
	for i := uint(0); i < (b.m/64 + 1); i++ {
		tot += popcnt(b.set[i])
	}

	return
}

// Initialize a new bloom filter, computing
// m and k from the estimated number of inserts
// nr and the false positive probability p
func NewFilter(n int, p float64) *Filter {
	b := &Filter{}

	b.m = uint(math.Ceil((float64(n) * math.Log(p)) / math.Log(1.0/(math.Pow(2.0, math.Log(2.0))))))
	b.k = uint(math.Ceil(math.Log(2.0) * float64(b.m) / float64(n)))

	b.set = make([]uint64, b.m/64+1)

	return b
}

// Set string in bloom filter
func (b *Filter) Set(what string) {
	for _, n := range b.bits(what) {
		b.setbit(n)
	}
}

// Find out whether what is set in the filter with
// false positive probability p
func (b *Filter) IsSet(what string) bool {
	for _, n := range b.bits(what) {
		if !b.bitisset(n) {
			return false
		}
	}
	return true
}

// Estimate number of elements in filter
func (b *Filter) EstimateN() uint {
	X := b.countbits()
	return uint(-float64(b.m) * math.Log(1-float64(X)/float64(b.m)) / float64(b.k))
}
