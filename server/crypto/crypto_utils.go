package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/hkdf"
	"io"
)

func AesGCMEncrypt(key, plainText []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("", err)
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, plainText, nil)
	return ciphertext, nonce, nil
}

func AesGCMDecrypt(key, cipherText, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("", err)
	}

	return aesgcm.Open(nil, nonce, cipherText, nil)
}

func DoKDF(keyMaterial []byte) ([]byte, error) {
	kdf := hkdf.New(sha256.New, keyMaterial, nil, nil)
	result := make([]byte, 16)
	_, err := kdf.Read(result)
	if err != nil {
		fmt.Println("Cannot do KDF")
		return nil, err
	}
	return result, nil
}
