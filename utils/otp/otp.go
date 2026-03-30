package otp

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOTP(digits int) (string, error) {
	limit := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil)

	n, err := rand.Int(rand.Reader, limit)
	if err != nil {
		return "", err
	}

	otp := fmt.Sprintf("%0*d", digits, n)

	return otp, nil
}
