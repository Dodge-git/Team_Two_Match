package repository

import (
	"context"
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
	List(models.MatchFilter) ([]models.Match, int64, error)
	GetActive() ([]models.Match, error)
	UpdateScoreFromKafka(ctx context.Context, matchID uint, teamID uint, newScore int) error
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

	if err := r.db.
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("Sport").
		First(&match, id).Error; err != nil {
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

func (r *matchRepository) List(filter models.MatchFilter) ([]models.Match, int64, error) {
	var matches []models.Match
	var total int64
	query := r.db.Model(&models.Match{})

	if filter.SportID != nil {
		query = query.Where("sport_id = ?", *filter.SportID)
	}

	if filter.DateFrom != nil {
		query = query.Where("scheduled_at >= ?", *filter.DateFrom)
	}

	if filter.DateTo != nil {
		query = query.Where("scheduled_at <= ?", *filter.DateTo)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)

	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize

	if err := query.
		Limit(filter.PageSize).
		Offset(offset).
		Order("scheduled_at ASC").
		Find(&matches).Error; err != nil {

		return nil, 0, err
	}
	return matches, total, nil
}

func (r *matchRepository) GetActive() ([]models.Match, error) {
	var matches []models.Match

	if err := r.db.Where("status = ?", "live").Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

func (r *matchRepository) UpdateScoreFromKafka(ctx context.Context, matchID uint, teamID uint, newScore int) error {
	result := r.db.WithContext(ctx).Model(&models.Match{}).
		Where("id = ? AND  status = ? AND home_team_id = ? AND home_score < ?", matchID, models.MatchStatusLive, teamID, newScore).
		UpdateColumn("home_score", newScore)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		return nil
	}

	result = r.db.WithContext(ctx).Model(&models.Match{}).
		Where("id = ? AND  status = ? AND away_team_id = ? AND away_score < ?", matchID, models.MatchStatusLive, teamID, newScore).
		UpdateColumn("away_score", newScore)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		return nil
	}
	return errs.ErrInvalidGoalEvent
}
