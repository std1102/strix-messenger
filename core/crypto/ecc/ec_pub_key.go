package ecc

import (
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"
	"lidx-core-lib/common"
	"log"
)

type IECPublicKey interface {
	Serialize() ([]byte, error)
	PublicKey() *ecdsa.PublicKey
}

type MyPublicKey struct {
	publicKey *ecdsa.PublicKey
}

func DeserializePublicKey(input []byte) (*MyPublicKey, error) {
	pubKey, err := x509.ParsePKIXPublicKey(input)
	if err != nil {
		return nil, fmt.Errorf("Cannot read private key")
	}
	return &MyPublicKey{
		publicKey: pubKey.(*ecdsa.PublicKey),
	}, nil
}

// Convert this public key to byte array
// Can convert this byte array to base64
func (pub *MyPublicKey) Serialize() ([]byte, error) {
	x509encode, e := x509.MarshalPKIXPublicKey(pub.publicKey)
	if e != nil {
		log.Println("Cannot read internal public key ", e)
		return nil, e
	}
	return x509encode, nil
}

func (pub *MyPublicKey) PublicKey() *ecdsa.PublicKey {
	return pub.publicKey
}

func ParsePublicKey(key string) IECPublicKey {
	genericPublicKey, _ := x509.ParsePKIXPublicKey(common.DecodeToByte(key))
	return &MyPublicKey{
		publicKey: genericPublicKey.(*ecdsa.PublicKey),
	}
}

func NewPublicKey(pub *ecdsa.PublicKey) IECPublicKey {
	return &MyPublicKey{
		publicKey: pub,
	}
}
