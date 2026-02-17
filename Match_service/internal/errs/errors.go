package errs

import "errors"

var (
	ErrPlayerNotFound  = errors.New("player not found")
	ErrInvalidPlayerID = errors.New("invalid player id")

	ErrTeamNotFound  = errors.New("team not found")
	ErrInvalidTeamID = errors.New("invalid team id")

	
	ErrMatchNotFound = errors.New("match not found")
)
