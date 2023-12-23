package ratchet

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"lidx-core-lib/common"
	"lidx-core-lib/crypto/ecc"
	"lidx-core-lib/crypto/kdf"
	"lidx-core-lib/keys"
)

type IRatchet interface {
	GetId() string
	GetTotalSent() uint
	GetTotalRecieved() uint
	InitNewSession()
	InitRecievedSession(yourEphemeralPubKey ecc.IECPublicKey)
	PopulateMessage(content []byte) *Message
	OnSend(message *Message)
	OnRecieved(message *Message)
	Save(PIN []byte) *RachetStore
}

type RachetStore struct {
	RachetId           string   `json:"rachet_id"`
	RootKey            string   `json:"root_key"`
	ChainSendKey       string   `json:"chain_send_key"`
	ChainRecvKey       string   `json:"chain_recv_key"`
	TotalMessageSent   uint     `json:"total_message_sent"`
	TotalMessageRecv   uint     `json:"total_message_recv"`
	MissingMessageKeys []string `json:"missing_message_keys"`
}

type Ratchet struct {
	RatchetId            string
	MyKeyBundle          *keys.InternalKeyBundle
	YourKeyBundle        *keys.ExternalKeyBundle
	RootKey              []byte
	ChainSendKey         []byte
	ChainRecieveKey      []byte
	MissingMessageKeys   [][]byte
	TotalMessageSent     uint
	TotalMessageRecieved uint
	RootKeyEncrypted     bool
}

func NewRachetFromInternal(internalKeyBundle *keys.InternalKeyBundle, externalBundle *keys.ExternalKeyBundle) (*Ratchet, error) {
	if !externalBundle.Verify() {
		fmt.Println("Cannot verify external key bundle")
		return nil, fmt.Errorf("Cannot verify external key bundle")
	}
	id, _ := uuid.NewUUID()
	ratchet := &Ratchet{
		RatchetId:            id.String(),
		MyKeyBundle:          internalKeyBundle,
		YourKeyBundle:        externalBundle,
		TotalMessageSent:     0,
		TotalMessageRecieved: 0,
		RootKeyEncrypted:     true,
	}
	ratchet.InitNewSession()
	return ratchet, nil
}

func NewRachetFromExternal(internalKeyBundle *keys.InternalKeyBundle, externalBundle *keys.ExternalKeyBundle, yourEphemeralPubKey ecc.IECPublicKey, ratchetId string) (*Ratchet, error) {
	ratchet := &Ratchet{
		RatchetId:            ratchetId,
		MyKeyBundle:          internalKeyBundle,
		YourKeyBundle:        externalBundle,
		TotalMessageSent:     0,
		TotalMessageRecieved: 0,
		RootKeyEncrypted:     true,
	}
	ratchet.InitRecievedSession(yourEphemeralPubKey)
	return ratchet, nil
}

func LoadRachet(rachetJsonString string, PIN []byte) *Ratchet {
	var rachetStore RachetStore
	err := json.Unmarshal([]byte(rachetJsonString), &rachetStore)
	if err != nil {
		fmt.Println("Cannot parase json")
		return nil
	}
	rootKey, err := common.DecryptHashedData(common.DecodeToByte(rachetStore.RootKey), PIN)
	if err != nil {
		fmt.Println("Cannot decrypt key")
		return nil
	}
	chainSendKey, err := common.DecryptHashedData(common.DecodeToByte(rachetStore.ChainSendKey), PIN)
	if err != nil {
		fmt.Println("Cannot decrypt key")
		return nil
	}
	chainRecvKey, err := common.DecryptHashedData(common.DecodeToByte(rachetStore.ChainRecvKey), PIN)
	if err != nil {
		fmt.Println("Cannot decrypt key")
		return nil
	}
	var missingKeys [][]byte
	for _, missKey := range rachetStore.MissingMessageKeys {
		decryptedMissKey, err := common.DecryptHashedData(common.DecodeToByte(missKey), PIN)
		if err != nil {
			fmt.Println("Cannot ecnrypt key")
			return nil
		}
		missingKeys = append(missingKeys, decryptedMissKey)
	}
	return &Ratchet{
		RatchetId:            rachetStore.RachetId,
		MyKeyBundle:          nil,
		YourKeyBundle:        nil,
		RootKey:              rootKey,
		ChainSendKey:         chainSendKey,
		ChainRecieveKey:      chainRecvKey,
		MissingMessageKeys:   missingKeys,
		TotalMessageSent:     rachetStore.TotalMessageSent,
		TotalMessageRecieved: rachetStore.TotalMessageRecv,
		RootKeyEncrypted:     false,
	}
}

func (r *Ratchet) GetId() string {
	return r.RatchetId
}

func (r *Ratchet) GetTotalSent() uint {
	return r.TotalMessageSent
}

func (r *Ratchet) GetTotalRecieved() uint {
	return r.TotalMessageRecieved
}

func (r *Ratchet) PopulateMessage(content []byte) *Message {
	return &Message{
		Index:        r.TotalMessageSent,
		RatchetID:    r.RatchetId,
		PlainMessage: content,
	}
}

func (r *Ratchet) InitNewSession() {
	if !r.RootKeyEncrypted {
		fmt.Println("Not a new session")
		return
	}
	// Performance X3DH
	ikA := r.MyKeyBundle.IdentityKey.PrivateKey()
	pkB := r.YourKeyBundle.GetPreKey()
	dh1, _ := ikA.CalculateCommonSecret(pkB)

	ekA := r.MyKeyBundle.EphemeralKey.PrivateKey()
	ikB := r.YourKeyBundle.GetIndentityKey()
	dh2, _ := ekA.CalculateCommonSecret(ikB)

	dh3, _ := ekA.CalculateCommonSecret(pkB)

	var preKdf []byte
	opkB, _ := r.YourKeyBundle.GetOneTimePreKey()
	if opkB != nil {
		dh4, _ := ekA.CalculateCommonSecret(opkB)
		preKdf = common.ConcatBytes(dh1, dh2, dh3, dh4)
	} else {
		preKdf = common.ConcatBytes(dh1, dh2, dh3)
	}

	rootKey, e := kdf.DoKDF(preKdf)
	if e != nil {
		print("Cannot generate root key", e)
		return
	}
	r.RootKey = rootKey
	r.RootKeyEncrypted = false
}

func (r *Ratchet) InitRecievedSession(ephemeralKey ecc.IECPublicKey) {
	if !r.RootKeyEncrypted {
		fmt.Println("Not a new session")
		return
	}
	// Performance X3DH
	pkA := r.MyKeyBundle.PreKey.PrivateKey()
	ikB := r.YourKeyBundle.GetIndentityKey()
	dh1, _ := pkA.CalculateCommonSecret(ikB)

	ikA := r.MyKeyBundle.IdentityKey.PrivateKey()
	ekB := ephemeralKey
	dh2, _ := ikA.CalculateCommonSecret(ekB)

	dh3, _ := pkA.CalculateCommonSecret(ekB)

	var preKdf []byte
	preKdf = common.ConcatBytes(dh1, dh2, dh3)
	_, opkId := r.YourKeyBundle.GetOneTimePreKey()
	if !common.IsStringEmpty(&opkId) {
		var opkA *ecc.ECKeyPair
		opkA = nil
		for _, v := range r.MyKeyBundle.OneTimeKey {
			opkA = v
		}
		if opkA != nil {
			dh4, _ := opkA.PrivateKey().CalculateCommonSecret(ekB)
			preKdf = common.ConcatBytes(dh1, dh2, dh3, dh4)
		}
	} else {
		preKdf = common.ConcatBytes(dh1, dh2, dh3)
	}

	rootKey, e := kdf.DoKDF(preKdf)
	if e != nil {
		fmt.Println("Cannot generate root key", e)
		return
	}
	r.RootKey = rootKey
	r.RootKeyEncrypted = false
}

func (r *Ratchet) OnSend(message *Message) {
	message.RatchetID = r.RatchetId
	var encyptedKey []byte
	if r.TotalMessageSent == 0 {
		encyptedKey, _ = kdf.DoKDF(r.RootKey)
	} else {
		encyptedKey, _ = kdf.DoKDF(r.ChainSendKey)
	}
	message.Encrypt(encyptedKey)
	r.TotalMessageSent++
	r.ChainSendKey = encyptedKey
	message.Index = r.TotalMessageSent
}

func (r *Ratchet) OnRecieved(message *Message) {
	if message.RatchetID != r.RatchetId {
		fmt.Println("Wrong rachet")
		return
	}
	var decryptKey []byte
	if (r.TotalMessageRecieved + 1) < message.Index {
		offset := (message.Index - r.TotalMessageRecieved)
		decryptKey, _ = kdf.DoKDF(r.RootKey)
		for i := 0; i < int(offset)-1; i++ {
			decryptKey, _ = kdf.DoKDF(decryptKey)
			r.ChainRecieveKey = decryptKey
			r.MissingMessageKeys = append(r.MissingMessageKeys, decryptKey)
		}
	} else if r.TotalMessageRecieved == 0 {
		decryptKey, _ = kdf.DoKDF(r.RootKey)
	} else {
		decryptKey, _ = kdf.DoKDF(r.ChainRecieveKey)
	}
	message.Decrypt(decryptKey)
	r.TotalMessageRecieved++
	r.ChainRecieveKey = decryptKey
}

func (r *Ratchet) Save(PIN []byte) *RachetStore {
	encryptedRootKey, err := common.EncryptAndHash(r.RootKey, PIN)
	if err != nil {
		fmt.Println("Cannot ecnrypt key")
		return nil
	}
	encryptedChainSendKey, err := common.EncryptAndHash(r.ChainSendKey, PIN)
	if err != nil {
		fmt.Println("Cannot ecnrypt key")
		return nil
	}
	encryptedRecvSendKey, err := common.EncryptAndHash(r.ChainRecieveKey, PIN)
	if err != nil {
		fmt.Println("Cannot ecnrypt key")
		return nil
	}
	var missingKeys []string
	for _, missKey := range r.MissingMessageKeys {
		encryptedMissKey, err := common.EncryptAndHash(missKey, PIN)
		if err != nil {
			fmt.Println("Cannot ecnrypt key")
			return nil
		}
		missingKeys = append(missingKeys, common.EncodeToString(encryptedMissKey))
	}
	return &RachetStore{
		RachetId:           r.RatchetId,
		RootKey:            common.EncodeToString(encryptedRootKey),
		ChainSendKey:       common.EncodeToString(encryptedChainSendKey),
		ChainRecvKey:       common.EncodeToString(encryptedRecvSendKey),
		TotalMessageSent:   r.GetTotalSent(),
		TotalMessageRecv:   r.GetTotalRecieved(),
		MissingMessageKeys: missingKeys,
	}
}
