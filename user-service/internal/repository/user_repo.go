// internal/repository/user_repo.go
package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Ubaid-Rza-08/go-rest-api/internal/models"
)

// ErrNotFound is returned when a record is not found.
var ErrNotFound = errors.New("record not found")

// ErrDuplicateEmail is returned when the email is already registered.
var ErrDuplicateEmail = errors.New("email already exists")

// UserRepository defines the contract for user data operations.
type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	FindByID(id string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindAll(limit, offset int) ([]models.User, int, error)
	Update(id string, fields map[string]interface{}) (*models.User, error)
	SoftDelete(id string) error
}

type userRepo struct {
	db *sqlx.DB
}

// NewUserRepository returns a UserRepository backed by PostgreSQL.
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

// Create inserts a new user and returns the created record.
func (r *userRepo) Create(user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (name, email, password_hash, role)
		VALUES (:name, :email, :password_hash, :role)
		RETURNING id, name, email, password_hash, role, is_active, created_at, updated_at`

	rows, err := r.db.NamedQuery(query, user)
	if err != nil {
		if isPgUniqueViolation(err) {
			return nil, ErrDuplicateEmail
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	defer rows.Close()

	var created models.User
	if rows.Next() {
		if err := rows.StructScan(&created); err != nil {
			return nil, fmt.Errorf("scan created user: %w", err)
		}
	}
	return &created, nil
}

// FindByID retrieves an active user by UUID.
func (r *userRepo) FindByID(id string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL`
	if err := r.db.Get(&user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &user, nil
}

// FindByEmail retrieves an active user by email address.
func (r *userRepo) FindByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL`
	if err := r.db.Get(&user, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return &user, nil
}

// FindAll returns a paginated list of active users and the total count.
func (r *userRepo) FindAll(limit, offset int) ([]models.User, int, error) {
	var users []models.User
	query := `SELECT * FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	if err := r.db.Select(&users, query, limit, offset); err != nil {
		return nil, 0, fmt.Errorf("find all users: %w", err)
	}

	var total int
	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
	if err := r.db.Get(&total, countQuery); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}
	return users, total, nil
}

// Update applies a map of column→value updates and returns the updated user.
// Only non-zero-value fields are updated.
func (r *userRepo) Update(id string, fields map[string]interface{}) (*models.User, error) {
	if len(fields) == 0 {
		return r.FindByID(id)
	}

	// Build dynamic SET clause safely
	setClauses := ""
	args := []interface{}{}
	i := 1
	for col, val := range fields {
		if setClauses != "" {
			setClauses += ", "
		}
		setClauses += fmt.Sprintf("%s = $%d", col, i)
		args = append(args, val)
		i++
	}
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE users SET %s
		WHERE id = $%d AND deleted_at IS NULL
		RETURNING id, name, email, password_hash, role, is_active, created_at, updated_at`,
		setClauses, i)

	var user models.User
	if err := r.db.Get(&user, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		if isPgUniqueViolation(err) {
			return nil, ErrDuplicateEmail
		}
		return nil, fmt.Errorf("update user: %w", err)
	}
	return &user, nil
}

// SoftDelete marks a user as deleted without physically removing the row.
func (r *userRepo) SoftDelete(id string) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("soft delete user: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// isPgUniqueViolation checks the PostgreSQL error code for unique constraint violations.
func isPgUniqueViolation(err error) bool {
	return err != nil && len(err.Error()) > 0 &&
		contains(err.Error(), "unique constraint") || contains(err.Error(), "23505")
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && findSubstring(s, sub))
}

func findSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}