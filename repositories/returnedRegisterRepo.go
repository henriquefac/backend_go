package repositories

import (
	"errors"
	"github.com/henriquefac/backend_go/models/data_models"
	"github.com/henriquefac/backend_go/models/db_models"
	"gorm.io/gorm"
	"time"
)

// Reutilização de erros
var (
	ErrReturnedRegisterExists = errors.New("O registro de retorno para este animal já existe")
)

// Estrutura do Repositório (Assumindo que está em outro arquivo ou você adicionará aqui)
type ReturnedRegisterRepository struct {
	db *gorm.DB
}

func NewReturnedRegisterRepository(db *gorm.DB) *ReturnedRegisterRepository {
	return &ReturnedRegisterRepository{db: db}
}

func (r *ReturnedRegisterRepository) CreateReturnedRegisterFromCreateRequest(
	createRequest *data_models.ReturnedRegisterCreateRequest,
	publicResponse *data_models.ReturnedRegisterResponse,
) error {
	var existingAnimal db_models.MissingAnimal
	// Busca otimizada
	result := r.db.Select("id").First(&existingAnimal, createRequest.MissingAnimalID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrMissingAnimalNotFound // O animal que está sendo devolvido não existe
		}
		return result.Error
	}

	returnedRegisterDB := db_models.AnimalReturnedRegister{
		MissingAnimalID: createRequest.MissingAnimalID,
		RescuerID:       createRequest.RescuerID,
		ReturnDate:      createRequest.ReturnDate,
	}
	if result := r.db.Create(&returnedRegisterDB); result.Error != nil {
		// Trata erro de violação de unicidade (se o animal já tiver um registro de retorno)
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) || result.RowsAffected == 0 {
			return ErrReturnedRegisterExists
		}
		return result.Error
	}

	publicResponse.ID = returnedRegisterDB.ID
	publicResponse.CreatedAt = returnedRegisterDB.CreatedAt
	publicResponse.MissingAnimalID = returnedRegisterDB.MissingAnimalID
	publicResponse.RescuerID = returnedRegisterDB.RescuerID
	publicResponse.ReturnDate = returnedRegisterDB.ReturnDate

	return nil
}

func (r *ReturnedRegisterRepository) GetReturnedRegisterFromAnimal (
	animalID uint, publicResponse *data_models.ReturnedRegisterResponse,
) error {
	var existingAnimal 
}
