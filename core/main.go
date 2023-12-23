package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"lidx-core-lib/common"
	"lidx-core-lib/crypto/ecc"
	"lidx-core-lib/keys"
	"lidx-core-lib/ratchet"
	"log"
	"syscall/js"
)

var INTERNAL_KEY_PREFIX = "INTERNAL_KEY_"
var EXTERNAL_KEY_PREFIX = "EXTERNAL_KEY_PREFIX_"

var INTERNAL_KEY_STORAGE = make(map[string]*keys.InternalKeyBundle)

var EXTERNAL_KEY_STORAGE = make(map[string]*keys.ExternalKeyBundle)

var RATCHET_STORAGE = make(map[string]*ratchet.Ratchet)

var PIN = ""

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	done := make(chan struct{}, 0)

	go js.Global().Set("startUp", js.FuncOf(startUp))
	go js.Global().Set("generateInternalKeyBundle", js.FuncOf(generateInternalKeyBundle))
	go js.Global().Set("loadInternalKey", js.FuncOf(loadInternalKey))
	go js.Global().Set("saveInternalKey", js.FuncOf(saveInternalKey))
	go js.Global().Set("regeneratePreKey", js.FuncOf(regeneratePreKey))
	go js.Global().Set("populateExternalKeyBundle", js.FuncOf(populateExternalKeyBundle))
	go js.Global().Set("initRatchetFromInternal", js.FuncOf(initRatchetFromInternal))
	go js.Global().Set("initRatchetFromExternal", js.FuncOf(initRatchetFromExternal))
	go js.Global().Set("saveRatchet", js.FuncOf(saveRatchet))
	go js.Global().Set("loadRatchet", js.FuncOf(loadRatchet))
	go js.Global().Set("isRatchetExist", js.FuncOf(isRatchetExist))
	go js.Global().Set("sendMessage", js.FuncOf(sendMessage))
	go js.Global().Set("receiveMessage", js.FuncOf(receiveMessage))
	go js.Global().Set("initVoipSessionFromInternal", js.FuncOf(initVoipSessionFromInternal))
	go js.Global().Set("initVoipSessionFromExternal", js.FuncOf(initVoipSessionFromExternal))

	<-done
}

// TODO make function call for js client
// First arg is PIN
func startUp(this js.Value, args []js.Value) interface{} {
	PIN = args[0].String()
	return nil
}

// Internal Key API
func generateInternalKeyBundle(this js.Value, args []js.Value) interface{} {
	internalKey := keys.NewInternalKeyBundle()
	return insertInternalKeyToStorage(internalKey)
}

func regeneratePreKey(this js.Value, args []js.Value) interface{} {
	internalKey := loadInternalKeyFromStorage()
	internalKey.GeneratePreKey()
	externalKeyBundle := internalKey.GenerateExternalKey()
	insertExternalKeyToStorage(externalKeyBundle)
	return convertToJsObject(externalKeyBundle.ToDto())
}

// param 1 : key json string
func loadInternalKey(this js.Value, args []js.Value) interface{} {
	decodedPin := common.StringToByte(PIN)
	internalKey := keys.LoadInternalKey(args[0].String(), decodedPin)
	return insertInternalKeyToStorage(internalKey)
}

func saveInternalKey(this js.Value, args []js.Value) interface{} {
	internalKey := loadInternalKeyFromStorage()
	decodedPin := common.StringToByte(PIN)
	return convertToJsObject(internalKey.Save(decodedPin))
}

func populateExternalKeyBundle(this js.Value, args []js.Value) interface{} {
	internalKey := loadInternalKeyFromStorage()
	externalKeyBundle := internalKey.GenerateExternalKey()
	insertExternalKeyToStorage(externalKeyBundle)
	/*resultMap := make(map[string]interface{})
	resultMap["keyId"] = externalKeyId
	resultMap["keyBundle"] = externalKeyBundle.ToDto()*/
	return convertToJsObject(externalKeyBundle.ToDto())
}

// Rachet API
// First param is other user external key bundle
func initRatchetFromInternal(this js.Value, args []js.Value) interface{} {
	externalKeyString := args[0].String()
	externalKeyBundle, err := keys.NewExternalKeyFromJson(externalKeyString)
	if err != nil {
		log.Println("cannot read external key bundle")
		return nil
	}
	internalKey := loadInternalKeyFromStorage()
	if internalKey.EphemeralKey == nil {
		internalKey.GenerateEphemeralKey()
	}

	rachet, _ := ratchet.NewRachetFromInternal(internalKey, externalKeyBundle)
	insertRatchetToStorage(rachet)
	ePubKey, _ := internalKey.EphemeralKey.PublicKey().Serialize()

	resultMap := make(map[string]interface{})
	resultMap["ratchetId"] = rachet.GetId()
	resultMap["keyBundle"] = externalKeyBundle.ToDto()
	resultMap["ephemeralKey"] = common.EncodeToString(ePubKey)
	return convertToJsObject(resultMap)
}

// (1) arg is externalKeyJsonString
// (2) is external ephemeralPubKeyString
// (3) is other ratchetId
func initRatchetFromExternal(this js.Value, args []js.Value) interface{} {
	externalKeyString := args[0].String()
	externalEphemeralPubKeyString := args[1].String()
	externalRatchetId := args[2].String()
	externalKeyBundle, err := keys.NewExternalKeyFromJson(externalKeyString)
	if err != nil {
		log.Println("cannot read external key bundle")
		return nil
	}

	internalKey := loadInternalKeyFromStorage()

	externalEphemeralPubKey, _ := ecc.DeserializePublicKey(common.DecodeToByte(externalEphemeralPubKeyString))

	rachet, _ := ratchet.NewRachetFromExternal(internalKey, externalKeyBundle, externalEphemeralPubKey, externalRatchetId)

	insertRatchetToStorage(rachet)
	resultMap := make(map[string]interface{})
	resultMap["ratchetId"] = rachet.GetId()
	return convertToJsObject(resultMap)
}

func initVoipSessionFromInternal(this js.Value, args []js.Value) interface{} {
	externalKeyString := args[0].String()
	externalKeyBundle, err := keys.NewExternalKeyFromJson(externalKeyString)
	if err != nil {
		log.Println("cannot read external key bundle")
		return nil
	}
	internalKey := loadInternalKeyFromStorage()
	if internalKey.EphemeralKey == nil {
		internalKey.GenerateEphemeralKey()
	}

	rachet, _ := ratchet.NewRachetFromInternal(internalKey, externalKeyBundle)
	ePubKey, _ := internalKey.EphemeralKey.PublicKey().Serialize()

	resultMap := make(map[string]interface{})
	resultMap["sessionCommonSecretKey"] = common.EncodeToString(rachet.RootKey)
	resultMap["ephemeralKey"] = common.EncodeToString(ePubKey)
	return convertToJsObject(resultMap)
}

// (1) arg is externalKeyJsonString
// (2) is external ephemeralPubKeyString
func initVoipSessionFromExternal(this js.Value, args []js.Value) interface{} {
	externalKeyString := args[0].String()
	externalEphemeralPubKeyString := args[1].String()
	externalKeyBundle, err := keys.NewExternalKeyFromJson(externalKeyString)
	if err != nil {
		log.Println("cannot read external key bundle")
		return nil
	}

	internalKey := loadInternalKeyFromStorage()

	externalEphemeralPubKey, _ := ecc.DeserializePublicKey(common.DecodeToByte(externalEphemeralPubKeyString))

	rachet, _ := ratchet.NewRachetFromExternal(internalKey, externalKeyBundle, externalEphemeralPubKey, "")

	resultMap := make(map[string]interface{})
	resultMap["sessionCommonSecretKey"] = common.EncodeToString(rachet.RootKey)
	return convertToJsObject(resultMap)
}

func saveRatchet(this js.Value, args []js.Value) interface{} {
	ratchetId := args[0].String()
	rachet := loadRatchetFromStorage(ratchetId)
	if rachet == nil {
		log.Println("cannot load ratchet")
		return nil
	}
	return convertToJsObject(rachet.Save(common.StringToByte(PIN)))
}

func loadRatchet(this js.Value, args []js.Value) interface{} {
	rachetJson := args[0].String()
	rachet := ratchet.LoadRachet(rachetJson, common.StringToByte(PIN))
	return insertRatchetToStorage(rachet)
}

func isRatchetExist(this js.Value, args []js.Value) interface{} {
	messageJson := args[0].String()
	storedRatchet := RATCHET_STORAGE[messageJson]
	if storedRatchet != nil {
		log.Println("ratchet existed")
		return true
	}
	return false
}

// Message API
// (1) argument is ratchet ID or conversion ID, (2) argument is dedicate this message is binary or not, (3) is content
func sendMessage(this js.Value, args []js.Value) interface{} {
	ratchetId := args[0].String()
	rachet := loadRatchetFromStorage(ratchetId)
	if rachet == nil {
		log.Println("cannot find ratchet")
		return nil
	}
	isBinary := args[1].Bool()
	content := args[2].String()
	var msg *ratchet.Message
	if isBinary {
		msg = rachet.PopulateMessage(common.DecodeToByte(content))
	} else {
		msg = rachet.PopulateMessage(common.StringToByte(content))
	}
	rachet.OnSend(msg)
	msgDto := msg.ToDto()
	msgDto.IsBinary = isBinary
	return convertToJsObject(msgDto)
}

// (1) argument is message dto
func receiveMessage(this js.Value, args []js.Value) interface{} {
	messageJson := args[0].String()
	var messageDto ratchet.MessageDto
	err := json.Unmarshal([]byte(messageJson), &messageDto)
	if err != nil {
		log.Println("Cannot parase json")
		return nil
	}
	recvMsg := ratchet.CreateMessageFromDto(&messageDto)
	rachet := loadRatchetFromStorage(recvMsg.RatchetID)
	if rachet == nil {
		log.Println("cannot find ratchet")
		return nil
	}
	rachet.OnRecieved(recvMsg)
	if messageDto.IsBinary {
		return common.EncodeToString(recvMsg.PlainMessage)
	} else {
		return string(recvMsg.PlainMessage)
	}
}

// Utils
func convertToJsObject(data any) map[string]interface{} {
	jsString, _ := json.Marshal(data)
	var result map[string]interface{}
	err := json.Unmarshal(jsString, &result)
	if err != nil {
		log.Println("cannot parse object", err)
		return nil
	}
	return result
}

func insertInternalKeyToStorage(internalKey *keys.InternalKeyBundle) string {
	if len(INTERNAL_KEY_STORAGE) != 0 {
		log.Println("internal key existed")
		return ""
	}
	internalKeyId, _ := uuid.NewUUID()
	INTERNAL_KEY_STORAGE[internalKeyId.String()] = internalKey
	return internalKeyId.String()
}

func loadInternalKeyFromStorage() *keys.InternalKeyBundle {
	for _, v := range INTERNAL_KEY_STORAGE {
		if v != nil {
			return v
		}
	}
	_ = fmt.Errorf("cannot find internal key")
	return nil
}

func insertExternalKeyToStorage(externalKey *keys.ExternalKeyBundle) string {
	mapKey, _ := uuid.NewUUID()
	EXTERNAL_KEY_STORAGE[mapKey.String()] = externalKey
	return mapKey.String()
}

func loadExternalKeyFromStorage(keyId string) *keys.ExternalKeyBundle {
	externalKey := EXTERNAL_KEY_STORAGE[keyId]
	if externalKey == nil {
		log.Println("cannot find external key")
		return nil
	}
	return externalKey
}

func insertRatchetToStorage(rachet *ratchet.Ratchet) string {
	if rachet == nil {
		log.Println("ratchet is null")
		return ""
	}
	storedRatchet := RATCHET_STORAGE[rachet.GetId()]
	if storedRatchet != nil {
		log.Println("ratchet existed")
		return ""
	}
	RATCHET_STORAGE[rachet.GetId()] = rachet
	return rachet.GetId()
}

func loadRatchetFromStorage(ratchetId string) *ratchet.Ratchet {
	storedRatchet := RATCHET_STORAGE[ratchetId]
	if storedRatchet == nil {
		log.Println("cannot find rachet")
		return nil
	}
	return storedRatchet
}
