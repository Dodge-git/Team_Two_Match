package service

import (
	"context"
	"errors"

	"github.com/mountainman199231/event_service/internal/client/match"
	"github.com/mountainman199231/event_service/internal/dto"
	"github.com/mountainman199231/event_service/internal/models"
	"github.com/mountainman199231/event_service/internal/repository"
	"gorm.io/gorm"
)

type CommentaryService interface {
	CreateCommentary(ctx context.Context, req dto.CreateCommentaryRequest) (*dto.CommentaryResponse, error)
	GetCommentaryByMatchID(ctx context.Context, matchID uint64, limit, offset int) ([]*dto.CommentaryResponse, error)
	DeleteCommentary(ctx context.Context, id uint64) error
	PinCommentary(ctx context.Context, id uint64) error
}

type commentaryService struct {
	commentaryRepo repository.CommentaryRepository
	db             *gorm.DB
	matchClient    match.Client
}

func NewCommentaryService(
	commentaryRepo repository.CommentaryRepository,
	db *gorm.DB,
	matchClient match.Client,
) CommentaryService {
	return &commentaryService{
		commentaryRepo: commentaryRepo,
		db:             db,
		matchClient:    matchClient,
	}
}

func (s *commentaryService) CreateCommentary(ctx context.Context, req dto.CreateCommentaryRequest) (*dto.CommentaryResponse, error) {
	match, err := s.matchClient.GetMatch(ctx, req.MatchID)
	if err != nil {
		return nil, err
	}
	if match.Status != "live" {
		return nil, errors.New("match is not live")
	}

	commentary := &models.Commentary{
		MatchID:  req.MatchID,
		Minute:   req.Minute,
		Text:     req.Text,
		IsPinned: false,
	}

	if err := s.commentaryRepo.Create(commentary); err != nil {
		return nil, err
	}

	return &dto.CommentaryResponse{
		ID:        commentary.ID,
		MatchID:   commentary.MatchID,
		Minute:    commentary.Minute,
		Text:      commentary.Text,
		IsPinned:  commentary.IsPinned,
		CreatedAt: commentary.CreatedAt,
		UpdatedAt: commentary.UpdatedAt,
	}, nil
}

func (s *commentaryService) GetCommentaryByMatchID(ctx context.Context, matchID uint64, limit, offset int) ([]*dto.CommentaryResponse, error) {

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	commentaries, err := s.commentaryRepo.GetByMatchID(matchID, limit, offset)
	if err != nil {
		return nil, err
	}

	var result []*dto.CommentaryResponse
	for _, c := range commentaries {
		result = append(result, &dto.CommentaryResponse{
			ID:        c.ID,
			MatchID:   c.MatchID,
			Minute:    c.Minute,
			Text:      c.Text,
			IsPinned:  c.IsPinned,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})

	}

	return result, nil
}

func (s *commentaryService) DeleteCommentary(ctx context.Context, id uint64) error {

	return s.commentaryRepo.DeleteByID(id)
}

func (s *commentaryService) PinCommentary(ctx context.Context, id uint64) error {
	commentary, err := s.commentaryRepo.GetByID(id)
	if err != nil {
		return err
	}

	match, err := s.matchClient.GetMatch(ctx, commentary.MatchID)
	if err != nil {
		return err
	}

	if match.Status != "live" {
		return errors.New("can not pin comment: match is not live")
	}

	if commentary.IsPinned {
		return nil
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	txRepo := s.commentaryRepo.WithDB(tx)

	if err := txRepo.UnpinAllByMatchID(commentary.MatchID); err != nil {
		tx.Rollback()
		return err
	}

	if err := txRepo.SetPinned(id); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
