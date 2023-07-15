package config

import (
	"log"

	"github.com/joeshaw/envdecode"
)

type Config struct {
	ChannelIds                 string `env:"SM_CHANNEL_IDS,required"`
	StreamersTableName         string `env:"SM_STREAMERS_TABLE_NAME,required"`
	StreamersTablePKey         string `env:"SM_STREAMERS_TABLE_PKEY,required"`
	StreamersTableNameAttr     string `env:"SM_STREAMERS_TABLE_NAME_ATTR,required"`
	HighlightsTableName        string `env:"SM_HIGHLIGHTS_TABLE_NAME,required"`
	HighlightsTablePKey        string `env:"SM_HIGHLIGHTS_TABLE_PKEY,required"`
	HighlightsTableSKey        string `env:"SM_HIGHLIGHTS_TABLE_SKEY,required"`
	HighlightsTableTopAttr     string `env:"SM_HIGHLIGHTS_TABLE_TOP_ATTR,required"`
	HighlightsTableVideoIdAttr string `env:"SM_HIGHLIGHTS_TABLE_VIDEO_ID_ATTR,required"`

	JobsRepoUrl string `env:"SM_JOBS_REPO_URL,required"`
}

func LoadConfig() (config Config, err error) {
	log.Println("Loading configurations...")

	err = envdecode.Decode(&config)
	return
}
