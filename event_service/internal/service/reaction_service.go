package service

import (
	"context"
	"errors"

	"github.com/mountainman199231/event_service/internal/dto"
	"github.com/mountainman199231/event_service/internal/models"
	"github.com/mountainman199231/event_service/internal/repository"
	"gorm.io/gorm"
)

type ReactionService interface {
	SetReaction(ctx context.Context, userID uint64, req dto.SetReactionRequest) (*dto.ReactionResponse, error)

	GetEventSummary(ctx context.Context, eventID uint64) (*dto.ReactionSummaryResponse, error)
	GetCommentarySummary(ctx context.Context, commentaryID uint64) (*dto.ReactionSummaryResponse, error)

	GetUserEventReaction(ctx context.Context, userID uint64, eventID uint64) (*dto.UserReactionResponse, error)
	GetUserCommentaryReaction(ctx context.Context, userID uint64, commentaryID uint64) (*dto.UserReactionResponse, error)
}

type reactionService struct {
	reactionRepo   repository.ReactionRepository
	eventRepo      repository.MatchEventRepository
	commentaryRepo repository.CommentaryRepository
}

func NewReactionService(
	reactionRepo repository.ReactionRepository,
	eventRepo repository.MatchEventRepository,
	commentaryRepo repository.CommentaryRepository,
) ReactionService {
	return &reactionService{
		reactionRepo:   reactionRepo,
		eventRepo:      eventRepo,
		commentaryRepo: commentaryRepo,
	}
}

func (s *reactionService) SetReaction(ctx context.Context, userID uint64, req dto.SetReactionRequest) (*dto.ReactionResponse, error) {
	if (req.CommentaryID == nil && req.EventID == nil) ||
		(req.CommentaryID != nil && req.EventID != nil) {
		return nil, errors.New("reaction must be attached either to event or commentary")
	}

	if req.EventID != nil {
		_, err := s.eventRepo.GetByID(*req.EventID)
		if err != nil {
			return nil, errors.New("event not found")
		}
	}

	if req.CommentaryID != nil {
		_, err := s.commentaryRepo.GetByID(*req.CommentaryID)
		if err != nil {
			return nil, errors.New("commentary not found")
		}
	}

	reactionType := models.ReactionType(req.Type)

	var existing *models.Reaction
	var err error

	if req.EventID != nil {
		existing, err = s.reactionRepo.GetByUserAndEvent(userID, *req.EventID)
	} else {
		existing, err = s.reactionRepo.GetByUserAndCommentary(userID, *req.CommentaryID)
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			newReaction := &models.Reaction{
				UserID:       userID,
				EventID:      req.EventID,
				CommentaryID: req.CommentaryID,
				Type:         reactionType,
			}

			if err := s.reactionRepo.Create(newReaction); err != nil {
				return nil, err
			}
			return mapReactionToDTO(newReaction), nil
		}
		return nil, err
	}

	if existing.Type == reactionType {
		if err := s.reactionRepo.Delete(existing.ID); err != nil {
			return nil, err
		}
		return nil, nil
	}

	existing.Type = reactionType
	if err := s.reactionRepo.Update(existing); err != nil {
		return nil, err
	}

	return mapReactionToDTO(existing), nil
}

func (s *reactionService) GetEventSummary(ctx context.Context, eventID uint64) (*dto.ReactionSummaryResponse, error) {
	data, err := s.reactionRepo.GetGroupedByEvent(eventID)
	if err != nil {
		return nil, err
	}

	return buildReactionSummary(data), nil
}

func (s *reactionService) GetCommentarySummary(ctx context.Context, commentaryID uint64) (*dto.ReactionSummaryResponse, error) {
	data, err := s.reactionRepo.GetGroupedByCommentary(commentaryID)
	if err != nil {
		return nil, err
	}

	return buildReactionSummary(data), nil
}

func (s *reactionService) GetUserEventReaction(ctx context.Context, userID uint64, eventID uint64) (*dto.UserReactionResponse, error) {
	reaction, err := s.reactionRepo.GetByUserAndEvent(userID, eventID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &dto.UserReactionResponse{Type: nil}, nil
	}

	if err != nil {
		return nil, err
	}

	t := string(reaction.Type)

	return &dto.UserReactionResponse{Type: &t}, nil
}

func (s *reactionService) GetUserCommentaryReaction(ctx context.Context, userID uint64, commentaryID uint64) (*dto.UserReactionResponse, error) {
	reaction, err := s.reactionRepo.GetByUserAndCommentary(userID, commentaryID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &dto.UserReactionResponse{Type: nil}, nil
	}

	if err != nil {
		return nil, err
	}

	t := string(reaction.Type)
	return &dto.UserReactionResponse{Type: &t}, nil
}
