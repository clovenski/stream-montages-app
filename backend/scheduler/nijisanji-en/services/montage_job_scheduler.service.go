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

func (svc MontageJobScheduler) Schedule(request models.MontageJobRequest) (response models.MontageJobResponse, err error) {
	log.Printf("Scheduling job according to request %v\n", request)
	var jobJsonReq []byte
	jobJsonReq, err = json.Marshal(request)
	if err != nil {
		log.Printf("Failed to marshal the streamer montage job '%s'. Error: %v\n", request.Name, err)
		return
	}
	var resp *http.Response
	resp, err = http.DefaultClient.Post(svc.Config.JobsRepoUrl, "application/json", strings.NewReader(string(jobJsonReq)))
	if err != nil {
		log.Printf("Failed to post to schedule streamer montage job '%s'. Error: %v\n", request.Name, err)
		return
	}
	defer resp.Body.Close()
	var respStr []byte
	respStr, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read scheduling montage job '%s' response body. Error: %v\n", request.Name, err)
		return
	}
	log.Printf("Response from job '%s' being scheduled: status=%s body=%s\n", request.Name, resp.Status, string(respStr))

	if err = json.Unmarshal(respStr, &response); err != nil {
		log.Printf("Failed to unmarshal response string to DTO. Error: %v\n", err)
	}
	return
}
