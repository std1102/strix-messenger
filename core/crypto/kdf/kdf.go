package kdf

import (
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/hkdf"
)

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
