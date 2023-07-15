package services

import (
	"errors"

	"github.com/clovenski/stream-montages-app/backend/jobs/repo/models"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MontageJobService struct {
	DBClient *gorm.DB
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
		// TODO: kafka publish
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
