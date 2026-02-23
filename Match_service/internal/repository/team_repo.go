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
	List(filter models.TeamFilter) ([]models.Team, int64, error)
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

func (r *teamRepository) List(filter models.TeamFilter) ([]models.Team, int64, error) {
	var teams []models.Team
	var total int64

	query := r.db.Model(&models.Team{})

	if filter.SportID != nil {
		query = query.Where("sport_id = ?", *filter.SportID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize

	if err := query.Order("name ASC").Limit(filter.PageSize).Offset(offset).Find(&teams).Error; err != nil {
		return nil, 0, err
	}
	return teams, total, nil
}
