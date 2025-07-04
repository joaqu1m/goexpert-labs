package covertool

func Calc(income float64, taxRate float64) float64 {
	if income <= 0 || taxRate < 0 {
		return 0
	}
	return income * taxRate / 100
}
