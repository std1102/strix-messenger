package ecc

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"lidx-core-lib/common"
)

func NewECKeyPair(publicKey IECPublicKey, privateKey IECPrivateKey) *ECKeyPair {
	keypair := ECKeyPair{
		publicKey:  publicKey,
		privateKey: privateKey,
	}

	return &keypair
}

// ECKeyPair is a combination of both public and private elliptic curve keys.
type ECKeyPair struct {
	publicKey  IECPublicKey
	privateKey IECPrivateKey
}

type ECKeyPairStore struct {
	PublicKey   string `json:"public_key,omitempty"`
	PrivateKey  string `json:"private_key,omitempty"`
	PrivateHash string `json:"private_hash,omitempty"`
	PublicHash  string `json:"public_hash,omitempty"`
}

func DeSerializeKey(keyStore *ECKeyPairStore, PIN []byte) (*ECKeyPair, error) {
	var pubKey, privKey []byte
	var pubHash, privHash []byte
	pubKey = common.DecodeToByte(keyStore.PublicKey)
	privKey = common.DecodeToByte(keyStore.PrivateKey)
	pubHash = common.DecodeToByte(keyStore.PublicHash)
	privHash = common.DecodeToByte(keyStore.PrivateHash)
	currentPubHash := sha256.Sum256(pubKey)
	currentPrivHash := sha256.Sum256(privKey)
	if bytes.Compare(currentPubHash[:], pubHash) != 0 && bytes.Compare(currentPrivHash[:], privHash) != 0 {
		return nil, fmt.Errorf("Key store hash not match")
	}
	mPrivKey, err := DeserializePrivateKey(privKey, PIN)
	if err != nil {
		return nil, fmt.Errorf("", err)
	}
	mPubKey, err := DeserializePublicKey(pubKey)
	if err != nil {
		return nil, fmt.Errorf("", err)
	}
	return &ECKeyPair{
		privateKey: mPrivKey,
		publicKey:  mPubKey,
	}, nil
}

func DeSerializeKeyStoreString(jsonString string, PIN []byte) (*ECKeyPair, error) {
	var keyStore ECKeyPairStore
	err := json.Unmarshal([]byte(jsonString), &keyStore)
	if err != nil {
		return nil, fmt.Errorf("Cannot load key store")
	}
	var pubKey, privKey []byte
	var pubHash, privHash []byte
	pubKey = common.DecodeToByte(keyStore.PublicKey)
	privKey = common.DecodeToByte(keyStore.PrivateKey)
	pubHash = common.DecodeToByte(keyStore.PublicHash)
	privHash = common.DecodeToByte(keyStore.PrivateHash)
	currentPubHash := sha256.Sum256(pubKey)
	currentPrivHash := sha256.Sum256(privKey)
	if bytes.Compare(currentPubHash[:], pubHash) != 0 && bytes.Compare(currentPrivHash[:], privHash) != 0 {
		return nil, fmt.Errorf("Key store hash not match")
	}
	mPrivKey, err := DeserializePrivateKey(privKey, PIN)
	if err != nil {
		return nil, fmt.Errorf("", err)
	}
	mPubKey, err := DeserializePublicKey(pubKey)
	if err != nil {
		return nil, fmt.Errorf("", err)
	}
	return &ECKeyPair{
		privateKey: mPrivKey,
		publicKey:  mPubKey,
	}, nil
}

func GenerateKeyPair() *ECKeyPair {
	priv, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	return &ECKeyPair{
		privateKey: NewPrivateKey(priv),
		publicKey:  NewPublicKey(&priv.PublicKey),
	}
}

// PublicKey returns the public key from the key pair.
func (e *ECKeyPair) PublicKey() IECPublicKey {
	return e.publicKey
}

// PrivateKey returns the private key from the key pair.
func (e *ECKeyPair) PrivateKey() IECPrivateKey {
	return e.privateKey
}

func (e *ECKeyPair) Save(PIN []byte) *ECKeyPairStore {
	var pubKey, privKey []byte
	var pubHash, privHash [32]byte
	var _ error
	pubKey, _ = e.PublicKey().Serialize()
	privKey, _ = e.PrivateKey().Serialize(PIN)
	pubHash = sha256.Sum256(pubKey)
	privHash = sha256.Sum256(privKey)
	return &ECKeyPairStore{
		PublicKey:   common.EncodeToString(pubKey),
		PrivateKey:  common.EncodeToString(privKey),
		PublicHash:  common.EncodeToString(pubHash[:]),
		PrivateHash: common.EncodeToString(privHash[:]),
	}
}
