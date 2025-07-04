package usecase

import (
	"errors"
	"movie-api/internal/model"
	"movie-api/internal/repository"
	"movie-api/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(payload *model.RegisterPayload) (*model.User, error)
	Login(payload *model.LoginPayload) (map[string]string, error)
}

type authUsecase struct {
	userRepo repository.UserRepository
}

func NewAuthUsecase(userRepo repository.UserRepository) AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
	}
}

func (a *authUsecase) Register(payload *model.RegisterPayload) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: payload.Username,
		Password: string(hashedPassword),
	}

	return a.userRepo.CreateUser(user.Username, user.Password)
}

func (a *authUsecase) Login(payload *model.LoginPayload) (map[string]string, error) {
	user, err := a.userRepo.FindUserByUsername(payload.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	return token, nil
}
