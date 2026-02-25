package services

import (
	"context"
	"errors"
	"match_service/internal/dto"
	"match_service/internal/errs"
	"match_service/internal/models"
	"match_service/internal/ports"
	"match_service/internal/repository"
	"time"
)

type MatchService interface {
	CreateMatch(req dto.CreateMatchRequest) (*models.Match, error)
	GetMatchByID(id uint) (*models.Match, error)
	DeleteMatch(id uint) error
	ListMatches(filter models.MatchFilter) (*dto.MatchListResponse, error)
	StartMatch(ctx context.Context, id uint) error
	FinishMatch(ctx context.Context, id uint) error
	CancelMatch(id uint) error
	UpdateScoreFromKafka(ctx context.Context, goalEvent models.GoalEvent) error
	GetActive() ([]models.Match, error)
}

type matchService struct {
	matchRepo repository.MatchRepository
	sportRepo repository.SportRepository
	teamRepo  repository.TeamRepository
	producer  ports.Producer
	//kafkaConsumer kafka.Consumer
}

func NewMatchService(matchRepo repository.MatchRepository, sportRepo repository.SportRepository, teamRepo repository.TeamRepository, producer ports.Producer) MatchService {
	return &matchService{matchRepo: matchRepo, sportRepo: sportRepo, teamRepo: teamRepo, producer: producer}
}

func (s *matchService) CreateMatch(req dto.CreateMatchRequest) (*models.Match, error) {
	if req.HomeTeamID == req.AwayTeamID {
		return nil, errors.New("home team and away team cannot be the same")
	}

	_, err := s.sportRepo.GetByID(req.SportID)
	if err != nil {
		if errors.Is(err, errs.ErrSportNotFound) {
			return nil, errs.ErrSportNotFound
		}
		return nil, err
	}
	homeTeam, err := s.teamRepo.GetByID(req.HomeTeamID)
	if err != nil {
		if errors.Is(err, errs.ErrTeamNotFound) {
			return nil, errs.ErrTeamNotFound
		}
		return nil, err
	}

	awayTeam, err := s.teamRepo.GetByID(req.AwayTeamID)
	if err != nil {
		if errors.Is(err, errs.ErrTeamNotFound) {
			return nil, errs.ErrTeamNotFound
		}
		return nil, err
	}

	if homeTeam.SportID != req.SportID || awayTeam.SportID != req.SportID {
		return nil, errs.ErrTeamsNotInSport
	}

	if req.ScheduledAt == nil || req.ScheduledAt.Before(time.Now().UTC()) {
		return nil, errors.New("scheduled_at must be a future date and time")
	}

	match := &models.Match{
		SportID:        req.SportID,
		HomeTeamID:     homeTeam.ID,
		AwayTeamID:     awayTeam.ID,
		Status:         models.MatchStatusScheduled,
		ScheduledAt:    *req.ScheduledAt,
		Venue:          req.Venue,
		TournamentName: req.TournamentName,
	}

	if err := s.matchRepo.Create(match); err != nil {
		return nil, err
	}
	return match, nil
}

func (r *matchService) GetMatchByID(id uint) (*models.Match, error) {
	if id == 0 {
		return nil, errs.ErrInvalidMatchID
	}
	match, err := r.matchRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, errs.ErrMatchNotFound) {
			return nil, errs.ErrMatchNotFound
		}
		return nil, err
	}
	return match, nil
}

func (s *matchService) DeleteMatch(id uint) error {
	if id == 0 {
		return errs.ErrInvalidMatchID
	}

	match, err := s.matchRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, errs.ErrMatchNotFound) {
			return errs.ErrMatchNotFound
		}
		return err
	}
	if match.Status != models.MatchStatusScheduled {
		return errors.New("only scheduled matches can be deleted")
	}
	if err := s.matchRepo.Delete(id); err != nil {
		return err
	}
	return nil
}

func (s *matchService) ListMatches(filter models.MatchFilter) (*dto.MatchListResponse, error) {

	if filter.Page < 1 {
		filter.Page = 1
	}

	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	if filter.SportID != nil {
		if _, err := s.sportRepo.GetByID(*filter.SportID); err != nil {
			return nil, err
		}
	}
	if filter.Status != nil {
		status := *filter.Status
		if status != models.MatchStatusScheduled &&
			status != models.MatchStatusLive &&
			status != models.MatchStatusFinished &&
			status != models.MatchStatusCancelled {
			return nil, errs.ErrInvalidMatchStatus
		}
	}

	if filter.DateFrom != nil && filter.DateTo != nil && filter.DateFrom.After(*filter.DateTo) {
		return nil, errors.New("date_from cannot be after date_to")
	}

	matches, total, err := s.matchRepo.List(filter)
	if err != nil {
		return nil, err
	}

	data := make([]dto.MatchResponse, 0, len(matches))

	for _, match := range matches {
		data = append(data, dto.MatchResponse{
			ID:             match.ID,
			SportID:        match.SportID,
			HomeTeamID:     match.HomeTeamID,
			AwayTeamID:     match.AwayTeamID,
			Status:         match.Status,
			ScheduledAt:    match.ScheduledAt,
			FinishedAt:     match.FinishedAt,
			StartedAt:      match.StartedAt,
			HomeScore:      match.HomeScore,
			AwayScore:      match.AwayScore,
			Venue:          match.Venue,
			TournamentName: match.TournamentName,
		})
	}
	resp := &dto.MatchListResponse{
		Data:     data,
		Page:     filter.Page,
		PageSize: filter.PageSize,
		Total:    total,
	}
	return resp, nil
}

func (s *matchService) StartMatch(ctx context.Context, id uint) error {
	if id == 0 {
		return errs.ErrInvalidMatchID
	}
	match, err := s.matchRepo.GetByID(id)
	if err != nil {
		return err
	}
	if match.Status != models.MatchStatusScheduled {
		return errors.New("match cannot be started")
	}

	now := time.Now().UTC()
	match.Status = models.MatchStatusLive
	match.StartedAt = &now

	if err := s.matchRepo.Update(match); err != nil {
		return err
	}

	event := dto.MatchStartedEvent{
		MatchID:  match.ID,
		HomeTeam: match.HomeTeam.Name,
		AwayTeam: match.AwayTeam.Name,
		Sport:    string(match.Sport.Name),
	}

	if err := s.producer.Publish(ctx, "match.started", event); err != nil {
		return err
	}

	return nil
}

func (s *matchService) FinishMatch(ctx context.Context, id uint) error {
	if id == 0 {
		return errs.ErrInvalidMatchID
	}
	match, err := s.matchRepo.GetByID(id)
	if err != nil {
		return err
	}
	if match.Status != models.MatchStatusLive {
		return errors.New("match cannot be finished")
	}

	now := time.Now().UTC()
	match.Status = models.MatchStatusFinished
	match.FinishedAt = &now

	if err := s.matchRepo.Update(match); err != nil {
		return err
	}
	event := dto.MatchEndedEvent{
		MatchID:   match.ID,
		HomeScore: match.HomeScore,
		AwayScore: match.AwayScore,
	}

	if err := s.producer.Publish(ctx, "match.ended", event); err != nil {
		return err
	}

	return nil
}

func (s *matchService) CancelMatch(id uint) error {
	if id == 0 {
		return errs.ErrInvalidMatchID
	}
	match, err := s.matchRepo.GetByID(id)
	if err != nil {
		return err
	}
	if match.Status != models.MatchStatusScheduled {
		return errors.New("match cannot be canceled")
	}

	match.Status = models.MatchStatusCancelled

	return s.matchRepo.Update(match)
}

func (s *matchService) UpdateScoreFromKafka(ctx context.Context, goalEvent models.GoalEvent) error {
	if goalEvent.MatchID == 0 || goalEvent.TeamID == 0 {
		return errors.New("invalid id")
	}

	return s.matchRepo.UpdateScoreFromKafka(ctx, goalEvent.MatchID, goalEvent.TeamID, goalEvent.NewScore)
}

func (s *matchService) GetActive() ([]models.Match, error) {

	return s.matchRepo.GetActive()
}
