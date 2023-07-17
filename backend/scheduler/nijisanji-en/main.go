package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/clovenski/stream-montages-app/backend/scheduler/nijisanji-en/config"
	"github.com/clovenski/stream-montages-app/backend/scheduler/nijisanji-en/models"
	"github.com/clovenski/stream-montages-app/backend/scheduler/nijisanji-en/services"
)

type RequestDetails struct {
	StartFilter string `json:"startFilter"`
	EndFilter   string `json:"endFilter"`
}

func isValidISODate(date string) bool {
	re := regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`)
	return re.MatchString(date)
}

func parseRequest(request json.RawMessage) (startFilter, endFilter string, err error) {
	var details RequestDetails
	if err = json.Unmarshal(request, &details); err != nil {
		return
	}

	if !isValidISODate(details.StartFilter) {
		err = errors.New("Invalid start filter " + details.StartFilter)
		return
	}
	if !isValidISODate(details.EndFilter) {
		err = errors.New("Invalid end filter " + details.EndFilter)
		return
	}

	return details.StartFilter, details.EndFilter, nil
}

func handler(ctx context.Context, event events.CloudWatchEvent) (string, error) {
	// TODO: decompose this large function
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	scheduler := services.MontageJobScheduler{
		Config: services.SchedulerConfig{JobsRepoUrl: config.JobsRepoUrl},
	}
	mapper := services.MontageJobRequestMapper{}

	var startFilter, endFilter string
	startFilter, endFilter, err = parseRequest(event.Detail)
	if err != nil {
		log.Printf("Failed to parse event details. Error: %v\n", err)

		// default 1 week window, can come from request details if needed
		start := time.Now().AddDate(0, 0, -7)
		startFilter = strings.Split(start.Local().String(), " ")[0]    // only need the date part; local is OK as long as lambda runs in same timezone as db
		endFilter = strings.Split(time.Now().Local().String(), " ")[0] // only need the date part
	}

	log.Printf("Scheduling montage jobs for streams from %s to %s\n", startFilter, endFilter)

	dbClient := dynamodb.New(session.New())

	channelIds := strings.Split(config.ChannelIds, " ")
	topHighlights := make([]models.HighlightInfo, 0, len(channelIds))

	scheduledJobIds := make([]string, 0, len(channelIds)+1)

	for _, channelId := range channelIds {
		log.Printf("Getting streamer info for channel id %s\n", channelId)
		// get channel info
		streamerInfoProjExp := "#n"
		streamerInfoReq := &dynamodb.GetItemInput{
			TableName:            &config.StreamersTableName,
			ProjectionExpression: &streamerInfoProjExp,
			ExpressionAttributeNames: map[string]*string{
				"#n": &config.StreamersTableNameAttr,
			},
			Key: map[string]*dynamodb.AttributeValue{
				(config.StreamersTablePKey): {S: &channelId},
			},
		}
		streamerInfo, err := dbClient.GetItem(streamerInfoReq)
		if err != nil {
			log.Printf("Failed to get streamer info for channel id %s Error: %v\n", channelId, err)
			continue
		}
		streamerName := *streamerInfo.Item[config.StreamersTableNameAttr].S
		log.Printf("Getting past highlights for streamer %s\n", streamerName)

		// get past highlights for this streamer
		topHighlight := models.HighlightInfo{}
		var allHighlights []models.HighlightInfo
		var keyCondExp = "#PKey = :channelId AND #SKey BETWEEN :startFilter AND :endFilter"
		var lastEvalKey map[string]*dynamodb.AttributeValue
		var queryReq *dynamodb.QueryInput
		var projExp = strings.Join([]string{config.HighlightsTableTopAttr, config.HighlightsTableVideoIdAttr}, ",")
		for queryReq == nil || len(lastEvalKey) > 0 {
			queryReq = &dynamodb.QueryInput{
				ExclusiveStartKey:      lastEvalKey,
				ProjectionExpression:   &projExp,
				TableName:              &config.HighlightsTableName,
				KeyConditionExpression: &keyCondExp,
				ExpressionAttributeNames: map[string]*string{
					"#PKey": &config.HighlightsTablePKey,
					"#SKey": &config.HighlightsTableSKey,
				},
				ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
					":channelId":   {S: &channelId},
					":startFilter": {S: &startFilter},
					":endFilter":   {S: &endFilter},
				},
			}
			queryRes, err := dbClient.Query(queryReq)
			if err != nil {
				log.Printf("Failed to query highlights for channel id %s Error: %v\n", channelId, err)
				lastEvalKey = nil
				break
			}
			if *queryRes.Count == 0 {
				log.Printf("No highlights found for channel id %s\n", channelId)
				break
			}

			lastEvalKey = queryRes.LastEvaluatedKey
			if allHighlights == nil {
				allHighlights = make([]models.HighlightInfo, 0, *queryRes.Count)
			}
			results := []models.RawDbTopHighlightsResponse{}

			if err := dynamodbattribute.UnmarshalListOfMaps(queryRes.Items, &results); err != nil {
				log.Printf("Failed to unmarshal query result for channel id %s Error: %v\n", channelId, err)
				break
			}

			for _, result := range results {

				var topHighlightsArr = make([]models.DbTopHighlightsInfo, 0, 3)
				if err := json.Unmarshal([]byte(result.TopHighlights), &topHighlightsArr); err != nil {
					log.Printf("Failed to unmarshal for channel id %s the json string: %s", channelId, result.TopHighlights)
					break
				}

				transformed := models.HighlightInfo{
					YTVideoID: result.VideoId,
					Timestamp: topHighlightsArr[0].Timestamp, // assuming first elem is top-scoring
					Score:     topHighlightsArr[0].Score,     // assuming first elem is top-scoring
				}

				if topHighlight.Score < transformed.Score {
					topHighlight = transformed
				}

				allHighlights = append(allHighlights, transformed)
			}
		}

		// might be zero-value if no highlights found for this streamer
		if topHighlight.Score == 0 {
			continue
		}

		topHighlights = append(topHighlights, topHighlight)

		// schedule for streamer montage
		job, err := scheduler.Schedule(mapper.Map(allHighlights, fmt.Sprintf("%s montage from %v to %v", streamerName, startFilter, endFilter)))
		if err != nil {
			log.Printf("Failed to schedule streamer montage for channel id %s Error: %v\n", channelId, err)
			continue
		}

		scheduledJobIds = append(scheduledJobIds, job.ID)
	}

	// randomize order so not same order of streamers every time
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(topHighlights), func(i, j int) { topHighlights[i], topHighlights[j] = topHighlights[j], topHighlights[i] })

	// schedule for all streamers montage
	job, err := scheduler.Schedule(mapper.Map(topHighlights, fmt.Sprintf("All streamer montage from %v to %v", startFilter, endFilter)))
	if err != nil {
		log.Printf("Failed to schedule all streamers montage. Error: %v\n", err)
	} else {
		scheduledJobIds = append(scheduledJobIds, job.ID)
	}

	return strings.Join(scheduledJobIds, ","), nil
}

func main() {
	lambda.Start(handler)
}
