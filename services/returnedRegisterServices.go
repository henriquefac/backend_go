package services

import (
	"github.com/henriquefac/backend_go/models/data_models"
	"github.com/henriquefac/backend_go/repositories"

	"gorm.io/gorm"
)

// O serviço precisa de acesso ao repositório de retorno E ao repositório de MissingAnimal
type ReturnedRegisterService struct {
	returnedRegisterRepo *repositories.ReturnedRegisterRepository
	missingAnimalRepo    *repositories.MissingAnimalRepository // Necessário para atualizar o Status
}

func NewReturnedRegisterService(
	r *repositories.ReturnedRegisterRepository,
	m *repositories.MissingAnimalRepository,
) *ReturnedRegisterService {
	return &ReturnedRegisterService{
		returnedRegisterRepo: r,
		missingAnimalRepo:    m,
	}
}

func (s *ReturnedRegisterService) Create(createRequest *data_models.ReturnedRegisterCreateRequest) (
	*data_models.ReturnedRegisterResponse, error,
) {
	var publicResponse data_models.ReturnedRegisterResponse

	// 1. Inicia a Transação (usando o DB do Repositório de MissingAnimal)
	// Isso garante que tanto o registro quanto a atualização de status sejam atômicos.
	err := s.missingAnimalRepo.DB().Transaction(func(tx *gorm.DB) error {

		// --- 1. Persistência do Registro de Retorno ---

		// Cria uma cópia do repositório de retorno usando a Transação (tx)
		txReturnedRepo := repositories.NewReturnedRegisterRepository(tx)

		// Tenta criar o registro de retorno
		err := txReturnedRepo.CreateReturnedRegisterFromCreateRequest(createRequest, &publicResponse)
		if err != nil {
			return err // Se falhar (ex: já existe), faz Rollback
		}

		// --- 2. Regra de Negócio: Atualizar Status do Animal ---

		// Definir a mudança de status (1 = Encontrado/Devolvido)
		newStatus := 1

		// Cria o objeto de atualização
		updateRequest := data_models.MissingAnimalUpdateRequest{
			ID:     createRequest.MissingAnimalID,
			UserID: createRequest.RescuerID, // Usamos o RescuerID como o 'usuário autorizado' para fins de Service
			Status: &newStatus,
		}

		// O método UpdateMissingAnimalFromUpdateRequest precisa ser adaptado para receber o TX
		// Aqui, simulamos a chamada:

		txMissingRepo := repositories.NewMissingAnimalRepository(tx)

		// Chamamos o repositório de MissingAnimal para atualizar o status
		// Passamos 'nil' para a resposta, pois não precisamos dela aqui, apenas da atualização do DB
		err = txMissingRepo.UpdateStatusForReturnedAnimal(&updateRequest)
		if err != nil {
			return err // Se falhar a atualização do status, faz Rollback
		}

		return nil // Commit da transação se tudo ocorrer bem
	})

	if err != nil {
		return nil, err
	}

	return &publicResponse, nil
}

func (s *ReturnedRegisterService) GetReturnedRegisterByAnimalID(animalID uint) (
	*data_models.ReturnedRegisterResponse, error,
) {
	var response data_models.ReturnedRegisterResponse

	// O repositório retorna nil error se o registro não for encontrado (gorm.ErrRecordNotFound),
	// mas deixa 'response' com valores zero.
	err := s.returnedRegisterRepo.GetReturnedRegisterFromAnimal(animalID, &response)

	if err != nil {
		return nil, err
	}

	// Verifica se a struct foi populada. Se o ID for 0, o registro não foi encontrado.
	if response.ID == 0 {
		return nil, nil // Animal ainda não foi devolvido
	}

	return &response, nil
}
