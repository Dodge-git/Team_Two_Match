package repository

import (
	"errors"
	"match_service/internal/errs"
	"match_service/internal/models"

	"gorm.io/gorm"
)

type SportRepository interface {
	Create(sport *models.Sport) error
	GetByID(id uint) (*models.Sport, error)
	Delete(id uint) error
	List() ([]models.Sport, error)
}

type sportRepository struct {
	db *gorm.DB
}

func NewSportRepository(db *gorm.DB) SportRepository {
	return &sportRepository{db: db}
}

func (r *sportRepository) Create(sport *models.Sport) error {
	return r.db.Create(sport).Error
}

func (r *sportRepository) GetByID(id uint) (*models.Sport, error) {
	var sport models.Sport
	if err := r.db.First(&sport, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrSportNotFound
		}
		return nil, err
	}
	return &sport, nil
}

func (r *sportRepository) Delete(id uint) error {
	return r.db.Delete(&models.Sport{}, id).Error
}

func (r *sportRepository) List() ([]models.Sport, error) {
	var sports []models.Sport

	if err := r.db.Find(&sports).Error; err != nil {
		return nil, err
	}
	return sports, nil
}
