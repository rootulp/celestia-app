package chainspec

func numValidators() *int {
	numValidators := 1
	return &numValidators
}

func numFullNodes() *int {
	numValidators := 0
	return &numValidators
}

func gasAdjustment() *float64 {
	gasAdjustment := 2.0
	return &gasAdjustment
}
