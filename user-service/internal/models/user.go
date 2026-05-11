// internal/models/user.go
package models

import (
	"time"
)

// Role type for user roles.
type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// User is the database model for the users table.
type User struct {
	ID           string     `db:"id"            json:"id"`
	Name         string     `db:"name"          json:"name"`
	Email        string     `db:"email"         json:"email"`
	PasswordHash string     `db:"password_hash" json:"-"` // never serialized to JSON
	Role         Role       `db:"role"          json:"role"`
	IsActive     bool       `db:"is_active"     json:"is_active"`
	CreatedAt    time.Time  `db:"created_at"    json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"    json:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"    json:"-"`
}

// ── Request / Response DTOs ────────────────────────────────────

// RegisterRequest is the payload for POST /auth/register.
type RegisterRequest struct {
	Name     string `json:"name"     binding:"required,min=2,max=100"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

// LoginRequest is the payload for POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is returned on successful login.
type LoginResponse struct {
	AccessToken string   `json:"access_token"`
	TokenType   string   `json:"token_type"`
	ExpiresIn   int      `json:"expires_in"` // seconds
	User        UserResponse `json:"user"`
}

// UserResponse is the public-facing user shape (no password hash).
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateUserRequest is the payload for PUT /users/:id.
type UpdateUserRequest struct {
	Name     string `json:"name"     binding:"omitempty,min=2,max=100"`
	Email    string `json:"email"    binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=8,max=72"`
	IsActive *bool  `json:"is_active"`
}

// ToResponse maps a User model to a UserResponse DTO.
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}