package keys

import (
	"encoding/json"
	"fmt"
	"lidx-core-lib/common"
	"lidx-core-lib/crypto/ecc"
)

type ExternalKeyBundle struct {
	IdentityKey  ecc.IECPublicKey
	EphemeralKey ecc.IECPublicKey
	PreKeyId     string
	PreKey       ecc.IECPublicKey
	PreKeySig    []byte
}

// Encode in base64
type ExternalKeyBundleDto struct {
	IdentityKey   string `json:"identityKey,omitempty"`
	OneTimeKey    string `json:"oneTimeKey,omitempty"`
	OneTimeKeySig string `json:"oneTimeKeySig,omitempty"`
	PreKeyId      string `json:"preKeyId,omitempty"`
	PreKey        string `json:"preKey,omitempty"`
	PreKeySig     string `json:"preKeySig,omitempty"`
}

// TODO convert base64 string to key material
func ExternalKeyFromJson(jsonString string) *ExternalKeyBundle {
	var result ExternalKeyBundle
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		fmt.Println("Cannot read external key bundle json")
		return nil
	}
	return &result
}

func NewExternalKeyBundle(
	ik ecc.IECPublicKey,
	ek ecc.IECPublicKey,
	pkid string,
	pk ecc.IECPublicKey,
	spk []byte,
) *ExternalKeyBundle {
	return &ExternalKeyBundle{
		IdentityKey:  ik,
		EphemeralKey: ek,
		PreKeyId:     pkid,
		PreKey:       pk,
		PreKeySig:    spk,
	}
}

func NewExternalKeyFromJson(jsonString string) (*ExternalKeyBundle, error) {
	var dto ExternalKeyBundleDto
	err := json.Unmarshal([]byte(jsonString), &dto)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Cannot deserialize json")
	}
	iKey, _ := ecc.DeserializePublicKey(common.DecodeToByte(dto.IdentityKey))
	pKey, _ := ecc.DeserializePublicKey(common.DecodeToByte(dto.PreKey))
	result := &ExternalKeyBundle{
		IdentityKey: iKey,
		PreKeyId:    dto.PreKeyId,
		PreKey:      pKey,
		PreKeySig:   common.DecodeToByte(dto.PreKeySig),
	}
	return result, nil
}

func (keyBundle *ExternalKeyBundle) GetIdentityKey() ecc.IECPublicKey {
	return keyBundle.IdentityKey
}

func (keyBundle *ExternalKeyBundle) GetEphemeralKey() ecc.IECPublicKey {
	return keyBundle.EphemeralKey
}

func (keyBundle *ExternalKeyBundle) GetPreKey() (ecc.IECPublicKey, string) {
	return keyBundle.PreKey, keyBundle.PreKeyId
}

func (keyBundle *ExternalKeyBundle) SetEphemeralKey(ePubKey *ecc.IECPublicKey) {
	keyBundle.EphemeralKey = *ePubKey
}

func (keyBundle *ExternalKeyBundle) Verify() bool {
	userIdentityKey := keyBundle.GetIdentityKey()
	if userIdentityKey == nil {
		fmt.Println("Cannot get user identity key")
		return false
	}
	signer := ecc.FromPublicKey(userIdentityKey)
	pKey, _ := keyBundle.PreKey.Serialize()
	if !signer.Verify(pKey, keyBundle.PreKeySig) {
		fmt.Println("Cannot verify user pre key")
		return false
	}
	if keyBundle.PreKey != nil && !common.IsByteArrayEmpty(&keyBundle.PreKeySig) {
		oKey, _ := keyBundle.PreKey.Serialize()
		if !signer.Verify(oKey, keyBundle.PreKeySig) {
			fmt.Println("Cannot verify user one-time key")
			return false
		}
	}
	return true
}

func (keyBundle *ExternalKeyBundle) ToDto() *ExternalKeyBundleDto {
	var pid, pk, pks string
	if keyBundle.PreKey != nil && keyBundle.PreKeySig != nil && keyBundle.PreKeyId != "" {
		pid = keyBundle.PreKeyId
		okResult, _ := keyBundle.PreKey.Serialize()
		pk = common.EncodeToString(okResult)
		pks = common.EncodeToString(keyBundle.PreKeySig)
	} else {
		pid = ""
		pk = ""
		pks = ""
	}

	iKey, _ := keyBundle.IdentityKey.Serialize()

	dto := &ExternalKeyBundleDto{
		IdentityKey: common.EncodeToString(iKey),
		PreKeyId:    pid,
		PreKey:      pk,
		PreKeySig:   pks,
	}
	return dto
}
