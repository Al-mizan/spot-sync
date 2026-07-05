package user

import (
	"errors"
	"fmt"
	"spotsync/internal/apperror"
	"spotsync/internal/auth"
	"spotsync/internal/domain/user/dto"
)

var ErrInvalidCredentials = fmt.Errorf("invalid email or password")

// Service defines the contract for user business logic.
type Service interface {
	Register(req dto.RegisterRequest) (*dto.UserResponse, error)
	Login(req dto.LoginRequest) (*dto.LoginData, error)
}

type service struct {
	repo       Repository
	jwtService auth.JWTService
}

func NewService(repo Repository, jwtService auth.JWTService) Service {
	return &service{repo, jwtService}
}

func (s *service) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {

	// default role to driver if not provided
	role := req.Role
	if role == "" {
		role = "driver"
	}

	user := User{
		Name:  req.Name,
		Email: req.Email,
		Role:  role,
	}

	// hash password
	err := user.HashPassword(req.Password)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to hash password")
	}

	err = s.repo.Create(&user)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			return nil, apperror.NewBadRequest(err, "user with this email already exists")
		}
		return nil, apperror.NewInternal(err, "failed to create user")
	}

	response := &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return response, nil
}

func (s *service) Login(req dto.LoginRequest) (*dto.LoginData, error) {
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch user")
	}

	if user == nil {
		return nil, apperror.NewUnauthorized(ErrInvalidCredentials, "invalid email or password")
	}

	// check password
	err = user.CheckPassword(req.Password)
	if err != nil {
		return nil, apperror.NewUnauthorized(ErrInvalidCredentials, "invalid email or password")
	}

	// generate token with id and role
	token, err := s.jwtService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to generate token")
	}

	response := &dto.LoginData{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}

	return response, nil
}
