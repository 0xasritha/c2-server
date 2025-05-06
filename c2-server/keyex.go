package main

import (
	"crypto/rand"
	"math"
	"math/big"
)

func powMod(a, b, m int) int {
	result := 1
	a %= m
	for b > 0 {
		if b&1 == 1 {
			result = result * a % m
		}
		a = a * a % m
		b >>= 1
	}
	return result
}

func primeFactors(n int) []int {
	factors := []int{}
	for n%2 == 0 {
		factors = append(factors, 2)
		n /= 2
	}
	for i := 3; i <= int(math.Sqrt(float64(n))); i += 2 {
		for n%i == 0 {
			factors = append(factors, i)
			n /= i
		}
	}
	if n > 2 {
		factors = append(factors, n)
	}
	return factors
}

func generatePrimitiveRoot(p uint) int {
	int_p := int(p)
	phi := int_p - 1
	raw := primeFactors(phi)
	unique := []int{}
	for _, f := range raw {
		dup := false
		for _, u := range unique {
			if u == f {
				dup = true
				break
			}
		}
		if !dup {
			unique = append(unique, f)
		}
	}
	for g := 2; g < int_p; g++ {
		ok := true
		for _, q := range unique {
			if powMod(g, phi/q, int_p) == 1 {
				ok = false
				break
			}
		}
		if ok {
			return g
		}
	}
	return -1
}

func generatePrime() uint {
	prime, _ := rand.Prime(rand.Reader, 63)
	return uint(prime.Int64())
}

func generatePrivateKey() *big.Int {
	prime, _ := rand.Prime(rand.Reader, 128)
	privateKey, _ := rand.Int(rand.Reader, prime)
	return privateKey
}

func calculatePublicKey(prime int, privateKey *big.Int, base int) *big.Int {
	publicKey := new(big.Int).Exp(big.NewInt(int64(base)), privateKey, big.NewInt(int64(prime)))
	return publicKey
}

func calculateSharedSecret(prime int, privateKey, publicKey *big.Int) *big.Int {
	sharedSecret := new(big.Int).Exp(publicKey, privateKey, big.NewInt(int64(prime)))
	return sharedSecret
}
