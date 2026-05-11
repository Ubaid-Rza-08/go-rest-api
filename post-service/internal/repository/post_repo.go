package repository

import (
	"github.com/jmoiron/sqlx"

	"github.com/Ubaid-Rza-08/post-service/internal/models"
)

type PostRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *models.Post) error {
	query := `
	INSERT INTO posts (user_id, title, content)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		post.UserID,
		post.Title,
		post.Content,
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
}

func (r *PostRepository) GetAll() ([]models.Post, error) {
	var posts []models.Post

	query := `
	SELECT * FROM posts
	WHERE deleted_at IS NULL
	ORDER BY created_at DESC
	`

	err := r.db.Select(&posts, query)

	return posts, err
}

func (r *PostRepository) GetByID(id string) (*models.Post, error) {
	var post models.Post

	query := `
	SELECT * FROM posts
	WHERE id=$1 AND deleted_at IS NULL
	`

	err := r.db.Get(&post, query, id)

	return &post, err
}

func (r *PostRepository) Update(post *models.Post) error {
	query := `
	UPDATE posts
	SET title=$1, content=$2, updated_at=NOW()
	WHERE id=$3
	`

	_, err := r.db.Exec(
		query,
		post.Title,
		post.Content,
		post.ID,
	)

	return err
}

func (r *PostRepository) Delete(id string) error {
	_, err := r.db.Exec(
		`UPDATE posts SET deleted_at=NOW() WHERE id=$1`,
		id,
	)

	return err
}