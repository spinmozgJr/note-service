package service

import (
	"context"
	"github.com/spinmozgJr/note-service/internal/config"
	"github.com/spinmozgJr/note-service/internal/models"
	"github.com/spinmozgJr/note-service/internal/storage"
	"github.com/spinmozgJr/note-service/pkg/auth"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
)

type UserService struct {
	UserRepository storage.Storage
	TokenManager   auth.TokenManager
	Config         *config.Config
}

func NewUserService(userRepository storage.Storage, manager auth.TokenManager, config *config.Config) UserService {
	return UserService{
		UserRepository: userRepository,
		TokenManager:   manager,
		Config:         config,
	}
}

func (s *UserService) SignIn(ctx context.Context,
	input UserInput) (*models.BaseResponse[models.AuthData], error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	userID, err := s.UserRepository.AddUser(ctx, strings.ToLower(input.Username), string(hash))
	if err != nil {
		return nil, err
	}
	jwt, err := s.TokenManager.NewJWT(strconv.Itoa(userID), s.Config.JwtTTLDuration)
	if err != nil {
		return nil, err
	}
	response := &models.BaseResponse[models.AuthData]{
		Data: &models.AuthData{
			Username:    strings.ToLower(input.Username),
			AccessToken: jwt,
		},
	}
	return response, nil
}

func (s *UserService) Login(ctx context.Context, input UserInput) (*models.BaseResponse[models.AuthData], error) {
	user, err := s.UserRepository.GetUserByUsername(ctx, input.Username)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, err
	}
	jwt, err := s.TokenManager.NewJWT(strconv.Itoa(user.ID), s.Config.JwtTTLDuration)
	if err != nil {
		return nil, err
	}
	response := &models.BaseResponse[models.AuthData]{
		Data: &models.AuthData{
			Username:    strings.ToLower(input.Username),
			AccessToken: jwt,
		},
	}
	return response, nil
}
