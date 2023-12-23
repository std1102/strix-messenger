package common

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"lidx-core-lib/crypto/aes"
	"lidx-core-lib/crypto/kdf"
)

func ConcatBytes(bytes ...[]byte) []byte {
	result := *new([]byte)
	for i := 0; i < len(bytes); i++ {
		result = append(result, bytes[i]...)
	}
	return result
}

func IsStringEmpty(data *string) bool {
	return data == nil || len(*data) == 0
}

func IsAllStringEmpty(inputs ...*string) bool {
	counter := 0
	for i := 0; i < len(inputs); i++ {
		if !IsStringEmpty(inputs[i]) {
			counter++
		}
	}
	return counter == len(inputs)-1
}

func IsArrayEmpty(data *[]any) bool {
	return data == nil || len(*data) == 0
}

func IsAllArrayEmpty(inputs ...*[]any) bool {
	counter := 0
	for i := 0; i < len(inputs); i++ {
		if !IsArrayEmpty(inputs[i]) {
			counter++
		}
	}
	return counter == len(inputs)-1
}

func IsByteArrayEmpty(data *[]byte) bool {
	return data == nil || len(*data) == 0
}

func IsAllByteArrayEmpty(inputs ...*[]byte) bool {
	counter := 0
	for i := 0; i < len(inputs); i++ {
		if !IsByteArrayEmpty(inputs[i]) {
			counter++
		}
	}
	return counter == len(inputs)-1
}

func RandomByt(size int) ([]byte, error) {
	result := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, result); err != nil {
		return nil, err
	}
	return result, nil
}

// TODO IF ERROR, IT'S MUST BE HERE
func EncodeToString(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// TODO IF ERROR, IT'S MUST BE HERE
func DecodeToByte(data string) []byte {
	result, e := base64.RawURLEncoding.DecodeString(data)
	if e != nil {
		fmt.Println("Cannot parse ", e)
		return nil
	}
	return result
}

func StringToByte(data string) []byte {
	return []byte(data)
}

func SkipError(value interface{}, _ error) interface{} {
	return value
}

func EncryptAndHash(plainText, key []byte) ([]byte, error) {
	encrypKey, err := kdf.DoKDF(key)
	if err != nil {
		return nil, fmt.Errorf("Cannot make pass phase")
	}
	cipherText, nonce, err := aes.AesGCMEncrypt(encrypKey, plainText)
	if err != nil {
		return nil, fmt.Errorf("Cannot encrypt ", err)
	}
	hash := sha256.Sum256(plainText)
	return ConcatBytes(hash[:], nonce, cipherText), nil
}

func DecryptHashedData(cipherText, key []byte) ([]byte, error) {
	encrypKey, err := kdf.DoKDF(key)
	if err != nil {
		return nil, fmt.Errorf("Cannot make pass phase")
	}
	hash := cipherText[0:32]
	nonce := cipherText[32:44]
	cipherData := cipherText[44:]
	plainText, err := aes.AesGCMDecrypt(encrypKey, cipherData, nonce)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt")
	}
	plainHash := sha256.Sum256(plainText)
	if bytes.Compare(hash, plainHash[:]) != 0 {
		return nil, fmt.Errorf("Hash not equals")
	}
	return plainText, nil
}
