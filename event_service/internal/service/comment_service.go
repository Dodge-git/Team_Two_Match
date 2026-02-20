package service

import (
	"context"
	"errors"

	"github.com/mountainman199231/event_service/internal/dto"
	"github.com/mountainman199231/event_service/internal/models"
	"github.com/mountainman199231/event_service/internal/repository"
)

type CommentService interface {
	Create(ctx context.Context, userID uint64, req dto.CreateCommentRequest) (*dto.CommentResponse, error)
	GetByEvent(ctx context.Context, eventID uint64, page, pageSize int) (*dto.CommentListResponse, error)
	GetByCommentary(ctx context.Context, commentaryID uint64, page, pageSize int) (*dto.CommentListResponse, error)

	UpdateComment(ctx context.Context, userID, id uint64, req dto.UpdateCommentRequest) error
	DeleteComment(ctx context.Context, userID, id uint64) error
}

type commentService struct {
	commentRepo repository.CommentRepository
}

func NewCommentService(commentRepo repository.CommentRepository) CommentService {
	return &commentService{commentRepo: commentRepo}
}

func (s *commentService) Create(ctx context.Context, userID uint64, req dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	if (req.EventID == nil && req.CommentaryID == nil) ||
		(req.EventID != nil && req.CommentaryID != nil) {
		return nil, errors.New("comment must be attached to either event or commentary")
	}

	comment := models.Comment{
		UserID:       userID,
		EventID:      req.EventID,
		CommentaryID: req.CommentaryID,
		Text:         req.Text,
	}

	if err := s.commentRepo.Create(&comment); err != nil {
		return nil, err
	}

	return mapCommentToDTO(&comment), nil
}

func (s *commentService) GetByEvent(ctx context.Context, eventID uint64, page, pageSize int) (*dto.CommentListResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize

	comments, total, err := s.commentRepo.GetByEventID(eventID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	return buildListResponse(comments, total, page, pageSize), nil
}

func (s *commentService) GetByCommentary(ctx context.Context, commentaryID uint64, page, pageSize int) (*dto.CommentListResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize

	comments, total, err := s.commentRepo.GetByCommentaryID(commentaryID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	return buildListResponse(comments, total, page, pageSize), nil
}

func (s *commentService) UpdateComment(ctx context.Context, userID uint64, id uint64, req dto.UpdateCommentRequest) error {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return err
	}

	if comment.UserID != userID {
		return errors.New("forbidden")
	}
	if req.Text != nil {
		comment.Text = *req.Text
	}

	return s.commentRepo.Update(comment)
}

func (s *commentService) DeleteComment(ctx context.Context, userID, id uint64) error {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return err
	}

	if comment.UserID != userID {
		return errors.New("forbidden")
	}

	return s.commentRepo.Delete(id)
}

func normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	return page, pageSize
}
