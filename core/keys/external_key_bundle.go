package keys

import (
	"encoding/json"
	"fmt"
	"lidx-core-lib/common"
	"lidx-core-lib/crypto/ecc"
)

type ExternalKeyBundle struct {
	IdentityKey   ecc.IECPublicKey
	EphemeralKey  ecc.IECPublicKey
	PreKey        ecc.IECPublicKey
	PreKeySig     []byte
	OneTimeKeyId  string
	OneTimeKey    ecc.IECPublicKey
	OneTimeKeySig []byte
}

// Encode in base64
type ExternalKeyBundleDto struct {
	IdentityKey   string `json:"identityKey,omitempty"`
	PreKey        string `json:"preKey,omitempty"`
	PreKeySig     string `json:"preKeySig,omitempty"`
	OneTimeKeyId  string `json:"oneTimeKeyId,omitempty"`
	OneTimeKey    string `json:"oneTimeKey,omitempty"`
	OneTimeKeySig string `json:"oneTimeKeySig,omitempty"`
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
	pk ecc.IECPublicKey,
	spk []byte,
	okid string,
	ok ecc.IECPublicKey,
	sok []byte,
) *ExternalKeyBundle {
	return &ExternalKeyBundle{
		IdentityKey:   ik,
		EphemeralKey:  ek,
		PreKey:        pk,
		PreKeySig:     spk,
		OneTimeKeyId:  okid,
		OneTimeKey:    ok,
		OneTimeKeySig: sok,
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
		PreKey:      pKey,
		PreKeySig:   common.DecodeToByte(dto.PreKeySig),
	}
	if dto.OneTimeKey != "" && dto.OneTimeKeySig != "" && dto.OneTimeKeyId != "" {
		oKey, _ := ecc.DeserializePublicKey(common.DecodeToByte(dto.OneTimeKey))
		result.OneTimeKeyId = dto.OneTimeKeyId
		result.OneTimeKey = oKey
		result.OneTimeKeySig = common.DecodeToByte(dto.OneTimeKeySig)
	}
	return result, nil
}

func (keyBundle *ExternalKeyBundle) GetIndentityKey() ecc.IECPublicKey {
	return keyBundle.IdentityKey
}

func (keyBundle *ExternalKeyBundle) GetEphemeralKey() ecc.IECPublicKey {
	return keyBundle.EphemeralKey
}

func (keyBundle *ExternalKeyBundle) GetPreKey() ecc.IECPublicKey {
	return keyBundle.PreKey
}

func (keyBundle *ExternalKeyBundle) GetOneTimePreKey() (ecc.IECPublicKey, string) {
	return keyBundle.OneTimeKey, keyBundle.OneTimeKeyId
}

func (keyBundle *ExternalKeyBundle) SetEphemeralKey(ePubKey *ecc.IECPublicKey) {
	keyBundle.EphemeralKey = *ePubKey
}

func (keyBundle *ExternalKeyBundle) Verify() bool {
	userIdentityKey := keyBundle.GetIndentityKey()
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
	if keyBundle.OneTimeKey != nil && !common.IsByteArrayEmpty(&keyBundle.OneTimeKeySig) {
		oKey, _ := keyBundle.OneTimeKey.Serialize()
		if !signer.Verify(oKey, keyBundle.OneTimeKeySig) {
			fmt.Println("Cannot verify user one-time key")
			return false
		}
	}
	return true
}

func (keyBundle *ExternalKeyBundle) ToDto() *ExternalKeyBundleDto {
	var oid, ok, oks string
	if keyBundle.OneTimeKey != nil && keyBundle.OneTimeKeySig != nil && keyBundle.OneTimeKeyId != "" {
		oid = keyBundle.OneTimeKeyId
		okResult, _ := keyBundle.OneTimeKey.Serialize()
		ok = common.EncodeToString(okResult)
		oks = common.EncodeToString(keyBundle.OneTimeKeySig)
	} else {
		oid = ""
		ok = ""
		oks = ""
	}

	iKey, _ := keyBundle.IdentityKey.Serialize()
	pKey, _ := keyBundle.PreKey.Serialize()

	dto := &ExternalKeyBundleDto{
		IdentityKey:   common.EncodeToString(iKey),
		PreKey:        common.EncodeToString(pKey),
		PreKeySig:     common.EncodeToString(keyBundle.PreKeySig),
		OneTimeKeyId:  oid,
		OneTimeKey:    ok,
		OneTimeKeySig: oks,
	}
	return dto
}
