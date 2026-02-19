package errs

import "errors"

var (
	ErrPlayerNotFound  = errors.New("player not found")
	ErrInvalidPlayerID = errors.New("invalid player id")

	ErrTeamNotFound    = errors.New("team not found")
	ErrInvalidTeamID   = errors.New("invalid team id")
	ErrInvalidTeamName = errors.New("invalid team name")
	

	ErrMatchNotFound      = errors.New("match not found")
	ErrInvalidMatchStatus = errors.New("invalid match status")
	ErrInvalidMatchID     = errors.New("invalid match id")

	ErrSportNotFound   = errors.New("sport not found")
	ErrTeamsNotInSport = errors.New("one or both teams do not belong to the specified sport")

	ErrInvalidGoalEvent = errors.New("invalid goal event: team is not part of the match or match is not live")

	ErrInvalidSport = errors.New("invalid sport")
)
