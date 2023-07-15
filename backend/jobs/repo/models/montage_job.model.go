package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type HighlightInfo struct {
	VideoURL        string
	Timestamp       string
	DurationSeconds uint32
}

type MontageJobDefinition struct {
	Highlights datatypes.JSONSlice[HighlightInfo]
}

type MontageJob struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key"`
	Name          string    `gorm:"type:varchar(255);not null"`
	JobDefinition datatypes.JSONType[MontageJobDefinition]
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
