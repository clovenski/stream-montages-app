package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/clovenski/stream-montages-app/backend/montages/repo/models"
	"github.com/clovenski/stream-montages-app/backend/montages/repo/services"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type MontageRepository interface {
	Create(montage models.Montage) (models.Montage, error)
	Exists(id uuid.UUID) bool
	Read(id uuid.UUID) models.Montage
	ReadAll() []models.Montage
	Update(montage models.Montage) (models.Montage, bool, error)
	Delete(id uuid.UUID)
}

var service MontageRepository

func InitController(client *gorm.DB) {
	service = services.MontageService{DBClient: client}
}

func GetMontages(w http.ResponseWriter, r *http.Request) {
	montages := service.ReadAll()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(montages)
}

func CreateMontage(w http.ResponseWriter, r *http.Request) {
	var montage models.Montage
	json.NewDecoder(r.Body).Decode(&montage)

	savedMontage, err := service.Create(montage)
	if err != nil {
		log.Printf("Failed to create montage with id %v Error: %v\n", montage.ID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(savedMontage)
}

func GetMontageByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	montageId, err := uuid.Parse(id)

	if err != nil {
		log.Printf("Failed to parse id %s into uuid. Error: %v\n", id, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	montage := service.Read(montageId)

	if montage.ID == uuid.Nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(montage)
}

func UpdateMontageByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	montageId, err := uuid.Parse(id)

	if err != nil {
		log.Printf("Failed to parse id %s into uuid. Error: %v\n", id, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var montage models.Montage
	json.NewDecoder(r.Body).Decode(&montage)

	if montageId != montage.ID {
		log.Printf("Mismatched montage id between query param and request body. Param=%v Request=%v\n", montageId, montage.ID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedMontage, inserted, err := service.Update(montage)
	if err != nil {
		log.Printf("Failed to update montage with id %v Error: %v\n", montageId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	status := http.StatusOK
	if inserted {
		status = http.StatusCreated
	}
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(updatedMontage)
}

func DeleteMontageByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	montageId, err := uuid.Parse(id)

	if err != nil {
		log.Printf("Failed to parse id %s into uuid. Error: %v\n", id, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	service.Delete(montageId)
	w.WriteHeader(http.StatusNoContent)
}
