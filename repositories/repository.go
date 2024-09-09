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
