package outbox

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/plugin/optimisticlock"
)

type Outbox struct {
	Id            uuid.UUID
	AggregateType string
	AggregateId   string
	EventType     string
	Payload       []byte
	CreatedAt     time.Time
}

func (Outbox) TableName() string {
	return "outbox"
}

type OutboxLock struct {
	Id          int
	Locked      bool
	LockedAt    *time.Time
	LockedUntil *time.Time
	Version     optimisticlock.Version
}

func (OutboxLock) TableName() string {
	return "outbox_lock"
}
