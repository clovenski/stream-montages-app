package models

type HighlightInfo struct {
	YTVideoID string
	Timestamp string
	Score     float64
}

type DbTopHighlightsInfo struct {
	Timestamp string
	Score     float64
}

type RawDbTopHighlightsResponse struct {
	TopHighlights string
	VideoId       string
}
