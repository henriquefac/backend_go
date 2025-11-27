package services

import (
	"github.com/henriquefac/backend_go/models/data_models"
	"github.com/henriquefac/backend_go/repositories"
)

// Serviços associados a MissingAnimal

type MissingAnimalService struct {
	missingAnimalRepo *repositories.MissingAnimalRepository
}

func NewMissingAnimalService(msr *repositories.MissingAnimalRepository) *MissingAnimalService {
	return &MissingAnimalService{missingAnimalRepo: msr}
}

func (msr *MissingAnimalService) Create(missingAnimalRequest *data_models.MissingAnimalRequest) (
	*data_models.MissingAnimalResponse, error,
) {
	var publicResponse data_models.MissingAnimalResponse

	err := msr.missingAnimalRepo.CreateMissingAnimalFromCreateRequest(missingAnimalRequest, &publicResponse)

	if err != nil {
		return nil, err
	}

	return &publicResponse, nil
}

func (msr *MissingAnimalService) Update(missingAnimalUpdateRequest)
