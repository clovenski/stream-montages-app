package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/clovenski/stream-montages-app/backend/scheduler/nijisanji-en/models"
)

type SchedulerConfig struct {
	JobsRepoUrl string
}

type MontageJobScheduler struct {
	Config SchedulerConfig
}

func (svc MontageJobScheduler) Schedule(request models.MontageJobRequest) error {
	log.Printf("Scheduling job according to request %v\n", request)
	jobJsonReq, err := json.Marshal(request)
	if err != nil {
		log.Printf("Failed to marshal the streamer montage job '%s'. Error: %v\n", request.Name, err)
		return err
	}
	resp, err := http.DefaultClient.Post(svc.Config.JobsRepoUrl, "application/json", strings.NewReader(string(jobJsonReq)))
	if err != nil {
		log.Printf("Failed to post to schedule streamer montage job '%s'. Error: %v\n", request.Name, err)
		return err
	}
	defer resp.Body.Close()
	respStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read scheduling montage job '%s' response body. Error: %v\n", request.Name, err)
		return err
	}
	log.Printf("Response from job '%s' being scheduled: status=%s body=%s\n", request.Name, resp.Status, string(respStr))
	return nil
}
