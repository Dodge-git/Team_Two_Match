package repository

import (
	"errors"
	"match_service/internal/errs"
	"match_service/internal/models"

	"gorm.io/gorm"
)

type TeamRepository interface {
	Create(team *models.Team) error
	GetByID(id uint) (*models.Team, error)
	Delete(id uint) error
	List(sportID uint) ([]models.Team, error)
}

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) Create(team *models.Team) error {
	if team == nil {
		return errors.New("team is nil")
	}

	if err := r.db.Create(team).Error; err != nil {
		return err
	}
	return nil
}

func (r *teamRepository) GetByID(id uint) (*models.Team, error) {

	var team models.Team

	if err := r.db.First(&team, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrTeamNotFound
		}
		return nil, err
	}
	return &team, nil
}

func (r *teamRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.Team{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *teamRepository) List(sportID uint) ([]models.Team, error) {
	var teams []models.Team

	if err := r.db.Where("sport_id = ?", sportID).Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}
