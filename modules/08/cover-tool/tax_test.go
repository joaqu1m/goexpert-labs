package covertool

import (
	"testing"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		name     string
		income   float64
		taxRate  float64
		expected float64
	}{
		{"Zero income", 0, 10, 0},
		{"Negative income", -1000, 10, 0},
		{"Zero tax rate", 1000, 0, 0},
		{"Valid input", 1000, 10, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Calc(tt.income, tt.taxRate)
			if result != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}
