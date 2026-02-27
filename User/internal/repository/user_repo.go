package repository

import (
	"User/internal/models"
	"errors"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	Delete(id uint) error
	Update(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByPhone(phone string) (*models.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (u *userRepo) GetByPhone(phone string)(*models.User,error){
	var user models.User
	if err := u.db.Where("phone = ?",phone).First(&user).Error; err != nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			return nil,ErrUserNotFound
		}
		return nil,err
	}
	return &user,nil
}
func (u *userRepo) Create(user *models.User) error {
	return u.db.Create(user).Error
}
func (u *userRepo) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := u.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err,gorm.ErrRecordNotFound){
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) Delete(id uint) error {
	result := u.db.Delete(&models.User{}, id)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
func (u *userRepo) Update(user *models.User) error {
	return  u.db.Model(&models.User{}).Where("id = ?",user.ID).Updates(user).Error
	 
}
func (u *userRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err,gorm.ErrRecordNotFound){
			return nil, ErrUserNotFound
		}
		return nil,err
	}
	return &user, nil
}
