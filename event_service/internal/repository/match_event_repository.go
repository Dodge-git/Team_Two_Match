package repository

import (
	"github.com/mountainman199231/event_service/internal/models"
	"gorm.io/gorm"
)

type MatchEventRepository interface {
	Create(event *models.MatchEvent) error

	GetByID(id uint64) (*models.MatchEvent, error)

	GetByMatchID(matchID uint64, limit, offset int) ([]*models.MatchEvent, error)
	CountByMatchID(matchID uint64) (int64, error)

	Update(event *models.MatchEvent) error
}

type gormMatchEventRepository struct {
	db *gorm.DB
}

func NewMatchEventRepository(db *gorm.DB) MatchEventRepository {
	return &gormMatchEventRepository{db: db}
}

func (r *gormMatchEventRepository) Create(event *models.MatchEvent) error {
	return r.db.Create(event).Error
}

func (r *gormMatchEventRepository) GetByID(id uint64) (*models.MatchEvent, error) {
	var event models.MatchEvent

	if err := r.db.First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *gormMatchEventRepository) GetByMatchID(matchID uint64, limit, offset int) ([]*models.MatchEvent, error) {
	var events []*models.MatchEvent

	err := r.db.
		Where("match_id = ?", matchID).
		Order("minute ASC, id ASC").
		Limit(limit).
		Offset(offset).
		Find(&events).Error

	return events, err
}

func (r *gormMatchEventRepository) CountByMatchID(matchID uint64) (int64, error) {
	var count int64

	err := r.db.
		Model(&models.MatchEvent{}).
		Where("match_id = ?", matchID).
		Count(&count).Error

	return count, err
}

func (r *gormMatchEventRepository) Update(event *models.MatchEvent) error {
	return r.db.Save(event).Error
}
