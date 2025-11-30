package repositories

import (
	"errors"
	"github.com/henriquefac/backend_go/models/data_models"
	"github.com/henriquefac/backend_go/models/db_models"
	"gorm.io/gorm"
)

var (
	ErrSpottedRegisterNotFound = errors.New("Registro não encontrado")
)

type SpottedRegisterRepository struct {
	db *gorm.DB
}

func NewSpottedRegisterRepository(db *gorm.DB) *SpottedRegisterRepository {
	return &SpottedRegisterRepository{db: db}
}

// Criar registro de ultima localização do animal em questão
// precisa do ID do regstro de MissingAnimal

func (r *SpottedRegisterRepository) CreateSpottedRegisterFromCreateRequest(
	createRequest *data_models.SpottedRegisterCreateRequest,
	publicResponse *data_models.SpottedRegisterResponse,
) error {

	// 1. Verificar se o registro de MissingAnimal existe
	var existingAnimal db_models.MissingAnimal

	// Nota: O Select("id") é otimizado, só precisamos saber que ele existe.
	result := r.db.Select("id").First(&existingAnimal, createRequest.MissingAnimalID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrMissingAnimalNotFound // Retorna erro se o animal pai não existir
		}
		return result.Error
	}

	// 2. Mapear Requisição para o Modelo de Banco
	spottedRegisterDB := db_models.AnimalSpottedRegister{
		MissingAnimalID: createRequest.MissingAnimalID,
		UserID:          createRequest.UserID,

		// Trata o ponteiro opcional da imagem
		AnimalPicture: func() []byte {
			if createRequest.AnimalPicture != nil {
				return *createRequest.AnimalPicture
			}
			return nil
		}(),

		Latitude:    createRequest.Latitude,
		Longitude:   createRequest.Longitude,
		SpottedTime: createRequest.SpottedTime,
		Description: createRequest.Description,
	}

	// 3. Criar o registro
	if result := r.db.Create(&spottedRegisterDB); result.Error != nil {
		return result.Error
	}

	// 4. Popular a Resposta Pública
	publicResponse.ID = spottedRegisterDB.ID
	publicResponse.UserID = spottedRegisterDB.UserID
	publicResponse.MissingAnimalID = spottedRegisterDB.MissingAnimalID
	publicResponse.Latitude = spottedRegisterDB.Latitude
	publicResponse.Longitude = spottedRegisterDB.Longitude
	publicResponse.SpottedTime = spottedRegisterDB.SpottedTime
	publicResponse.Description = spottedRegisterDB.Description
	// publicResponse.AnimalPicture não está incluído no modelo de resposta,
	// seguindo a boa prática de retornar URLs e não bytes brutos.

	return nil
}

func (r *SpottedRegisterRepository) ListSpottedRegistersByAnimalID(
	animalID uint,
) ([]data_models.SpottedRegisterResponse, error) {

	var spottedRegistersDB []db_models.AnimalSpottedRegister

	// 1. Consulta ao Banco de Dados
	// Busca todos os registros onde MissingAnimalID corresponde ao ID do animal.
	result := r.db.Where("missing_animal_id = ?", animalID).
		// Ordenamos por SpottedTime DESC para que os avistamentos mais recentes apareçam primeiro.
		Order("spotted_time DESC").
		Find(&spottedRegistersDB)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Retorna uma lista vazia se nenhum registro for encontrado (sem erro)
			return []data_models.SpottedRegisterResponse{}, nil
		}
		return nil, result.Error
	}

	// 2. Mapeamento para a Resposta Pública
	var responses []data_models.SpottedRegisterResponse

	for _, spottedDB := range spottedRegistersDB {
		response := data_models.SpottedRegisterResponse{
			ID:              spottedDB.ID,
			UserID:          spottedDB.UserID,
			MissingAnimalID: spottedDB.MissingAnimalID,
			Latitude:        spottedDB.Latitude,
			Longitude:       spottedDB.Longitude,
			SpottedTime:     spottedDB.SpottedTime,
			Description:     spottedDB.Description,
			// AnimalPicture pode ser opcionalmente mapeado aqui, dependendo da necessidade
		}
		responses = append(responses, response)
	}

	return responses, nil
}
