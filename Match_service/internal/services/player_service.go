package services

import (
	"errors"
	"match_service/internal/dto"
	"match_service/internal/errs"
	"match_service/internal/models"
	"match_service/internal/repository"
	"strings"
)

type PlayerService interface {
	Create(req dto.CreatePlayerRequest) (*models.Player, error)
	Update(playerID uint, req dto.UpdatePlayerRequest) (*models.Player, error)
	Delete(playerID uint) error
	GetByID(playerID uint) (*models.Player, error)
	List(teamID uint) ([]models.Player, error)
}

type playerService struct {
	playerRepo repository.PlayerRepository
}

func NewPlayerService(playerRepo repository.PlayerRepository) PlayerService {
	return &playerService{playerRepo: playerRepo}
}

func (s *playerService) Create(req dto.CreatePlayerRequest) (*models.Player, error) {
	name := strings.TrimSpace(req.Name)
	position := strings.TrimSpace(req.Position)

	player := &models.Player{
		TeamID:   req.TeamID,
		Name:     name,
		Number:   req.Number,
		Position: position,
	}

	if err := s.playerRepo.Create(player); err != nil {
		return nil, err
	}
	return player, nil
}

func (s *playerService) Update(playerID uint, req dto.UpdatePlayerRequest) (*models.Player, error) {
	if playerID == 0 {
		return nil, errs.ErrInvalidPlayerID
	}

	player, err := s.playerRepo.GetByID(playerID)
	if err != nil {
		if errors.Is(err, errs.ErrPlayerNotFound) {
			return nil, errs.ErrPlayerNotFound
		}
		return nil, err
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		player.Name = name
	}
	if req.Number != nil {
		player.Number = *req.Number
	}
	if req.Position != nil {
		position := strings.TrimSpace(*req.Position)
		player.Position = position
	}

	if err := s.playerRepo.Update(player); err != nil {
		return nil, err
	}
	return player, nil
}

func (s *playerService) Delete(id uint) error {
	if id == 0 {
		return errs.ErrInvalidPlayerID
	}

	_, err := s.playerRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, errs.ErrPlayerNotFound) {
			return errs.ErrPlayerNotFound
		}
		return err
	}
	if err := s.playerRepo.Delete(id); err != nil {
		return err
	}
	return nil
}

func (s *playerService) GetByID(id uint) (*models.Player, error) {
	if id == 0 {
		return nil, errs.ErrInvalidPlayerID
	}

	player, err := s.playerRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, errs.ErrPlayerNotFound) {
			return nil, errs.ErrPlayerNotFound
		}
		return nil, err
	}
	return player, nil
}

func (s *playerService) List(teamID uint) ([]models.Player, error) {
	if teamID == 0 {
		return nil, errs.ErrInvalidTeamID
	}

	players, err := s.playerRepo.List(teamID)
	if err != nil {
		return nil, err
	}
	return players, nil
}
