package services

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"insense-local/data"
	"insense-local/models"
	"time"
)

type UserServiceInterface interface {
	AddUser(ctx context.Context, user *models.User) error
	LoginUser(ctx context.Context, username, password string) (*models.User, error)
	UpdatePassword(ctx context.Context, id primitive.ObjectID, oldPassword, newPassword string) error
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
	UpdateUserInfo(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
}

type userService struct {
	userData       data.UserDataInterface
	contextTimeout time.Duration
}

func UserService(userData data.UserDataInterface, timeout time.Duration) UserServiceInterface {
	return &userService{
		userData:       userData,
		contextTimeout: timeout,
	}
}

func (us *userService) AddUser(ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, us.contextTimeout)
	defer cancel()

	// Check if user already exists
	existingUser, err := us.userData.FindByUserName(ctx, user.UserName)
	if err == nil && existingUser != nil {
		return errors.New("user already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return us.userData.Create(ctx, user)
}

func (us *userService) LoginUser(ctx context.Context, username, password string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, us.contextTimeout)
	defer cancel()

	user, err := us.userData.FindByUserName(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func (us *userService) UpdatePassword(ctx context.Context, id primitive.ObjectID, oldPassword, newPassword string) error {
	ctx, cancel := context.WithTimeout(ctx, us.contextTimeout)
	defer cancel()

	user, err := us.userData.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("incorrect old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return us.userData.Update(ctx, user)
}

func (us *userService) UpdateUserInfo(ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, us.contextTimeout)
	defer cancel()

	existingUser, err := us.userData.FindByID(ctx, user.ID)
	if err != nil {
		return err
	}

	if existingUser == nil {
		return errors.New("user not found")
	}

	user.Password = existingUser.Password // 保持原来的密码不变
	return us.userData.Update(ctx, user)
}

func (us *userService) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, us.contextTimeout)
	defer cancel()

	return us.userData.Delete(ctx, id)
}

func (us *userService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, us.contextTimeout)
	defer cancel()
	return us.userData.FindByID(ctx, id)
}
