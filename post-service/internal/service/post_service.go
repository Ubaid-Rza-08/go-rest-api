package service

import (
	"errors"

	"github.com/Ubaid-Rza-08/post-service/internal/models"
	"github.com/Ubaid-Rza-08/post-service/internal/repository"
)

var ErrForbidden = errors.New("forbidden")

type PostService struct {
	repo *repository.PostRepository
}

func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) Create(
	req *models.CreatePostRequest,
	userID string,
) (*models.Post, error) {

	post := &models.Post{
		UserID: userID,
		Title:  req.Title,
		Content: req.Content,
	}

	err := s.repo.Create(post)

	return post, err
}

func (s *PostService) GetAll() ([]models.Post, error) {
	return s.repo.GetAll()
}

func (s *PostService) GetByID(id string) (*models.Post, error) {
	return s.repo.GetByID(id)
}

func (s *PostService) Update(
	id string,
	req *models.UpdatePostRequest,
	callerID string,
) (*models.Post, error) {

	post, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if post.UserID != callerID {
		return nil, ErrForbidden
	}

	if req.Title != nil {
		post.Title = *req.Title
	}

	if req.Content != nil {
		post.Content = *req.Content
	}

	err = s.repo.Update(post)

	return post, err
}

func (s *PostService) Delete(
	id string,
	callerID string,
) error {

	post, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if post.UserID != callerID {
		return ErrForbidden
	}

	return s.repo.Delete(id)
}