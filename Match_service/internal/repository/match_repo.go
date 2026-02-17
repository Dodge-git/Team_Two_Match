package repository

import (
	"errors"
	"match_service/internal/errs"
	"match_service/internal/models"

	"gorm.io/gorm"
)

type MatchRepository interface {
	Create(match *models.Match) error
	GetByID(id uint) (*models.Match, error)
	Update(match *models.Match) error
	Delete(id uint) error
}

type matchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepository {
	return &matchRepository{db: db}
}

func (r *matchRepository) Create(match *models.Match) error {
	if match == nil {
		return errors.New("match is nil")
	}

	if err := r.db.Create(match).Error; err != nil {
		return err
	}
	return nil
}

func (r *matchRepository) GetByID(id uint) (*models.Match, error) {
	var match models.Match

	if err := r.db.First(&match, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrMatchNotFound
		}
		return nil, err
	}
	return &match, nil
}

func (r *matchRepository) Update(match *models.Match) error {
	if match == nil {
		return errors.New("match is nil")
	}
	return r.db.Save(match).Error
}

func (r *matchRepository) Delete(id uint) error {
	return r.db.Delete(&models.Match{}, id).Error
}
