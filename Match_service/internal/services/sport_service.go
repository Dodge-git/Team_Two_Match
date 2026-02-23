package services

import (
	"match_service/internal/dto"
	"match_service/internal/errs"
	"match_service/internal/models"
	"match_service/internal/repository"
)

type SportService interface {
	Create(req dto.CreateSportRequest) (*models.Sport, error)
	List() ([]models.Sport, error)
}

type sportService struct {
	sportRepo repository.SportRepository
}

func NewSportService(sportRepo repository.SportRepository) SportService {
	return &sportService{sportRepo: sportRepo}
}

func (s *sportService) Create(req dto.CreateSportRequest) (*models.Sport, error) {

	if _, ok := models.ValidSport[req.Name]; !ok {
		return nil, errs.ErrInvalidSport
	}

	sport := models.Sport{
		Name: req.Name,
	}

	if err := s.sportRepo.Create(&sport); err != nil {
		return nil, err
	}
	return &sport, nil
}

func (s *sportService) List() ([]models.Sport, error) {

	sports, err := s.sportRepo.List()
	if err != nil {
		return nil, err
	}
	return sports, nil
}
