package service

import (
	"context"
	"fmt"
	"time"

	"github.com/toncek345/userservice/storage"

	"golang.org/x/crypto/bcrypt"
)

var _ UserService = (*UserServiceImpl)(nil)
var _ UserService = (*UsersMock)(nil)

type UserService interface {
	AddUser(ctx context.Context, user *AddUser) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, user *UpdateUser) (*User, error)
	// SearchUser returns a list of users and optinally filters them by country.
	SearchUser(ctx context.Context, page, page_size int64, country string) ([]*User, error)
}

type UserServiceImpl struct {
	UserStorage storage.UserStorage
}

var GeneratePasswordHash func(password []byte, cost int) ([]byte, error) = bcrypt.GenerateFromPassword

type User struct {
	// ID is represented in UUID.
	ID        string
	FirstName string
	LastName  string
	Email     string
	Country   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func storageUserToServiceUser(u *storage.UserModel) *User {
	return &User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Country:   u.Country,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type AddUser struct {
	FirstName string
	LastName  string
	Email     string
	Country   string
	Password  string
}

func (u *UserServiceImpl) AddUser(ctx context.Context, user *AddUser) (*User, error) {
	hashedPw, err := GeneratePasswordHash([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	storageUser, err := u.UserStorage.InsertUser(
		ctx,
		&storage.InsertUser{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Country:   user.Country,
			Password:  string(hashedPw),
		})
	if err != nil {
		return nil, fmt.Errorf("adding user: %w", err)
	}

	return storageUserToServiceUser(storageUser), nil
}

func (u *UserServiceImpl) DeleteUser(ctx context.Context, id string) error {
	if err := u.UserStorage.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("user storage: %w", err)
	}

	return nil
}

type UpdateUser struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Country   string
	Password  string
}

func (u *UserServiceImpl) UpdateUser(ctx context.Context, user *UpdateUser) (*User, error) {
	hashedPw, err := GeneratePasswordHash([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	storageUser, err := u.UserStorage.UpdateUser(
		ctx,
		&storage.UpdateUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Country:   user.Country,
			Password:  string(hashedPw),
		})
	if err != nil {
		return nil, fmt.Errorf("updating user: %w", err)
	}

	return storageUserToServiceUser(storageUser), nil
}

func (u *UserServiceImpl) SearchUser(ctx context.Context, page, page_size int64, country string) ([]*User, error) {
	// TODO: extract pagination in other file/module

	usersS, err := u.UserStorage.SearchUser(ctx, &storage.Filters{Country: country}, page_size*(page-1), page_size)
	if err != nil {
		return nil, fmt.Errorf("search user storage: %w", err)
	}

	users := make([]*User, 0, len(usersS))
	for _, u := range usersS {
		users = append(users, storageUserToServiceUser(u))
	}

	return users, nil
}
