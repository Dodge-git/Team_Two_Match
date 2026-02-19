package repository

import (
	"github.com/mountainman199231/event_service/internal/models"
	"gorm.io/gorm"
)

type CommentaryRepository interface {
	WithDB(db *gorm.DB) CommentaryRepository
	Create(commentary *models.Commentary) error
	GetByID(id uint64) (*models.Commentary, error)
	GetByMatchID(matchID uint64, limit, offset int) ([]*models.Commentary, error)
	Update(commentary *models.Commentary) error
	DeleteByID(id uint64) error

	UnpinAllByMatchID(matchID uint64) error
	SetPinned(id uint64) error
}

type gormCommentaryRepository struct {
	db *gorm.DB
}

func NewCommentaryRepository(db *gorm.DB) CommentaryRepository {
	return &gormCommentaryRepository{db: db}
}

func (r *gormCommentaryRepository) WithDB(db *gorm.DB) CommentaryRepository {
	return &gormCommentaryRepository{db: db}
}

func (r *gormCommentaryRepository) Create(commentary *models.Commentary) error {
	return r.db.Create(commentary).Error
}

func (r *gormCommentaryRepository) GetByID(id uint64) (*models.Commentary, error) {
	var commentary models.Commentary

	if err := r.db.First(&commentary, id).Error; err != nil {
		return nil, err
	}

	return &commentary, nil
}

func (r *gormCommentaryRepository) GetByMatchID(matchID uint64, limit, offset int) ([]*models.Commentary, error) {
	var commentaries []*models.Commentary

	err := r.db.
		Where("match_id = ?", matchID).
		Order("minute ASC").
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&commentaries).Error

	return commentaries, err

}

func (r *gormCommentaryRepository) Update(commentary *models.Commentary) error {
	return r.db.Save(commentary).Error
}

func (r *gormCommentaryRepository) DeleteByID(id uint64) error {
	return r.db.Delete(&models.Commentary{}, id).Error
}

func (r *gormCommentaryRepository) UnpinAllByMatchID(matchID uint64) error {
	return r.db.
		Model(&models.Commentary{}).
		Where("match_id = ?", matchID).
		Update("is_pinned", false).Error

}

func (r *gormCommentaryRepository) SetPinned(id uint64) error {
	result := r.db.
		Model(&models.Commentary{}).
		Where("id = ?", id).
		Update("is_pinned", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil

}
