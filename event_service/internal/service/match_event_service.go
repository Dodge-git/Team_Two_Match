package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/mountainman199231/event_service/internal/client/match"
	"github.com/mountainman199231/event_service/internal/dto"
	"github.com/mountainman199231/event_service/internal/models"
	"github.com/mountainman199231/event_service/internal/repository"
)

type MatchEventService interface {
	CreateMatchEvent(ctx context.Context, req dto.CreateMatchEventRequest) (*dto.MatchEventResponse, error)

	GetMatchEventByID(ctx context.Context, id uint64) (*dto.MatchEventResponse, error)

	GetMatchEvents(ctx context.Context, matchID uint64, page int, pageSize int) (*dto.MatchEventListResponse, error)

	UpdateMatchEvent(ctx context.Context, id uint64, req dto.UpdateMatchEventRequest) (*dto.MatchEventResponse, error)

	GetMatchTimeline(ctx context.Context, matchID uint64, page int, pageSize int) (*dto.TimelineResponse, error)
}

type KafkaProduser interface {
	PublishMatchEventCreated(msg dto.MatchEventCreatedMessage) error
	PublishMatchGoal(msg dto.MatchGoalMessage) error
}

type matchEventService struct {
	matchEventRepo repository.MatchEventRepository
	commentaryRepo repository.CommentaryRepository
	reactionRepo   repository.ReactionRepository
	produser       KafkaProduser
	matchClient    match.Client
}

func NewMatchEventService(
	matchEventRepo repository.MatchEventRepository,
	commentaryRepo repository.CommentaryRepository,
	reactionRepo repository.ReactionRepository,
	produser KafkaProduser,
	matchClient match.Client,
) MatchEventService {
	return &matchEventService{
		matchEventRepo: matchEventRepo,
		commentaryRepo: commentaryRepo,
		reactionRepo:   reactionRepo,
		produser:       produser,
		matchClient:    matchClient,
	}
}

func (s *matchEventService) CreateMatchEvent(ctx context.Context, req dto.CreateMatchEventRequest) (*dto.MatchEventResponse, error) {

	match, err := s.matchClient.GetMatch(ctx, req.MatchID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch match: %w", err)
	}

	if match.Status != "live" {
		return nil, errors.New("cannot add event to non-live match")
	}

	if req.EventType == "goal" && req.TeamID == nil {
		return nil, errors.New("team_id is required for goal event")
	}

	if (req.EventType == "yellow_card" || req.EventType == "red_card") && req.PlayerID == nil {
		return nil, errors.New("player_id is required for card event")
	}

	event := &models.MatchEvent{
		MatchID:     req.MatchID,
		EventType:   models.EventType(req.EventType),
		Minute:      req.Minute,
		Period:      req.Period,
		TeamID:      req.TeamID,
		PlayerID:    req.PlayerID,
		Description: req.Description,
	}

	if err := s.matchEventRepo.Create(event); err != nil {
		return nil, fmt.Errorf("failed to save event: %w", err)
	}

	eventMsg := dto.MatchEventCreatedMessage{
		MatchID:     event.MatchID,
		EventType:   string(event.EventType),
		Minute:      event.Minute,
		Description: event.Description,
	}

	if err := s.produser.PublishMatchEventCreated(eventMsg); err != nil {
		log.Printf("failed to publish match.event.created: %v", err)
	}

	if event.EventType == models.EventGoal && event.TeamID != nil {
		goalMsg := dto.MatchGoalMessage{
			MatchID:  event.MatchID,
			TeamID:   *event.TeamID,
			PlayerID: event.PlayerID,
			Minute:   event.Minute,
		}

		if err := s.produser.PublishMatchGoal(goalMsg); err != nil {
			log.Printf("failed to publish match goal: %v", err)
		}
	}

	resp := dto.MatchEventResponse{
		ID:          event.ID,
		MatchID:     event.MatchID,
		EventType:   string(event.EventType),
		Minute:      event.Minute,
		Period:      event.Period,
		TeamID:      event.TeamID,
		PlayerID:    event.PlayerID,
		Description: event.Description,
		CreatedAt:   event.CreatedAt,
	}

	return &resp, nil

}

func (s *matchEventService) GetMatchEventByID(ctx context.Context, id uint64) (*dto.MatchEventResponse, error) {
	event, err := s.matchEventRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return mapToMatchEventResponse(event), nil
}

func (s *matchEventService) GetMatchEvents(ctx context.Context, matchID uint64, page int, pageSize int) (*dto.MatchEventListResponse, error) {

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 20
	}

	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	events, err := s.matchEventRepo.GetByMatchID(matchID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.matchEventRepo.CountByMatchID(matchID)
	if err != nil {
		return nil, err
	}

	items := make([]dto.MatchEventResponse, 0, len(events))
	for _, e := range events {
		items = append(items, *mapToMatchEventResponse(e))
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &dto.MatchEventListResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil

}

func (s *matchEventService) UpdateMatchEvent(ctx context.Context, id uint64, req dto.UpdateMatchEventRequest) (*dto.MatchEventResponse, error) {
	event, err := s.matchEventRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Minute != nil {
		event.Minute = *req.Minute
	}

	if req.Period != nil {
		event.Period = req.Period
	}

	if req.Description != nil {
		event.Description = *req.Description
	}

	if err := s.matchEventRepo.Update(event); err != nil {
		return nil, err
	}

	return mapToMatchEventResponse(event), nil
}

func (s *matchEventService) GetMatchTimeline(ctx context.Context, matchID uint64, page int, pageSize int) (*dto.TimelineResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	events, err := s.matchEventRepo.GetByMatchID(matchID, 1000, 0)
	if err != nil {
		return nil, err
	}

	commentaries, err := s.commentaryRepo.GetByMatchID(matchID, 1000, 0)
	if err != nil {
		return nil, err
	}

	var eventIDs []uint64
	for _, e := range events {
		eventIDs = append(eventIDs, e.ID)
	}

	var commentaryIDs []uint64
	for _, c := range commentaries {
		commentaryIDs = append(commentaryIDs, c.ID)
	}

	eventReactions := make(map[uint64]map[string]int)
	commentaryReactions := make(map[uint64]map[string]int)

	if len(eventIDs) > 0 {
		eventReactions, err = s.reactionRepo.GetEventReactionSummary(eventIDs)
		if err != nil {
			return nil, err
		}
	}

	if len(commentaryIDs) > 0 {
		commentaryReactions, err = s.reactionRepo.GetCommentaryReactionSummary(commentaryIDs)
		if err != nil {
			return nil, err
		}
	}

	var timeline []dto.TimelineItemDTO

	for _, e := range events {

		reactions := eventReactions[e.ID]
		if reactions == nil {
			reactions = map[string]int{}
		}

		timeline = append(timeline, dto.TimelineItemDTO{
			ID:        e.ID,
			Type:      dto.TimelineTypeEvent,
			EventType: string(e.EventType),
			Minute:    e.Minute,
			TeamID:    e.TeamID,
			PlayerID:  e.PlayerID,
			Reactions: reactions,
		})
	}

	for _, c := range commentaries {
		reactions := commentaryReactions[c.ID]
		if reactions == nil {
			reactions = map[string]int{}
		}

		timeline = append(timeline, dto.TimelineItemDTO{
			ID:        c.ID,
			Type:      dto.TimelineTypeCommentary,
			Minute:    c.Minute,
			Text:      c.Text,
			IsPinned:  c.IsPinned,
			Reactions: reactions,
		})
	}

	sort.Slice(timeline, func(i, j int) bool {
		if timeline[i].IsPinned != timeline[j].IsPinned {
			return timeline[i].IsPinned
		}

		if timeline[i].Minute == timeline[j].Minute {
			return timeline[i].ID < timeline[j].ID
		}
		return timeline[i].Minute < timeline[j].Minute
	})

	total := int64(len(timeline))

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(timeline) {
		start = len(timeline)
	}

	if end > len(timeline) {
		end = len(timeline)
	}

	pagedItems := timeline[start:end]

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &dto.TimelineResponse{
		Items:      pagedItems,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
