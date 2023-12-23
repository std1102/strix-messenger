package ecc

import (
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"
	"lidx-core-lib/common"
)

type IECPrivateKey interface {
	Serialize(PIN []byte) ([]byte, error)
	PrivateKey() *ecdsa.PrivateKey
	CalculateCommonSecret(otherPub IECPublicKey) ([]byte, error)
}

type MyPrivateKey struct {
	privateKey *ecdsa.PrivateKey
}

func DeserializePrivateKey(input []byte, PIN []byte) (*MyPrivateKey, error) {
	decryptData, err := common.DecryptHashedData(input, PIN)
	if err != nil {
		return nil, fmt.Errorf("Cannot decryp private key", err)
	}
	privKey, err := x509.ParseECPrivateKey(decryptData)
	if err != nil {
		return nil, fmt.Errorf("Cannot read private key", err)
	}
	return &MyPrivateKey{
		privateKey: privKey,
	}, nil
}

func (priv *MyPrivateKey) CalculateCommonSecret(otherPub IECPublicKey) ([]byte, error) {
	dhPriv, _ := priv.privateKey.ECDH()
	dhPub, _ := otherPub.PublicKey().ECDH()
	result, err := dhPriv.ECDH(dhPub)
	if err != nil {
		return nil, fmt.Errorf("Cannot caculate common secret")
	}
	return result, nil
}

// TODO Should encrypt this
func (priv *MyPrivateKey) Serialize(PIN []byte) ([]byte, error) {
	privateKey, err := x509.MarshalECPrivateKey(priv.privateKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot get private key")
	}
	ePKey, err := common.EncryptAndHash(privateKey, PIN)
	if err != nil {
		return nil, fmt.Errorf("Cannot encrypt private key ", err)
	}
	return ePKey, nil
}

func (priv *MyPrivateKey) PrivateKey() *ecdsa.PrivateKey {
	return priv.privateKey
}

func NewPrivateKey(priv *ecdsa.PrivateKey) *MyPrivateKey {
	return &MyPrivateKey{
		privateKey: priv,
	}
}
