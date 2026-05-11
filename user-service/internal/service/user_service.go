// internal/service/user_service.go
package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Ubaid-Rza-08/go-rest-api/internal/config"
	"github.com/Ubaid-Rza-08/go-rest-api/internal/models"
	"github.com/Ubaid-Rza-08/go-rest-api/internal/repository"
	"github.com/Ubaid-Rza-08/go-rest-api/internal/utils"
)

// Common service errors — callers should check with errors.Is.
var (
	ErrNotFound      = errors.New("user not found")
	ErrInvalidCreds  = errors.New("invalid email or password")
	ErrEmailTaken    = errors.New("email is already registered")
	ErrForbidden     = errors.New("insufficient permissions")
)

// UserService handles business logic for user operations.
type UserService interface {
	Register(req *models.RegisterRequest) (*models.UserResponse, error)
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
	GetByID(id string) (*models.UserResponse, error)
	GetAll(page, pageSize int) ([]models.UserResponse, int, error)
	Update(id string, req *models.UpdateUserRequest, callerID, callerRole string) (*models.UserResponse, error)
	Delete(id, callerID, callerRole string) error
}

type userService struct {
	repo repository.UserRepository
	cfg  *config.Config
}

// NewUserService creates a UserService with dependencies injected.
func NewUserService(repo repository.UserRepository, cfg *config.Config) UserService {
	return &userService{repo: repo, cfg: cfg}
}

// Register creates a new user with a hashed password.
func (s *userService) Register(req *models.RegisterRequest) (*models.UserResponse, error) {
	hash, err := utils.HashPassword(req.Password, s.cfg.App.BcryptCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
		Role:         models.RoleUser,
	}

	created, err := s.repo.Create(user)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return nil, ErrEmailTaken
		}
		return nil, fmt.Errorf("creating user: %w", err)
	}

	resp := created.ToResponse()
	return &resp, nil
}

// Login validates credentials and returns a signed JWT.
func (s *userService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidCreds // Don't reveal which field was wrong
		}
		return nil, fmt.Errorf("finding user: %w", err)
	}

	if !user.IsActive {
		return nil, ErrInvalidCreds
	}

	if err := utils.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, ErrInvalidCreds
	}

	expiry := s.cfg.JWT.AccessExpiry
	token, err := utils.GenerateAccessToken(
		user.ID, user.Email, string(user.Role),
		s.cfg.JWT.Secret, expiry,
	)
	if err != nil {
		return nil, fmt.Errorf("generating token: %w", err)
	}

	resp := user.ToResponse()
	return &models.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(expiry / time.Second),
		User:        resp,
	}, nil
}

// GetByID returns a user's public profile.
func (s *userService) GetByID(id string) (*models.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	resp := user.ToResponse()
	return &resp, nil
}

// GetAll returns a paginated list of users (admin only — enforce in handler).
func (s *userService) GetAll(page, pageSize int) ([]models.UserResponse, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	users, total, err := s.repo.FindAll(pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.UserResponse, len(users))
	for i, u := range users {
		responses[i] = u.ToResponse()
	}
	return responses, total, nil
}

// Update modifies a user record. Users can update themselves; admins can update anyone.
func (s *userService) Update(id string, req *models.UpdateUserRequest, callerID, callerRole string) (*models.UserResponse, error) {
	if callerID != id && callerRole != string(models.RoleAdmin) {
		return nil, ErrForbidden
	}

	fields := map[string]interface{}{}
	if req.Name != "" {
		fields["name"] = req.Name
	}
	if req.Email != "" {
		fields["email"] = req.Email
	}
	if req.Password != "" {
		hash, err := utils.HashPassword(req.Password, s.cfg.App.BcryptCost)
		if err != nil {
			return nil, fmt.Errorf("hashing new password: %w", err)
		}
		fields["password_hash"] = hash
	}
	if req.IsActive != nil {
		fields["is_active"] = *req.IsActive
	}

	updated, err := s.repo.Update(id, fields)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return nil, ErrEmailTaken
		}
		return nil, err
	}

	resp := updated.ToResponse()
	return &resp, nil
}

// Delete soft-deletes a user. Only admins or the user themselves may delete.
func (s *userService) Delete(id, callerID, callerRole string) error {
	if callerID != id && callerRole != string(models.RoleAdmin) {
		return ErrForbidden
	}
	if err := s.repo.SoftDelete(id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}