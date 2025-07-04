package repository

import (
	"movie-api/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(username, password string) (*model.User, error)
	FindUserByUsername(username string) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) CreateUser(username string, password string) (*model.User, error) {
	var user model.User
	user.Username = username
	user.Password = password
	if err := u.db.Create(&user).Error; err != nil {
		return nil, err
	}
	emptyID := [16]byte{}
	if user.ID == emptyID {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}

func (u *userRepository) FindUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	if user.ID == [16]byte{} {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}
