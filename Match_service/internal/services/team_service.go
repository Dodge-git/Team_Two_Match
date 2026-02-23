package services

import (
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
	List(filter models.TeamFilter) ([]models.Team, int64, error)
}

type teamService struct {
	teamRepo  repository.TeamRepository
	sportRepo repository.SportRepository
}

func NewTeamService(teamRepo repository.TeamRepository, sportRepo repository.SportRepository) TeamService {
	return &teamService{teamRepo: teamRepo, sportRepo: sportRepo}
}

func (s *teamService) Create(req dto.CreateTeamRequest) (*models.Team, error) {
	if req.SportID == 0 {
		return nil, errs.ErrInvalidSportID
	}

	if _, err := s.sportRepo.GetByID(req.SportID); err != nil {
		return nil, err
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
		return nil, errs.ErrInvalidCity
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

func (s *teamService) List(filter models.TeamFilter) ([]models.Team, int64, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}

	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	if filter.SportID != nil {
		if *filter.SportID == 0 {
			return nil, 0, errs.ErrInvalidSportID
		}
	}
	teams, total, err := s.teamRepo.List(filter)
	if err != nil {
		return nil, 0, err
	}
	return teams, total, nil
}
