package ratchet

import (
	"encoding/json"
	"fmt"
	"lidx-core-lib/common"
)

type Message struct {
	RatchetID     string
	Index         uint
	PlainMessage  []byte
	CipherMessage []byte
}

type MessageDto struct {
	RatchetID     string `json:"chatSessionId"`
	Index         uint   `json:"index"`
	CipherMessage string `json:"cipherMessage"`
	IsBinary      bool   `json:"isBinary"`
}

func CreateMessageFromDto(messageDto *MessageDto) *Message {
	return &Message{
		RatchetID:     messageDto.RatchetID,
		Index:         messageDto.Index,
		PlainMessage:  nil,
		CipherMessage: common.DecodeToByte(messageDto.CipherMessage),
	}
}

func CreateMessageFromJson(jsonString string) *Message {
	var messageDto MessageDto
	err := json.Unmarshal([]byte(jsonString), &messageDto)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Cannot parase json")
		return nil
	}
	return &Message{
		RatchetID:     messageDto.RatchetID,
		Index:         messageDto.Index,
		PlainMessage:  nil,
		CipherMessage: common.DecodeToByte(messageDto.CipherMessage),
	}
}

func (m *Message) Encrypt(key []byte) {
	var e error
	m.CipherMessage, e = common.EncryptAndHash(m.PlainMessage, key)
	if e != nil {
		fmt.Println("Cannot encrypt message ", e)
		return
	}
}

func (m *Message) Decrypt(key []byte) {
	prePlainText, err := common.DecryptHashedData(m.CipherMessage, key)
	if err != nil {
		fmt.Println("Cannot decrypt message ", err)
		return
	}
	m.PlainMessage = prePlainText
}

func (m *Message) ToDto() *MessageDto {
	return &MessageDto{
		RatchetID:     m.RatchetID,
		Index:         m.Index,
		CipherMessage: common.EncodeToString(m.CipherMessage),
	}
}
