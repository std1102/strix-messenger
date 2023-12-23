package test

import (
	"bytes"
	"fmt"
	"lidx-core-lib/crypto/ecc"
	"testing"
)

func TestEccSign(t *testing.T) {
	fmt.Println("Signing test")
	keyPair := ecc.GenerateKeyPair()
	msg := []byte("TEST STRING")
	signer := ecc.FromKeyPair(keyPair)
	sig, _ := signer.Sign(msg)
	if !signer.Verify(msg, sig) {
		t.Error("Cannot sign")
	}
}

func TestEccDh(t *testing.T) {
	fmt.Println("DH test")
	aKeyPair := ecc.GenerateKeyPair()
	bKeyPair := ecc.GenerateKeyPair()

	aSceret, _ := aKeyPair.PrivateKey().CalculateCommonSecret(bKeyPair.PublicKey())
	bSecret, _ := bKeyPair.PrivateKey().CalculateCommonSecret(aKeyPair.PublicKey())

	if len(aSceret) == 0 || len(bSecret) == 0 {
		t.Error("Common secret fail")
	}

	if bytes.Compare(aSceret, bSecret) != 0 {
		t.Error("Fail DH")
	}
}
