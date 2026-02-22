package repository

import (
	"github.com/mountainman199231/event_service/internal/models"
	"gorm.io/gorm"
)

type ReactionRepository interface {
	Create(reaction *models.Reaction) error
	Update(reaction *models.Reaction) error
	Delete(id uint64) error

	GetByID(id uint64) (*models.Reaction, error)

	GetByUserAndEvent(userID, eventID uint64) (*models.Reaction, error)
	GetByUserAndCommentary(userID, commentaryID uint64) (*models.Reaction, error)

	GetGroupedByEvent(eventID uint64) (map[string]int64, error)
	GetGroupedByCommentary(commentaryID uint64) (map[string]int64, error)

	GetEventReactionSummary(eventIDs []uint64) (map[uint64]map[string]int, error)
	GetCommentaryReactionSummary(commentaryIDs []uint64) (map[uint64]map[string]int, error)
}

type gormReactionRepository struct {
	db *gorm.DB
}

func NewReactionRepository(db *gorm.DB) ReactionRepository {
	return &gormReactionRepository{db: db}
}

func (r *gormReactionRepository) Create(reaction *models.Reaction) error {
	return r.db.Create(reaction).Error
}

func (r *gormReactionRepository) Update(reaction *models.Reaction) error {
	return r.db.Save(reaction).Error
}

func (r *gormReactionRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Reaction{}, id).Error
}

func (r *gormReactionRepository) GetByID(id uint64) (*models.Reaction, error) {
	var reaction models.Reaction

	if err := r.db.First(&reaction, id).Error; err != nil {
		return nil, err
	}

	return &reaction, nil

}

func (r *gormReactionRepository) GetByUserAndEvent(userID, eventID uint64) (*models.Reaction, error) {
	var reaction models.Reaction

	err := r.db.
		Where("user_id = ? AND event_id = ?", userID, eventID).
		First(&reaction).Error

	if err != nil {
		return nil, err
	}

	return &reaction, nil
}

func (r *gormReactionRepository) GetByUserAndCommentary(userID, commentaryID uint64) (*models.Reaction, error) {
	var reaction models.Reaction

	err := r.db.
		Where("user_id = ? AND commentary_id = ?", userID, commentaryID).
		First(&reaction).Error

	if err != nil {
		return nil, err
	}

	return &reaction, nil
}

func (r *gormReactionRepository) GetGroupedByEvent(eventID uint64) (map[string]int64, error) {
	type result struct {
		Type  string
		Count int64
	}

	var rows []result

	err := r.db.Model(&models.Reaction{}).
		Select("type, COUNT(*) as count").
		Where("event_id = ?", eventID).
		Group("type").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	response := make(map[string]int64)
	for _, row := range rows {
		response[row.Type] = row.Count
	}

	return response, nil
}

func (r *gormReactionRepository) GetGroupedByCommentary(commentaryID uint64) (map[string]int64, error) {
	type result struct {
		Type  string
		Count int64
	}

	var rows []result

	err := r.db.Model(&models.Reaction{}).
		Select("type, COUNT(*) as count").
		Where("commentary_id = ?", commentaryID).
		Group("type").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	response := make(map[string]int64)
	for _, row := range rows {
		response[row.Type] = row.Count
	}

	return response, nil

}

func (r *gormReactionRepository) GetEventReactionSummary(eventIDs []uint64) (map[uint64]map[string]int, error) {

	type result struct {
		EventID uint64
		Type    string
		Count   int
	}

	var rows []result

	err := r.db.Model(&models.Reaction{}).
		Select("event_id, type, COUNT(*) as count").
		Where("event_id IN ?", eventIDs).
		Group("event_id, type").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	response := make(map[uint64]map[string]int)

	for _, row := range rows {
		if response[row.EventID] == nil {
			response[row.EventID] = make(map[string]int)
		}
		response[row.EventID][row.Type] = row.Count
	}

	return response, nil
}

func (r *gormReactionRepository) GetCommentaryReactionSummary(commentaryIDs []uint64) (map[uint64]map[string]int, error) {

	type result struct {
		CommentaryID uint64
		Type         string
		Count        int
	}

	var rows []result

	err := r.db.Model(&models.Reaction{}).
		Select("commentary_id, type, COUNT(*) as count").
		Where("commentary_id IN ?", commentaryIDs).
		Group("commentary_id, type").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	response := make(map[uint64]map[string]int)

	for _, row := range rows {
		if response[row.CommentaryID] == nil {
			response[row.CommentaryID] = make(map[string]int)
		}
		response[row.CommentaryID][row.Type] = row.Count
	}

	return response, nil
}
