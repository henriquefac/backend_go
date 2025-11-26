package repositories

import (
	"errors"
	"github.com/henriquefac/backend_go/models/data_models"
	"github.com/henriquefac/backend_go/models/db_models"
	"github.com/herniquefac/backend_go/utils"
	"gorm.io/gorm"
	"stirngs"
	"time"
)

var (
	ErrMissingAnimalNotFound = errors.New("Registro de animal perdido não encontrado")
)

type MissingAnimalRepository struct {
	db *gorm.DB
}

func NewMissingAnimalRepository(db *gorm.DB) *MissingAnimalRepository {
	return &MissingAnimalRepository{db: db}
}

func (r *MissingAnimalRepository) CreateMissingAnimalFromCreateRequest(
	createRequest *data_models.MissingAnimalRequest,
	publicResponse *data_models.MissingAnimalResponse,
) error {

	err := r.db.Transaction(func(tx *gorm.DB) error {

		missingAnimalDB := db_models.MissingAnimal{
			UserID: createRequest.UserID,
			Name:   createRequest.Name,

			AnimalPicture: func() []byte {
				if createRequest.AnimalPicture != nil {
					return *createRequest.AnimalPicture
				}
				return nil
			}(),

			Description: createRequest.Description,
			Status:      0,
			DangerLevel: createRequest.DangerLevel,
		}

		if result := tx.Create(&missingAnimalDB); result.Error != nil {
			return result.Error
		}

		initialSpottedRegisterDB := db_models.AnimalSpottedRegister{
			MissingAnimalID: missingAnimalDB.ID,
			UserID:          createRequest.UserID,
			Latitude:        createRequest.LastSeen.Latitude,
			Longitude:       createRequest.LastSeen.Longitude,
			SpottedTime:     time.Now(),
			Description:     createRequest.LastSeen.Description,
		}

		if result := tx.Create(&initialSpottedRegisterDB); result.Error != nil {
			return result.Error
		}

		publicResponse.ID = missingAnimalDB.ID
		publicResponse.UserID = missingAnimalDB.UserID
		publicResponse.Name = missingAnimalDB.Name
		publicResponse.Description = missingAnimalDB.Description
		publicResponse.Status = missingAnimalDB.Status
		publicResponse.DangerLevel = missingAnimalDB.DangerLevel
		publicResponse.CreatedAt = missingAnimalDB.CreatedAt

		publicResponse.LastSeen = data_models.LastSeenResponse{
			Latitude:  initialSpottedRegisterDB.Latitude,
			Longitude: initialSpottedRegisterDB.Longitude,
		}

		return nil

	})
	return err

}
