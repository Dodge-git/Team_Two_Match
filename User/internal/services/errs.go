package services

import "errors"

var (
	ErrInvalidEmail                  = errors.New("invalid email")
	ErrInvalidPassword               = errors.New("invalid password")
	ErrNotFoundUser                  = errors.New("user not found")
	ErrInvalidValue                  = errors.New("invalid value")
	ErrInvalidName                   = errors.New("invalid name")
	ErrInvalidPhone                  = errors.New("invalid phone")
	ErrNotFoundUserWithThisEmail     = errors.New("not found user with this email")
	ErrNotFoundUserWithThisPhone     = errors.New("not found user with this phone")
	ErrUserAlreadyExist              = errors.New("user already exist")
	ErrInvalidFormat                 = errors.New("invalid format")
	ErrTokenWasntGenerate            = errors.New("token wasn,t generate")
	ErrAlreadyExistUserWithThisEmail = errors.New("already exist user with this email")
	ErrInvalidCredentials            = errors.New("invalid credentials")
	ErrRequiredNewAndOldPassword     = errors.New("new and old password required")
	ErrOldPasswordIsntCorrect        = errors.New("old password isnt correct")
	ErrNewPasswordIsntCorrect        = errors.New("new password isnt correct")
)
