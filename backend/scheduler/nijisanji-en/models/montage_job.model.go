package models

import "gorm.io/datatypes"

type MontageJobHighlightInfo struct {
	VideoURL        string
	Timestamp       string
	DurationSeconds uint32
}

type MontageJobDefinition struct {
	Highlights datatypes.JSONSlice[MontageJobHighlightInfo]
}

type MontageJobRequest struct {
	Name          string
	JobDefinition datatypes.JSONType[MontageJobDefinition]
}
