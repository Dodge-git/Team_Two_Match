package services

import (
	"User/internal/auth"
	"User/internal/dto"
	"User/internal/models"
	"User/internal/repository"
	"errors"
	"strings"
)

type UserService interface {
	UserRegister(req dto.UserRegister) (dto.UserResponse, error)
	GetProfile(id uint) (*models.User, error)
	DeleteUser(id uint) error
	UpdateUser(id uint, req dto.UserUpdate) error
	LoginUser(req dto.UserLogin) (string, error)
}
type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}
func (u *userService) UserRegister(req dto.UserRegister) (dto.UserResponse, error) {

	name := strings.TrimSpace(req.Name)
	if len(name) < 3 {
		return dto.UserResponse{}, ErrInvalidName
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if len(email) < 7 {
		return dto.UserResponse{}, ErrInvalidEmail
	}

	phone := strings.TrimSpace(req.Phone)
	if len(phone) < 4 {
		return dto.UserResponse{}, ErrInvalidPhone
	}

	var favorite_sport *string
	if req.FavoriteSport != nil {
		fs := strings.TrimSpace(*req.FavoriteSport)
		if len(fs) > 1 {
			favorite_sport = &fs
		}
	}

	// Тут проверка и хэширования пароля

	checkEmail, err := u.userRepo.GetByEmail(email)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return dto.UserResponse{}, err
	}
	if checkEmail != nil {
		return dto.UserResponse{}, ErrUserAlreadyExist
	}
	checkPhone, err := u.userRepo.GetByPhone(phone)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return dto.UserResponse{}, err
	}
	if checkPhone != nil {
		return dto.UserResponse{}, ErrUserAlreadyExist
	}

	if err := auth.ValidatePassword(req.Password); err != nil {
		return dto.UserResponse{}, err
	}
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return dto.UserResponse{}, err
	}

	user := &models.User{
		Name:          name,
		Email:         email,
		Phone:         phone,
		PasswordHash:  hashedPassword,
		FavoriteSport: favorite_sport,
		Role:          models.RoleFan,
	}

	if err := u.userRepo.Create(user); err != nil {
		return dto.UserResponse{}, err
	}

	req_response := &dto.UserResponse{
		Name:          user.Name,
		Email:         user.Email,
		FavoriteSport: user.FavoriteSport,
	}
	return *req_response, nil
}

func (u *userService) GetProfile(id uint) (*models.User, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userService) DeleteUser(id uint) error {
	return u.userRepo.Delete(id); 
	
}
func (u *userService) UpdateUser(id uint, req dto.UserUpdate) error {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrNotFoundUser
		}
		return err
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if len(name) < 3 {
			return ErrInvalidName
		}
		user.Name = name
	}

	if req.Email != nil {
		email := strings.ToLower(strings.TrimSpace(*req.Email))
		if len(email) < 7 {
			return ErrInvalidEmail
		}
		if email != user.Email {
			checkEmail, err := u.userRepo.GetByEmail(email)
			if err == nil && checkEmail.ID != user.ID {
				return ErrAlreadyExistUserWithThisEmail
			}
			if err != nil {
				if !errors.Is(err, repository.ErrUserNotFound) {
					return err
				}
			}
		}
		user.Email = email
	}

	if req.Phone != nil {
		phone := strings.TrimSpace(*req.Phone)
		if len(phone) < 4 {
			return ErrInvalidPhone
		}
		user.Phone = phone
	}

	if req.FavoriteSport != nil {
		fs := strings.TrimSpace(*req.FavoriteSport)
		if len(fs) > 1 {
			user.FavoriteSport = &fs
		}
	}

	if req.OldPassword == nil && req.NewPassword == nil {
// сюда добавим log 
	} else {

		if req.OldPassword == nil || req.NewPassword == nil {
			return ErrRequiredNewAndOldPassword
		}
		if err := auth.CheckPassword(*req.OldPassword, user.PasswordHash); err != nil {
			return ErrOldPasswordIsntCorrect
		}
		if err := auth.ValidatePassword(*req.NewPassword); err != nil {
			return ErrNewPasswordIsntCorrect
		}
		password, err := auth.HashPassword(*req.NewPassword)
		if err != nil {
			return err
		}
		user.PasswordHash = password
	}

	return u.userRepo.Update(user)
}
func (u *userService) LoginUser(req dto.UserLogin) (string, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}
	if err := auth.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", ErrTokenWasntGenerate
	}
	return token, nil

}
