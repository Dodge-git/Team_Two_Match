package repository

import (
	"github.com/mountainman199231/event_service/internal/models"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *models.Comment) error
	GetByID(id uint64) (*models.Comment, error)

	GetByEventID(eventID uint64, limit, offset int) ([]*models.Comment, int64, error)
	GetByCommentaryID(commentaryID uint64, limit, offset int) ([]*models.Comment, int64, error)

	Update(comment *models.Comment) error
	Delete(id uint64) error
}

type gormCommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &gormCommentRepository{db: db}
}

func (r *gormCommentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *gormCommentRepository) GetByID(id uint64) (*models.Comment, error) {
	var comment models.Comment

	if err := r.db.First(&comment, id).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *gormCommentRepository) GetByEventID(eventID uint64, limit, offset int) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var total int64

	query := r.db.Model(&models.Comment{}).
		Where("event_id = ?", eventID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

func (r *gormCommentRepository) GetByCommentaryID(commentaryID uint64, limit, offset int) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var total int64

	query := r.db.Model(&models.Comment{}).
		Where("commentary_id = ?", commentaryID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

func (r *gormCommentRepository) Update(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

func (r *gormCommentRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Comment{}, id).Error
}
