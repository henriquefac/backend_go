package services

import (
	"github.com/henriquefac/backend_go/models/data_models"
	"github.com/henriquefac/backend_go/repositories"
)

type SpottedRegisterService struct {
	spottedRegisterRepo *repositories.SpottedRegisterRepository
}

func NewSpottedRegisterService(s *repositories.SpottedRegisterRepository) *SpottedRegisterService {
	return &SpottedRegisterService{spottedRegisterRepo: s}
}

func (s *SpottedRegisterService) Create(spottedRegisterCreateRequest *data_models.SpottedRegisterCreateRequest) (
	*data_models.SpottedRegisterResponse, error,
) {
	// 1. Inicializa a struct que receberá a resposta pública
	var publicResponse data_models.SpottedRegisterResponse

	// 2. Chama o método de criação do Repositório
	// O Repositório cuida da validação de integridade (se MissingAnimal existe) e persistência.
	err := s.spottedRegisterRepo.CreateSpottedRegisterFromCreateRequest(
		spottedRegisterCreateRequest,
		&publicResponse, // Passamos o endereço para o repositório preencher
	)

	if err != nil {
		// Retorna o erro vindo da camada de Repositório (ex: MissingAnimalNotFound)
		return nil, err
	}

	// 3. Retorna o endereço da resposta pública preenchida
	return &publicResponse, nil
}
