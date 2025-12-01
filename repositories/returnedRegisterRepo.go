package repositories

import (
	"errors"
	"github.com/henriquefac/backend_go/models/data_models"
	"github.com/henriquefac/backend_go/models/db_models"
	"gorm.io/gorm"
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

func (r *ReturnedRegisterRepository) GetReturnedRegisterFromAnimal(
	animalID uint,
	publicResponse *data_models.ReturnedRegisterResponse,
) error {

	var returnedRegister db_models.AnimalReturnedRegister

	// 1. CORREÇÃO CRÍTICA: Usar .First() para buscar um único registro
	result := r.db.Where("missing_animal_id = ?", animalID).First(&returnedRegister)

	// 2. Tratamento de Erro
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Se o registro NÃO for encontrado, retornamos nil (sem erro),
			// indicando que o animal ainda não foi devolvido.
			// 'publicResponse' permanece zerado (ID=0), o que o Service verifica.
			return nil
		}
		// Se for outro erro de sistema, retornamos o erro.
		return result.Error
	}

	// 3. Mapeamento dos dados encontrados para o publicResponse
	// O registro foi encontrado, então populamos a resposta.
	publicResponse.ID = returnedRegister.ID
	publicResponse.CreatedAt = returnedRegister.CreatedAt
	publicResponse.MissingAnimalID = returnedRegister.MissingAnimalID
	publicResponse.RescuerID = returnedRegister.RescuerID
	publicResponse.ReturnDate = returnedRegister.ReturnDate

	return nil
}

func (r *ReturnedRegisterRepository) ListReturnedRegisterFromUser(
	userID uint,
) (*[]data_models.ReturnedRegisterResponse, error) {
	var returndeRegistersDB []db_models.AnimalReturnedRegister

	result := r.db.Where("rescuer_id = ?", userID).Find(&returndeRegistersDB)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &[]data_models.ReturnedRegisterResponse{}, nil
		}
		return nil, result.Error
	}

	var responses []data_models.ReturnedRegisterResponse

	for _, returnedDB := range returndeRegistersDB {
		response := data_models.ReturnedRegisterResponse{
			ID:              returnedDB.ID,
			RescuerID:       returnedDB.RescuerID,
			MissingAnimalID: returnedDB.MissingAnimalID,
			ReturnDate:      returnedDB.ReturnDate,
		}
		responses = append(responses, response)
	}

	return &responses, nil
}
