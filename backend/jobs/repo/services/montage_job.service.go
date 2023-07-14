package services

import (
	"github.com/clovenski/stream-montages-app/backend/jobs/repo/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MontageJobService struct {
	DBClient *gorm.DB
}

func (svc MontageJobService) Create(job models.MontageJob) models.MontageJob {
	if job.ID == uuid.Nil {
		job.ID = uuid.New()
	}

	svc.DBClient.Create(&job)

	go func() {
		// TODO: kafka publish
	}()

	return job
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

func (svc MontageJobService) Update(job models.MontageJob) (_ models.MontageJob, inserted bool) {
	if !svc.Exists(job.ID) {
		inserted = true
	}
	svc.DBClient.Save(job)
	return job, inserted
}

func (svc MontageJobService) Delete(id uuid.UUID) {
	svc.DBClient.Delete(&models.MontageJob{}, id)
	return
}
