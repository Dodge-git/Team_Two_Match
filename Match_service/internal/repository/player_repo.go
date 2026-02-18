package repository

import (
	"errors"
	"match_service/internal/errs"
	"match_service/internal/models"

	"gorm.io/gorm"
)

type PlayerRepository interface {
	Create(player *models.Player) error
	GetByID(playerID uint) (*models.Player, error)
	Delete(playerID uint) error
	Update(player *models.Player) error
	List(teamID uint) ([]models.Player, error)
}

type playerRepository struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return &playerRepository{db: db}
}

func (r *playerRepository) Create(player *models.Player) error {
	return r.db.Create(player).Error
}

func (r *playerRepository) GetByID(playerID uint) (*models.Player, error) {
	var player models.Player
	if err := r.db.First(&player, playerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrPlayerNotFound
		}
		return nil, err
	}
	return &player, nil
}

func (r *playerRepository) Delete(playerID uint) error {
	return r.db.Delete(&models.Player{}, playerID).Error
}

func (r *playerRepository) Update(player *models.Player) error {
	return r.db.Save(player).Error
}

func (r *playerRepository) List(teamID uint) ([]models.Player, error) {
	var players []models.Player

	if err := r.db.Where("teamID = ?", teamID).Find(&players).Error; err != nil {
		return nil, err
	}
	return players, nil
}
