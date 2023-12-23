package router

type RegisterDto struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	AliasName string `json:"aliasName"`
	Email     string `json:"email"`
}

type LoginDto struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	RememberMe   bool   `json:"rememberMe"`
	RefreshToken string `json:"refreshToken"`
	LoginType    string `json:"loginType"`
}

type LoginResponseDto struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	LoggedInAt   string `json:"loggedInAt"`
}

type ExternalKeyBundleDto struct {
	IdentityKey   string `json:"identityKey,omitempty"`
	PreKey        string `json:"preKey,omitempty"`
	PreKeySig     string `json:"preKeySig,omitempty"`
	OneTimeKeyId  string `json:"oneTimeKeyId,omitempty"`
	OneTimeKey    string `json:"oneTimeKey,omitempty"`
	OneTimeKeySig string `json:"oneTimeKeySig,omitempty"`
}

type UserDto struct {
	Id        string `json:"id"`
	UserName  string `json:"userName"`
	AliasName string `json:"aliasName"`
	Avatar    string `json:"avatar"`
}

type UserRequestDto struct {
	UserNames []string `json:"usernames"`
}

const (
	CHAT_NEW    = "CHAT_NEW"
	CHAT_TEXT   = "CHAT_TEXT"
	CHAT_FILE   = "CHAT_FILE"
	CHAT_VOIP   = "CHAT_VOIP"
	CHAT_VIDEO  = "CHAT_VIDEO"
	CHAT_AUDIO  = "CHAT_AUDIO"
	CHAT_ACCEPT = "CHAT_ACCEPT"
	CHAT_CLOSE  = "CHAT_CLOSE"
)

type MessageDto struct {
	Type           string      `json:"type"`
	SenderUsername string      `json:"senderUsername"`
	PlainMessage   *string     `json:"plainMessage"`
	ChatSessionId  string      `json:"chatSessionId"`
	Index          uint64      `json:"index"`
	CipherMessage  string      `json:"cipherMessage"`
	FilePath       *string     `json:"filePath"`
	IsBinary       bool        `json:"isBinary"`
	AdditionalData interface{} `json:"additionalData"`
}

type ChatSessionDto struct {
	ChatSessionId    string               `json:"chatSessionId"`
	EphemeralKey     string               `json:"ephemeralKey"`
	ReceiverUserName string               `json:"receiverUserName"`
	SenderUserName   string               `json:"senderUserName"`
	SenderKeyBundle  ExternalKeyBundleDto `json:"senderKeyBundle"`
}
