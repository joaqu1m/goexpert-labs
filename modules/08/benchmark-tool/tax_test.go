package covertool

import (
	"testing"
	"time"
)

func BenchmarkTax(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Calc(1000, 10)
	}
}

func BenchmarkTaxZeroIncome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Calc(0, 10)
	}
}

func BenchmarkTaxDelayed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Sleep(10 * time.Nanosecond)
		Calc(1000, 10)
	}
}
