package benchmark

import "math/big"

// What does it do?
func mysteriousFunctionRec(n int) *big.Int {
	if n <= 1 {
		return big.NewInt(1)
	}

	res := mysteriousFunctionRec(n - 1)
	return res.Mul(big.NewInt(int64(n)), res)
}

// What does it do?
func mysteriousFunction(n int) *big.Int {
	if n <= 1 {
		return big.NewInt(1)
	}

	result := big.NewInt(1)
	for i := 1; i <= n; i++ {
		result.Mul(result, big.NewInt(int64(i)))
	}
	return result
}
