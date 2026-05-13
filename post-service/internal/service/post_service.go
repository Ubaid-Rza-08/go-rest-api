package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/Ubaid-Rza-08/post-service/internal/models"
	"github.com/Ubaid-Rza-08/post-service/internal/repository"
)

var ErrForbidden = errors.New("forbidden")

type PostService struct {
	repo  *repository.PostRepository
	redis *goredis.Client
}

func NewPostService(
	repo *repository.PostRepository,
	redis *goredis.Client,
) *PostService {

	return &PostService{
		repo:  repo,
		redis: redis,
	}
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

	if err != nil {
		return nil, err
	}

	// INVALIDATE CACHE
	s.redis.Del(context.Background(), "posts:all")

	return post, nil
}

func (s *PostService) GetAll() ([]models.Post, error) {

	cacheKey := "posts:all"

	// --------------------------------
	// CHECK REDIS CACHE
	// --------------------------------

	cachedPosts, err := s.redis.Get(
		context.Background(),
		cacheKey,
	).Result()

	if err == nil {

		var posts []models.Post

		json.Unmarshal([]byte(cachedPosts), &posts)

		log.Println("POSTS FETCHED FROM REDIS")

		return posts, nil
	}

	// --------------------------------
	// CACHE MISS → DATABASE
	// --------------------------------

	posts, err := s.repo.GetAll()

	if err != nil {
		return nil, err
	}

	// --------------------------------
	// SAVE TO REDIS
	// --------------------------------

	jsonData, _ := json.Marshal(posts)

	s.redis.Set(
		context.Background(),
		cacheKey,
		jsonData,
		5*time.Minute,
	)

	log.Println("POSTS FETCHED FROM POSTGRES")

	return posts, nil
}

func (s *PostService) GetByID(id string) (*models.Post, error) {

	cacheKey := "post:" + id

	// --------------------------------
	// CHECK REDIS
	// --------------------------------

	cachedPost, err := s.redis.Get(
		context.Background(),
		cacheKey,
	).Result()

	if err == nil {

		var post models.Post

		json.Unmarshal([]byte(cachedPost), &post)

		log.Println("POST FETCHED FROM REDIS")

		return &post, nil
	}

	// --------------------------------
	// CACHE MISS → DATABASE
	// --------------------------------

	post, err := s.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	jsonData, _ := json.Marshal(post)

	s.redis.Set(
		context.Background(),
		cacheKey,
		jsonData,
		5*time.Minute,
	)

	log.Println("POST FETCHED FROM POSTGRES")

	return post, nil
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

	if err != nil {
		return nil, err
	}

	// INVALIDATE CACHE
	s.redis.Del(context.Background(), "posts:all")
	s.redis.Del(context.Background(), "post:"+id)

	return post, nil
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

	err = s.repo.Delete(id)

	if err != nil {
		return err
	}

	// INVALIDATE CACHE
	s.redis.Del(context.Background(), "posts:all")
	s.redis.Del(context.Background(), "post:"+id)

	return nil
}