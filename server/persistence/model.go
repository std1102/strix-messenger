package persistence

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID                uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Username          string     `gorm:"type:varchar(255);not null;unique"`
	Password          string     `gorm:"type:varchar(255);not null"`
	AliasName         string     `gorm:"type:varchar(255)"`
	Email             string     `gorm:"type:varchar(255)"`
	Avatar            *string    `gorm:"type:varchar(500)"`
	IdentityKey       string     `gorm:"type:varchar(255)"`
	PreKeyCreatedTime *time.Time `gorm:"type:time"`
	PreKeys           []*PreKeys `gorm:"foreignKey:UserId"`
	Devices           []*Device  `gorm:"foreignKey:UserId"`
	CreatedAt         time.Time  `gorm:"type:time;default:current_timestamp;not null"`
}

type PreKeys struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserId       uuid.UUID `gorm:"type:uuid"`
	Key          string    `gorm:"type:varchar(255);not null"`
	KeySignature string    `gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time `gorm:"type:time;default:current_timestamp;not null"`
	Owner        *User     `gorm:"foreignKey:UserId"`
}

type Device struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserId         uuid.UUID `gorm:"type:uuid"`
	PhysicDeviceId string    `gorm:"type:varchar(255);not null"`
	IsPrimary      bool      `gorm:"type:varchar(255);not null"`
	PublicKey      string    `gorm:"type:varchar(255);not null"`
	Owner          *User     `gorm:"foreignKey:UserId"`
	CreatedAt      time.Time `gorm:"type:varchar(255);default:current_timestamp;not null"`
	LastLoggedIn   time.Time `gorm:"type:time;default:current_timestamp;not null"`
}

type ChatSession struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key"`
	SenderId      uuid.UUID  `gorm:"type:uuid"`
	ReceiverId    uuid.UUID  `gorm:"type:uuid"`
	IsInitialized bool       `gorm:"default:false"`
	EphemeralKey  string     `gorm:"type:varchar(255)"`
	DeletedAt     *time.Time `gorm:"type:time"`
	CreatedAt     time.Time  `gorm:"type:time;default:current_timestamp;not null"`
	Sender        *User      `gorm:"foreignKey:SenderId"`
	Receiver      *User      `gorm:"foreignKey:ReceiverId"`
}

type PendingMessage struct {
	ID             uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Type           string       `gorm:"type:varchar(255)"`
	Index          uint64       `gorm:"type:bigint"`
	OwnerId        uuid.UUID    `gorm:"type:uuid"`
	SenderId       uuid.UUID    `gorm:"type:uuid"`
	SenderUsername string       `gorm:"type:varchar(255)"`
	ChatSessionId  uuid.UUID    `gorm:"type:uuid"`
	CipherMessage  string       `gorm:"type:text"`
	PlainMessage   *string      `gorm:"type:text"`
	FilePath       *string      `gorm:"type:text"`
	IsBinary       bool         `gorm:"default:false"`
	IsRead         bool         `gorm:"default:false"`
	Owner          *User        `gorm:"foreignKey:OwnerId"`
	Sender         *User        `gorm:"foreignKey:SenderId"`
	ChatSession    *ChatSession `gorm:"foreignKey:ChatSessionId"`
	CreatedAt      time.Time    `gorm:"type:time;default:current_timestamp;not null"`
}

type UploadedFile struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Type      string    `gorm:"type:varchar(255)"`
	Size      uint64    `gorm:"type:bigint"`
	CreatedAt time.Time `gorm:"type:time;default:current_timestamp;not null"`
	OwnerId   uuid.UUID `gorm:"type:uuid"`
	Owner     *User     `gorm:"foreignKey:OwnerId"`
}
