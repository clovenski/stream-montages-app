package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/clovenski/stream-montages-app/backend/jobs/repo/models"
	"github.com/clovenski/stream-montages-app/backend/jobs/repo/services"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type MontageJobRepository interface {
	Create(job models.MontageJob) (models.MontageJob, error)
	Exists(id uuid.UUID) bool
	Read(id uuid.UUID) models.MontageJob
	ReadAll() []models.MontageJob
	Update(job models.MontageJob) (models.MontageJob, bool, error)
	Delete(id uuid.UUID)
}

var service MontageJobRepository

func InitController(client *gorm.DB) {
	service = services.MontageJobService{DBClient: client}
}

func GetMontageJobs(w http.ResponseWriter, r *http.Request) {
	jobs := service.ReadAll()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jobs)
}

func CreateMontageJob(w http.ResponseWriter, r *http.Request) {
	var job models.MontageJob
	json.NewDecoder(r.Body).Decode(&job)

	savedJob, err := service.Create(job)
	if err != nil {
		log.Printf("Failed to create job with id %v Error: %v\n", job.ID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(savedJob)
}

func GetMontageJobByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	jobId, err := uuid.Parse(id)

	if err != nil {
		log.Printf("Failed to parse id %s into uuid. Error: %v\n", id, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	job := service.Read(jobId)

	if job.ID == uuid.Nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(job)
}

func UpdateMontageJobByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	jobId, err := uuid.Parse(id)

	if err != nil {
		log.Printf("Failed to parse id %s into uuid. Error: %v\n", id, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var job models.MontageJob
	json.NewDecoder(r.Body).Decode(&job)

	if jobId != job.ID {
		log.Printf("Mismatched job id between query param and request body. Param=%v Request=%v\n", jobId, job.ID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedJob, inserted, err := service.Update(job)
	if err != nil {
		log.Printf("Failed to update job with id %v Error: %v\n", jobId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	status := http.StatusOK
	if inserted {
		status = http.StatusCreated
	}
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(updatedJob)
}

func DeleteMontageJobByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	jobId, err := uuid.Parse(id)

	if err != nil {
		log.Printf("Failed to parse id %s into uuid. Error: %v\n", id, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	service.Delete(jobId)
	w.WriteHeader(http.StatusNoContent)
}
