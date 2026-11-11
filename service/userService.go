package service

import (
	"database/sql"
	"errors"

	"github.com/SandroS22/chi-test/cmd/request"
	"github.com/SandroS22/chi-test/cmd/response"
	"github.com/SandroS22/chi-test/cmd/util"
	"github.com/SandroS22/chi-test/db"
	"github.com/SandroS22/chi-test/repository"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidUserData = errors.New("invalid user data")
)

type UserService interface {
	CreateUser(req request.CreateUserRequest) (*response.CreateUserResponse, error)
	GetUserByID(id int) (*response.UserResponse, error)
	GetAllUsers() ([]response.UserResponse, error)
	UpdateUser(id int, req request.UpdateUserRequest) (*response.UserResponse, error)
	DeleteUser(id int) error
	GetUserByUsername(username string) (*response.UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(req request.CreateUserRequest) (*response.CreateUserResponse, error) {
	userIdentifier := req.Email
	if userIdentifier == "" {
		userIdentifier = req.User
	}
	if req.Name == "" || userIdentifier == "" || req.Password == "" {
		return nil, ErrInvalidUserData
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &db.User{
		Name:     req.Name,
		User:     userIdentifier,
		Password: hashedPassword,
		Role:     req.Role,
	}

	id, err := s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return response.NewCreateUserResponse(id, user.Name, userIdentifier, hashedPassword, user.Role), nil
}

func (s *userService) GetUserByID(id int) (*response.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) || errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return response.NewUserResponse(user.ID, user.Name, user.User, user.Role), nil
}

func (s *userService) GetAllUsers() ([]response.UserResponse, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	res := make([]response.UserResponse, 0, len(users))
	for _, u := range users {
		res = append(res, *response.NewUserResponse(u.ID, u.Name, u.User, u.Role))
	}
	return res, nil
}

func (s *userService) UpdateUser(id int, req request.UpdateUserRequest) (*response.UserResponse, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) || errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	name := req.Name
	if name == "" {
		name = existing.Name
	}

	userIdentifier := req.Email
	if userIdentifier == "" {
		userIdentifier = req.User
	}
	if userIdentifier == "" {
		userIdentifier = existing.User
	}

	role := req.Role
	if role == "" {
		role = existing.Role
	}

	var hashedPassword string
	if req.Password != "" {
		hashed, err := util.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		hashedPassword = hashed
	}

	updatedUser := &db.User{
		ID:       id,
		Name:     name,
		User:     userIdentifier,
		Password: hashedPassword,
		Role:     role,
	}

	if err := s.repo.Update(updatedUser); err != nil {
		return nil, err
	}

	return response.NewUserResponse(id, name, userIdentifier, role), nil
}

func (s *userService) DeleteUser(id int) error {
	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) || errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

func (s *userService) GetUserByUsername(username string) (*response.UserResponse, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) || errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return response.NewUserResponse(user.ID, user.Name, user.User, user.Role), nil
}
