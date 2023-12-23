package ecc

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
)

type ISigner interface {
	Sign(hash []byte) (sig []byte, err error)
	Verify(pubKey []byte, sig []byte) (bool, error)
}

type ECDSASigner struct {
	privateKey IECPrivateKey
	publicKey  IECPublicKey
}

func FromPrivateKey(priv IECPrivateKey) *ECDSASigner {
	return &ECDSASigner{
		privateKey: priv,
	}
}

func FromPublicKey(pub IECPublicKey) *ECDSASigner {
	return &ECDSASigner{
		publicKey: pub,
	}
}

func FromKeyPair(keyPair *ECKeyPair) *ECDSASigner {
	return &ECDSASigner{
		publicKey:  keyPair.PublicKey(),
		privateKey: keyPair.PrivateKey(),
	}
}

// return ASN1 encoded
func (E *ECDSASigner) Sign(hash []byte) (sig []byte, err error) {
	sig, err = ecdsa.SignASN1(rand.Reader, E.privateKey.PrivateKey(), hash)
	if err != nil {
		fmt.Println("Cannot sign", err)
		return nil, err
	}
	return sig, nil
}

func (E *ECDSASigner) Verify(hash, sig []byte) bool {
	return ecdsa.VerifyASN1(E.publicKey.PublicKey(), hash, sig)
}
