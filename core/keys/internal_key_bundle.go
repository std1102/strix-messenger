package keys

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"lidx-core-lib/crypto/ecc"
)

type InternalKeyBundle struct {
	IdentityKey  *ecc.ECKeyPair
	EphemeralKey *ecc.ECKeyPair
	PreKey       *ecc.ECKeyPair
	OneTimeKey   map[string]*ecc.ECKeyPair
}

type InternalKeyBundleStore struct {
	IdentityKey *ecc.ECKeyPairStore            `json:"identity_key"`
	PreKey      *ecc.ECKeyPairStore            `json:"pre_key"`
	OneTimeKey  map[string]*ecc.ECKeyPairStore `json:"one_time_key"`
}

func LoadInternalKey(keyJsonString string, PIN []byte) *InternalKeyBundle {
	var internalBundleStore InternalKeyBundleStore
	err := json.Unmarshal([]byte(keyJsonString), &internalBundleStore)
	if err != nil {
		fmt.Println("Cannot parse json")
		return nil
	}
	identityKey, _ := ecc.DeSerializeKey(internalBundleStore.IdentityKey, PIN)
	preKey, _ := ecc.DeSerializeKey(internalBundleStore.PreKey, PIN)
	oneTimeKeyMap := make(map[string]*ecc.ECKeyPair)
	for k, v := range internalBundleStore.OneTimeKey {
		dKey, _ := ecc.DeSerializeKey(v, PIN)
		oneTimeKeyMap[k] = dKey
	}
	return &InternalKeyBundle{
		IdentityKey: identityKey,
		PreKey:      preKey,
		OneTimeKey:  oneTimeKeyMap,
	}
}

func NewInternalKeyBundle() *InternalKeyBundle {
	oneTimeKeys := make(map[string]*ecc.ECKeyPair)
	key, _ := uuid.NewUUID()
	oneTimeKeys[key.String()] = ecc.GenerateKeyPair()
	return &InternalKeyBundle{
		IdentityKey:  ecc.GenerateKeyPair(),
		EphemeralKey: ecc.GenerateKeyPair(),
		PreKey:       ecc.GenerateKeyPair(),
		OneTimeKey:   oneTimeKeys,
	}
}

func (internalKey *InternalKeyBundle) Save(PIN []byte) *InternalKeyBundleStore {
	identityKey := internalKey.IdentityKey.Save(PIN)
	preKey := internalKey.PreKey.Save(PIN)
	oneTimeKeyMap := make(map[string]*ecc.ECKeyPairStore)
	for k, v := range internalKey.OneTimeKey {
		oneTimeKeyMap[k] = v.Save(PIN)
	}
	return &InternalKeyBundleStore{
		IdentityKey: identityKey,
		PreKey:      preKey,
		OneTimeKey:  oneTimeKeyMap,
	}
}

func (internalKey *InternalKeyBundle) GenerateEphemeralKey() ecc.ECKeyPair {
	internalKey.EphemeralKey = ecc.GenerateKeyPair()
	return *internalKey.EphemeralKey
}

func (internalKey *InternalKeyBundle) GeneratePreKey() ecc.ECKeyPair {
	internalKey.PreKey = ecc.GenerateKeyPair()
	return *internalKey.PreKey
}

func (internalKey *InternalKeyBundle) GenerateExternalKey() *ExternalKeyBundle {
	yIk := internalKey.IdentityKey
	yPk := internalKey.PreKey

	singer := ecc.FromKeyPair(yIk)

	yPkR, _ := yPk.PublicKey().Serialize()

	bspk, _ := singer.Sign(yPkR)

	var iOKey *ecc.ECKeyPair
	var oKeyId string

	for k, v := range internalKey.OneTimeKey {
		oKeyId = k
		iOKey = v
	}

	iOkR, _ := iOKey.PublicKey().Serialize()

	okeySig, _ := singer.Sign(iOkR)

	return NewExternalKeyBundle(
		yIk.PublicKey(),
		nil,
		yPk.PublicKey(),
		bspk,
		oKeyId,
		iOKey.PublicKey(),
		okeySig,
	)
}
