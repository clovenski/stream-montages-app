package models

import (
	"time"

	"github.com/google/uuid"
)

type Montage struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	JobID     uuid.UUID `gorm:"type:uuid;not null"`
	Name      string    `gorm:"type:varchar(255);not null"`
	FilePath  string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
