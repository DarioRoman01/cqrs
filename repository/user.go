package repository

import (
	"context"

	"github.com/DarioRoman01/cqrs/models"
)

type UserRepository interface {
	Repository[*models.User]
	Login(ctx context.Context, name, password string) (*models.User, error)
}

var userRepo UserRepository

func SetUserRepository(repo UserRepository) {
	userRepo = repo
}

func CloseUserRepo() error {
	return userRepo.Close()
}

func InsertUser(ctx context.Context, user *models.User) error {
	return userRepo.Insert(ctx, user)
}

func ListUsers(ctx context.Context) ([]*models.User, error) {
	return userRepo.List(ctx)
}

func GetUser(ctx context.Context, id string) (*models.User, error) {
	return userRepo.Get(ctx, id)
}

func DeleteUser(ctx context.Context, id string) error {
	return userRepo.Delete(ctx, id)
}

func UpdateUser(ctx context.Context, user *models.User) error {
	return userRepo.Update(ctx, user)
}

func Login(ctx context.Context, name, password string) (*models.User, error) {
	return userRepo.Login(ctx, name, password)
}
