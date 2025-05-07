package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// TODO: fix nonce implementation
// TODO: is this the correct way to do it?
var (
	// 32 byte long secret key
	sharedSecretKey string = "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
)

func main() {
	// This will successfully encrypt & decrypt
	ciphertext1 := encrypt("This is some sensitive information")
	fmt.Printf("Encrypted ciphertext 1: %x \n", ciphertext1)

	plaintext1 := decrypt(ciphertext1)
	fmt.Printf("Decrypted plaintext 1: %s \n", plaintext1)

	// This will successfully encrypt & decrypt as well.
	ciphertext2 := encrypt("hi_this_is_a_flag")
	fmt.Printf("Encrypted ciphertext 2: %x \n", ciphertext2)

	plaintext2 := decrypt(ciphertext2)
	fmt.Printf("Decrypted plaintext 2: %s \n", plaintext2)
}

func encrypt(plaintext string) string {
	aes, err := aes.NewCipher([]byte(sharedSecretKey))
	if err != nil {
		panic(err)
	}
	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}
	// cipher text is nonce+ciphertext
	// when decrypting, we can use size of nonce to separate it from cipher text
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func decrypt(encodedCipherText string) string {
	ciphertext, err := base64.StdEncoding.DecodeString(string(encodedCipherText))
	if err != nil {
		panic(err)
	}

	aes, err := aes.NewCipher([]byte(sharedSecretKey))
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plainText, err := gcm.Open(nil, []byte(nonce), ciphertext, nil)
	if err != nil {
		panic(err)
	}

	return string(plainText)
}
