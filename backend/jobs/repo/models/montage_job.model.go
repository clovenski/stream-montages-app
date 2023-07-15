package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type MontageJobStatus string

const (
	Pending  MontageJobStatus = "PENDING"
	Started  MontageJobStatus = "STARTED"
	Complete MontageJobStatus = "COMPLETE"
	Failed   MontageJobStatus = "FAILED"
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
	ID            uuid.UUID        `gorm:"type:uuid;primary_key"`
	Name          string           `gorm:"type:varchar(255);not null"`
	Status        MontageJobStatus `gorm:"type:varchar(32);not null;default:PENDING"`
	JobDefinition datatypes.JSONType[MontageJobDefinition]
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
