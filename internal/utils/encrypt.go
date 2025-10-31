package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"os"
)

func HashText(key []byte, text []byte) ([]byte, error) {
	cipherText, err := Encrypt(key, text)
	if err != nil {
		return nil, err
	}

	return []byte(base64.StdEncoding.EncodeToString(cipherText)), nil
}

func Encrypt(key, text []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("cipher key must be 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCTR(block, iv)

	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func GetSecrets() (string, string) {
	key := []byte(os.Getenv("CIPHER_KEY"))

	passwordCiphertext, err := base64.StdEncoding.DecodeString(os.Getenv("PASSWORD"))
	if err != nil {
		log.Fatal("Error reading the password: ", err)
	}

	password, err := Decrypt(key, passwordCiphertext)
	if err != nil {
		log.Fatal("Error decrypting the password: ", err)
	}

	pinCiphertext, err := base64.StdEncoding.DecodeString(os.Getenv("PIN"))
	if err != nil {
		log.Fatal("Error reading the pin: ", err)
	}

	pin, err := Decrypt(key, pinCiphertext)
	if err != nil {
		log.Fatal("Error decrypting the pin: ", err)
	}

	return string(password), string(pin)
}

func Decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCTR(block, iv)

	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}
