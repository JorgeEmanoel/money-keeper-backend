package sec

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

func EncryptText(text string) string {
	keyString := os.Getenv("CIPHER_KEY")

	if len(keyString) < 1 {
		log.Fatalln("[!!!WARNING!!!] No cipher defined for cipher")
	}

	key, _ := hex.DecodeString(keyString)
	plainText := []byte(text)

	block, err := aes.NewCipher(key)

	if err != nil {
		log.Fatalf("Failed to create Cipher: %v\n", err)
	}

	aesGCM, err := cipher.NewGCM(block)

	if err != nil {
		log.Fatalf("Failed to create CGM: %v\n", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	ciphertext := aesGCM.Seal(nonce, nonce, plainText, nil)

	return fmt.Sprintf("%x", ciphertext)
}

func DecryptText(text string) string {
	keyString := os.Getenv("CIPHER_KEY")

	if len(keyString) < 1 {
		log.Fatalln("[!!!WARNING!!!] No cipher defined for cipher")
	}

	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(text)

	block, err := aes.NewCipher(key)

	if err != nil {
		log.Fatalf("Failed to create Cipher: %v\n", err)
	}

	aesGCM, err := cipher.NewGCM(block)

	if err != nil {
		log.Fatalf("Failed to create CGM: %v\n", err)
	}

	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plainText, err := aesGCM.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		log.Printf("Failed to Decrypt text: %v", err)
	}

	return string(plainText)
}
