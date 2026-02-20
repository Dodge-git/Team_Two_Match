package repository

import (
	"github.com/mountainman199231/event_service/internal/models"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *models.Comment) error
	GetByID(id uint64) (*models.Comment, error)

	GetByEventID(eventID uint64, limit, offset int) ([]*models.Comment, error)
	GetByCommentaryID(commentaryID uint64, limit, offset int) ([]*models.Comment, error)

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

func (r *gormCommentRepository) GetByEventID(eventID uint64, limit, offset int) ([]*models.Comment, error) {

	var comments []*models.Comment

	err := r.db.
		Where("event_id = ?", eventID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&comments).Error

	return comments, err
}

func (r *gormCommentRepository) GetByCommentaryID(commentaryID uint64, limit, offset int) ([]*models.Comment, error) {
	var comments []*models.Comment

	err := r.db.
		Where("commentary_id = ?", commentaryID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&comments).Error

	return comments, err
}

func (r *gormCommentRepository) Update(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

func (r *gormCommentRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Comment{}, id).Error
}
