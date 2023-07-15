package services

import (
	"errors"

	"github.com/clovenski/stream-montages-app/backend/montages/repo/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MontageService struct {
	DBClient *gorm.DB
}

func validateMontageDetails(montage models.Montage) error {
	if montage.JobID == uuid.Nil {
		return errors.New("Job ID must be filled")
	}
	if montage.Name == "" {
		return errors.New("Job name must be filled")
	}
	if montage.Filename == "" {
		return errors.New("Job name must be filled")
	}

	return nil
}

func (svc MontageService) Create(montage models.Montage) (models.Montage, error) {
	if montage.ID == uuid.Nil {
		montage.ID = uuid.New()
	}

	if err := validateMontageDetails(montage); err != nil {
		return models.Montage{}, err
	}

	svc.DBClient.Create(&montage)

	return montage, nil
}

func (svc MontageService) Exists(id uuid.UUID) bool {
	var montage models.Montage
	svc.DBClient.First(&montage, id)
	return montage.ID != uuid.Nil
}

func (svc MontageService) Read(id uuid.UUID) (montage models.Montage) {
	svc.DBClient.First(&montage, id)
	return
}

func (svc MontageService) ReadAll() (montages []models.Montage) {
	svc.DBClient.Find(&montages)
	return
}

func (svc MontageService) Update(montage models.Montage) (_ models.Montage, inserted bool, err error) {
	if !svc.Exists(montage.ID) {
		inserted = true
	}

	if err := validateMontageDetails(montage); err != nil {
		return models.Montage{}, false, err
	}

	svc.DBClient.Save(montage)
	return montage, inserted, nil
}

func (svc MontageService) Delete(id uuid.UUID) {
	svc.DBClient.Delete(&models.Montage{}, id)
	return
}
