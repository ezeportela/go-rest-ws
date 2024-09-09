package repositories

import (
	"context"

	"github.com/ezeportela/go-rest-ws/models"
)

type UserRepository interface {
	Close()
	InsertUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	InsertPost(ctx context.Context, post *models.Post) error
	GetPostById(ctx context.Context, id string) (*models.Post, error)
	UpdatePost(ctx context.Context, post *models.Post) error
	DeletePost(ctx context.Context, id string, userId string) error
	ListPosts(ctx context.Context, limit uint64, page uint64) ([]*models.Post, error)
}

var userRepository UserRepository

func SetUserRepository(repo UserRepository) {
	userRepository = repo
}

func InsertUser(ctx context.Context, user *models.User) error {
	return userRepository.InsertUser(ctx, user)
}

func GetUserById(ctx context.Context, id string) (*models.User, error) {
	return userRepository.GetUserById(ctx, id)
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return userRepository.GetUserByEmail(ctx, email)
}

func Close() {
	userRepository.Close()
}

func InsertPost(ctx context.Context, post *models.Post) error {
	return userRepository.InsertPost(ctx, post)
}

func GetPostById(ctx context.Context, id string) (*models.Post, error) {
	return userRepository.GetPostById(ctx, id)
}

func UpdatePost(ctx context.Context, post *models.Post) error {
	return userRepository.UpdatePost(ctx, post)
}

func DeletePost(ctx context.Context, id string, userId string) error {
	return userRepository.DeletePost(ctx, id, userId)
}

func ListPosts(ctx context.Context, limit uint64, page uint64) ([]*models.Post, error) {
	return userRepository.ListPosts(ctx, limit, page)
}
