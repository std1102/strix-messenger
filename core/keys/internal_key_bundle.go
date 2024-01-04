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
	OneTimeKey   *ecc.ECKeyPair
	PreKeys      map[string]*ecc.ECKeyPair
}

type InternalKeyBundleStore struct {
	IdentityKey *ecc.ECKeyPairStore            `json:"identity_key"`
	PreKeys     map[string]*ecc.ECKeyPairStore `json:"pre_keys"`
}

func LoadInternalKey(keyJsonString string, PIN []byte) *InternalKeyBundle {
	var internalBundleStore InternalKeyBundleStore
	err := json.Unmarshal([]byte(keyJsonString), &internalBundleStore)
	if err != nil {
		fmt.Println("Cannot parse json")
		return nil
	}
	identityKey, _ := ecc.DeSerializeKey(internalBundleStore.IdentityKey, PIN)
	preKeyMap := make(map[string]*ecc.ECKeyPair)
	for k, v := range internalBundleStore.PreKeys {
		dKey, _ := ecc.DeSerializeKey(v, PIN)
		preKeyMap[k] = dKey
	}
	return &InternalKeyBundle{
		IdentityKey: identityKey,
		PreKeys:     preKeyMap,
	}
}

func NewInternalKeyBundle() *InternalKeyBundle {
	oneTimeKeys := make(map[string]*ecc.ECKeyPair)
	key, _ := uuid.NewUUID()
	oneTimeKeys[key.String()] = ecc.GenerateKeyPair()
	return &InternalKeyBundle{
		IdentityKey:  ecc.GenerateKeyPair(),
		EphemeralKey: ecc.GenerateKeyPair(),
		OneTimeKey:   ecc.GenerateKeyPair(),
		PreKeys:      oneTimeKeys,
	}
}

func (internalKey *InternalKeyBundle) Save(PIN []byte) *InternalKeyBundleStore {
	identityKey := internalKey.IdentityKey.Save(PIN)
	oneTimeKeyMap := make(map[string]*ecc.ECKeyPairStore)
	for k, v := range internalKey.PreKeys {
		oneTimeKeyMap[k] = v.Save(PIN)
	}
	return &InternalKeyBundleStore{
		IdentityKey: identityKey,
		PreKeys:     oneTimeKeyMap,
	}
}

func (internalKey *InternalKeyBundle) GenerateEphemeralKey() ecc.ECKeyPair {
	internalKey.EphemeralKey = ecc.GenerateKeyPair()
	return *internalKey.EphemeralKey
}

func (internalKey *InternalKeyBundle) GeneratePreKey() ecc.ECKeyPair {
	internalKey.OneTimeKey = ecc.GenerateKeyPair()
	return *internalKey.OneTimeKey
}

func (internalKey *InternalKeyBundle) GenerateExternalKey() *ExternalKeyBundle {
	yIk := internalKey.IdentityKey
	singer := ecc.FromKeyPair(yIk)

	var pk *ecc.ECKeyPair
	var pkId string

	for k, v := range internalKey.PreKeys {
		pkId = k
		pk = v
	}

	pkPublic, _ := pk.PublicKey().Serialize()

	pkSig, _ := singer.Sign(pkPublic)

	return NewExternalKeyBundle(
		yIk.PublicKey(),
		nil,
		pkId,
		pk.PublicKey(),
		pkSig,
	)
}
