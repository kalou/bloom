package bloom

import (
	"fmt"
	"log"
	"math"
	"os"
	"testing"
)

var b *Filter

func TestMain(m *testing.M) {
	b = NewFilter(1500, .01)
	os.Exit(m.Run())
}

func TestBloomSize(t *testing.T) {
	if b.k != 7 || b.m != 14378 {
		t.Fatal("Bad size estimation")
	}

	if len(b.set) != 225 {
		t.Fatal("Bad array size", len(b.set))
	}
}

func TestBloomUnset(t *testing.T) {
	log.Println("Testing nothing is set")
	if b.IsSet("some string") {
		t.Fatal("some string set")
	}
}

func TestBloomSet(t *testing.T) {
	log.Println("Testing set works")
	b.Set("some string")
	if !b.IsSet("some string") {
		t.Fatal("some string unset")
	}
}

func TestEstimate(t *testing.T) {
	log.Println("Testing estimate")
	tot := 0
	for _, j := range []int{30, 50, 200, 850, 1400} {
		for i := 0; i < j; i++ {
			str := fmt.Sprintf("Test string", i)
			b.Set(str)
			tot++
		}

		log.Println("After", j, "count is", b.EstimateN())
		if math.Abs(float64(j-int(b.EstimateN()))) > float64(10) {
			t.Fatal("Counting test failed", j, b.EstimateN())
		}
	}
}
