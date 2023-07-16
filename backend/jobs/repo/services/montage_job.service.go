package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/clovenski/stream-montages-app/backend/jobs/repo/config"
	"github.com/clovenski/stream-montages-app/backend/jobs/repo/models"
	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MontageJobService struct {
	DBClient *gorm.DB
	Config   config.Config
}

func validateJobDefinition(def datatypes.JSONType[models.MontageJobDefinition]) error {
	highlights := def.Data().Highlights
	if len(highlights) == 0 {
		return errors.New("Highlights must be non-empty")
	}

	return nil
}

func (svc MontageJobService) Create(job models.MontageJob) (models.MontageJob, error) {
	if job.ID == uuid.Nil {
		job.ID = uuid.New()
	}

	if err := validateJobDefinition(job.JobDefinition); err != nil {
		return models.MontageJob{}, err
	}

	svc.DBClient.Create(&job)

	go func() {
		// kafka publish
		dialer := kafka.DefaultDialer
		dialer.ClientID = svc.Config.ClientId
		conn, err := dialer.DialLeader(context.Background(), "tcp", svc.Config.BootstrapServer, svc.Config.JobsTopic, 0)
		if err != nil {
			log.Printf("Job id=%s - Failed to dial leader: %v\n", job.ID, err)
			return
		}

		jobJson, err := json.Marshal(job) // # may not be up to date with db (i.e. hidden manip to record after saving to db)
		if err != nil {
			log.Printf("Job id=%s - Failed to marshal record to json: %v\n", job.ID, err)
			return
		}

		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		_, err = conn.Write(jobJson)
		if err != nil {
			log.Printf("Job id=%s - Failed to write messages: %v\n", job.ID, err)
			return
		}

		log.Printf("Job id=%s - Successfully published to kafka topic %s\n", job.ID, svc.Config.JobsTopic)

		if err := conn.Close(); err != nil {
			log.Printf("Job id=%s - Failed to close writer: %v\n", job.ID, err)
			return
		}
	}()

	return job, nil
}

func (svc MontageJobService) Exists(id uuid.UUID) bool {
	var job models.MontageJob
	svc.DBClient.First(&job, id)
	return job.ID != uuid.Nil
}

func (svc MontageJobService) Read(id uuid.UUID) (job models.MontageJob) {
	svc.DBClient.First(&job, id)
	return
}

func (svc MontageJobService) ReadAll() (jobs []models.MontageJob) {
	svc.DBClient.Find(&jobs)
	return
}

func (svc MontageJobService) Update(job models.MontageJob) (_ models.MontageJob, inserted bool, err error) {
	if !svc.Exists(job.ID) {
		inserted = true
	}

	if err := validateJobDefinition(job.JobDefinition); err != nil {
		return models.MontageJob{}, false, err
	}

	svc.DBClient.Save(job)
	return job, inserted, nil
}

func (svc MontageJobService) Delete(id uuid.UUID) {
	svc.DBClient.Delete(&models.MontageJob{}, id)
	return
}
