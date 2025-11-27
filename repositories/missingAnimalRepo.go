package repositories

import (
	"errors"
	"github.com/henriquefac/backend_go/models/data_models"
	"github.com/henriquefac/backend_go/models/db_models"
	"gorm.io/gorm"
	"time"
)

var (
	ErrMissingAnimalNotFound = errors.New("Registro de animal perdido não encontrado")
	ErrUnauthorized          = errors.New("Usuário não autorizado a alterar este registro")
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

		existingAnimal := db_models.MissingAnimal{
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

		if result := tx.Create(&existingAnimal); result.Error != nil {
			return result.Error
		}

		initialSpottedRegisterDB := db_models.AnimalSpottedRegister{
			MissingAnimalID: existingAnimal.ID,
			UserID:          createRequest.UserID,
			Latitude:        createRequest.LastSeen.Latitude,
			Longitude:       createRequest.LastSeen.Longitude,
			SpottedTime:     time.Now(),
			Description:     createRequest.LastSeen.Description,
		}

		if result := tx.Create(&initialSpottedRegisterDB); result.Error != nil {
			return result.Error
		}

		publicResponse.ID = existingAnimal.ID
		publicResponse.UserID = existingAnimal.UserID
		publicResponse.Name = existingAnimal.Name
		publicResponse.Description = existingAnimal.Description
		publicResponse.Status = existingAnimal.Status
		publicResponse.DangerLevel = existingAnimal.DangerLevel
		publicResponse.CreatedAt = existingAnimal.CreatedAt

		publicResponse.LastSeen = data_models.LastSeenResponse{
			Latitude:  initialSpottedRegisterDB.Latitude,
			Longitude: initialSpottedRegisterDB.Longitude,
		}

		return nil

	})
	return err
}

// alterar um registro de missing animal como usuário
// autoridade para alterar e editar informações

func (r *MissingAnimalRepository) UpdateMissingAnimalFromEditRequest(
	updateRequest *data_models.MissingAnimalUpdateRequest,
	publicResponse *data_models.MissingAnimalResponse,
) error {

	var existingAnimal db_models.MissingAnimal
	result := r.db.Select("id").First(&existingAnimal, updateRequest.ID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrMissingAnimalNotFound
		}
		return result.Error
	}

	if existingAnimal.UserID != updateRequest.UserID {
		return ErrUnauthorized
	}

	updates := make(map[string]interface{})

	if updateRequest.Name != nil {
		updates["Name"] = *updateRequest.Name
	}
	if updateRequest.Description != nil {
		updates["Description"] = *updateRequest.Description
	}
	if updateRequest.DangerLevel != nil {
		updates["DangerLevel"] = *updateRequest.DangerLevel
	}
	if updateRequest.Status != nil {
		updates["Status"] = *updateRequest.Status
	}
	if updateRequest.AnimalPicture != nil {
		updates["AnimalPicture"] = *updateRequest.AnimalPicture
	}

	// 4. Executar a Atualização
	result = r.db.Model(&db_models.MissingAnimal{}).
		Where("id = ?", updateRequest.ID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	var finalAnimalDB db_models.MissingAnimal

	result = r.db.Preload("SpottedRegister", func(db *gorm.DB) *gorm.DB {
		// Ordenamos por CreatedAt ou ID para garantir que pegamos o registro inicial (o mais antigo)
		return db.Order("created_at ASC").Limit(1)
	}).First(&finalAnimalDB, updateRequest.ID)

	// Localização inicial/última vista (SpottedRegister)
	if len(finalAnimalDB.SpottedRegister) > 0 {
		firstSpotted := finalAnimalDB.SpottedRegister[0]
		publicResponse.LastSeen = data_models.LastSeenResponse{
			Latitude:  firstSpotted.Latitude,
			Longitude: firstSpotted.Longitude,
			// CORREÇÃO: Popula o campo Description
			Description: firstSpotted.Description,
		}
	} else {
		return errors.New("Registro do animal sem localização inicial")
	}

	return nil
}
