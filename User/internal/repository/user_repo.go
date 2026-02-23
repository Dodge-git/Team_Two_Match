package repository

import "User/internal/models"

type UserRepository interface {
	Create(user *models.User)error
	Login(user *models.User)
}