package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/mountainman199231/event_service/internal/dto"
	"github.com/mountainman199231/event_service/internal/models"
	"github.com/mountainman199231/event_service/internal/repository"

	"github.com/redis/go-redis/v9"
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
	rdb         *redis.Client // redis
}

func NewCommentService(commentRepo repository.CommentRepository, rdb *redis.Client) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		rdb:         rdb, // redis
	}
}

// redis helper
func (s *commentService) invalidateCommentaryCache(ctx context.Context, commentaryID uint64) {
	setKey := fmt.Sprintf("comments:index:%d", commentaryID)

	keys, err := s.rdb.SMembers(ctx, setKey).Result()
	if err == nil && len(keys) > 0 {
		s.rdb.Del(ctx, keys...)
	}

	s.rdb.Del(ctx, setKey)
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
	//redis
	if comment.CommentaryID != nil {
		s.invalidateCommentaryCache(ctx, *comment.CommentaryID)
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

// формируем ключ кеша
func buildCacheKey(commentaryID uint64, page, pageSize int) string {
	return fmt.Sprintf("comments:%d:%d:%d", commentaryID, page, pageSize)
}

func (s *commentService) GetByCommentary(ctx context.Context, commentaryID uint64, page, pageSize int) (*dto.CommentListResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)

	cacheKey := buildCacheKey(commentaryID, page, pageSize)
	setKey := fmt.Sprintf("comments:index:%d", commentaryID)

	cached, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var resp dto.CommentListResponse
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			return &resp, nil
		}
	} else if err != redis.Nil {
		fmt.Println("redis error:", err)
	}

	offset := (page - 1) * pageSize

	comments, total, err := s.commentRepo.GetByCommentaryID(commentaryID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	resp := buildListResponse(comments, total, page, pageSize)

	data, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	if err := s.rdb.Set(ctx, cacheKey, data, 5*time.Minute).Err(); err == nil {
		s.rdb.SAdd(ctx, setKey, cacheKey)
		s.rdb.Expire(ctx, setKey, 5*time.Minute)
	}

	return resp, nil
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

	if err := s.commentRepo.Update(comment); err != nil {
		return err
	}

	if comment.CommentaryID != nil {
		s.invalidateCommentaryCache(ctx, *comment.CommentaryID)
	}

	return nil
}

func (s *commentService) DeleteComment(ctx context.Context, userID, id uint64) error {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return err
	}

	if comment.UserID != userID {
		return errors.New("forbidden")
	}

	if err := s.commentRepo.Delete(id); err != nil {
		return err
	}

	if comment.CommentaryID != nil {
		s.invalidateCommentaryCache(ctx, *comment.CommentaryID)
	}
	return nil
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
