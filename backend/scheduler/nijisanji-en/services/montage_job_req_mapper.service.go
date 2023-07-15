package services

import (
	"fmt"

	"github.com/clovenski/stream-montages-app/backend/scheduler/nijisanji-en/models"
	"gorm.io/datatypes"
)

type MontageJobRequestMapper struct{}

func (MontageJobRequestMapper) Map(highlights []models.HighlightInfo, jobName string) models.MontageJobRequest {
	reqHighlights := make([]models.MontageJobHighlightInfo, 0, len(highlights))
	for _, item := range highlights {
		reqHighlights = append(reqHighlights, models.MontageJobHighlightInfo{
			VideoURL:        fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.YTVideoID),
			Timestamp:       item.Timestamp,
			DurationSeconds: 60,
		})
	}

	return models.MontageJobRequest{
		Name: jobName,
		JobDefinition: datatypes.NewJSONType(models.MontageJobDefinition{
			Highlights: datatypes.NewJSONSlice(reqHighlights),
		}),
	}
}
