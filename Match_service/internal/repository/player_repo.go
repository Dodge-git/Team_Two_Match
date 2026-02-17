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
}

type gormPlayerRepository struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return &gormPlayerRepository{db: db}
}

func (r *gormPlayerRepository) Create(player *models.Player) error {
	return r.db.Create(player).Error
}

func (r *gormPlayerRepository) GetByID(playerID uint) (*models.Player, error) {
	var player models.Player
	if err := r.db.First(&player, playerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrPlayerNotFound
		}
		return nil, err
	}
	return &player, nil
}

func (r *gormPlayerRepository) Delete(playerID uint) error {
	return r.db.Delete(&models.Player{}, playerID).Error
}

func (r *gormPlayerRepository) Update(player *models.Player) error {
	return r.db.Save(player).Error
}
