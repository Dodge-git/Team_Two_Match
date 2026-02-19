package service

import (
	"errors"
	"match_service/internal/dto"
	"match_service/internal/errs"
	"match_service/internal/models"
	"match_service/internal/repository"
	"strings"
)

type TeamService interface {
	Create(req dto.CreateTeamRequest) (*models.Team, error)
	GetByID(id uint) (*models.Team, error)
	Delete(id uint) error
	List(filter models.TeamFilter) ([]models.Team, error)
}

type teamService struct {
	teamRepo repository.TeamRepository
}

func NewTeamService(teamRepo repository.TeamRepository) TeamService {
	return &teamService{teamRepo: teamRepo}
}

func (s *teamService) Create(req dto.CreateTeamRequest) (*models.Team, error) {
	if req.SportID == 0 {
		return nil, errs.ErrInvalidTeamID
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errs.ErrInvalidTeamName
	}
	shortName := strings.TrimSpace(req.ShortName)
	if shortName == "" {
		return nil, errs.ErrInvalidTeamName
	}
	city := strings.TrimSpace(req.City)
	if city == "" {
		return nil, errs.ErrInvalidTeamName
	}

	team := &models.Team{
		SportID:   req.SportID,
		Name:      name,
		ShortName: shortName,
		City:      city,
	}
	if err := s.teamRepo.Create(team); err != nil {
		return nil, err
	}
	return team, nil
}

func (s *teamService) GetByID(id uint) (*models.Team, error) {
	if id == 0 {
		return nil, errs.ErrInvalidTeamID
	}
	team, err := s.teamRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, errs.ErrTeamNotFound) {
			return nil, errs.ErrTeamNotFound
		}
		return nil, err
	}
	return team, nil
}

func (s *teamService) Delete(id uint) error {
	if id == 0 {
		return errs.ErrInvalidTeamID
	}
	return s.teamRepo.Delete(id)
}

func (s *teamService) List(filter models.TeamFilter) ([]models.Team, error) {
	teams, _, err := s.teamRepo.List(filter)
	if err != nil {
		return nil, err
	}
	return teams, nil
}
