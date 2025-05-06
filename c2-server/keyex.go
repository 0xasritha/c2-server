package main

import (
	"crypto/rand"
	"math"
	"math/big"
)

// CryptoConfig stores shared crypto settings w/ the C2Server
type CryptoConfig struct {
	Prime          uint     `json:"prime"`
	Generator      int      `json:"generator"`
	OtherPublicKey *big.Int `json:"other_public_key"`
	SelfPublicKey  *big.Int `json:"self_public_key"`
	sharedSecret   *big.Int // field will be empty until exchange occurs
	privateKey     *big.Int
}

func NewCryptoConfig() *CryptoConfig {
	c := new(CryptoConfig)
	return c
}

func (c *CryptoConfig) generatePrime() {
	prime, _ := rand.Prime(rand.Reader, 63)
	c.Prime = uint(prime.Int64())
}

func (c *CryptoConfig) generateGenerator(p uint) {
	c.Generator = generatePrimitiveRoot(p)
	// should i make this uint? -1 is for error
}

func (c *CryptoConfig) generatePrivateKey() {
	prime, _ := rand.Prime(rand.Reader, 128)
	c.privateKey, _ = rand.Int(rand.Reader, prime)
}

func (c *CryptoConfig) calculatePublicKey(prime uint, privateKey *big.Int, base int) {
	c.SelfPublicKey = new(big.Int).Exp(big.NewInt(int64(base)), privateKey, big.NewInt(int64(prime)))
}

func (c *CryptoConfig) calculateSharedSecret(prime int, privateKey, publicKey *big.Int) {
	c.sharedSecret = new(big.Int).Exp(publicKey, privateKey, big.NewInt(int64(prime)))
}

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
