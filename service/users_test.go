package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/toncek345/userservice/service"
	"github.com/toncek345/userservice/storage"
)

type testingTCtx struct{}

func testingTToCtx(ctx context.Context, t *testing.T) context.Context {
	return context.WithValue(ctx, testingTCtx{}, t)
}

func testingTFromCtx(ctx context.Context) *testing.T {
	return ctx.Value(testingTCtx{}).(*testing.T)
}

func TestAddUser(t *testing.T) {
	tests := []struct {
		name       string
		userIn     *service.AddUser
		userOut    *service.User
		service    service.UserService
		hasherFunc func(password []byte, cost int) ([]byte, error)
		isError    bool
	}{
		{
			name: "works",
			userIn: &service.AddUser{
				FirstName: "first_name",
				LastName:  "last_name",
				Email:     "email",
				Country:   "US",
				Password:  "password",
			},
			userOut: &service.User{
				ID:        "gen_id",
				FirstName: "first_name",
				LastName:  "last_name",
				Email:     "email",
				Country:   "US",
			},
			hasherFunc: func(password []byte, cost int) ([]byte, error) {
				return []byte("hashed_password!!!"), nil
			},
			service: &service.UserServiceImpl{
				UserStorage: &storage.MockUser{
					InsertUserFn: func(ctx context.Context, user *storage.InsertUser) (*storage.UserModel, error) {
						t := testingTFromCtx(ctx)
						if user.FirstName != "first_name" || user.LastName != "last_name" ||
							user.Email != "email" || user.Country != "US" || user.Password != "hashed_password!!!" {
							t.Fatal("user doesn't match insert user")
						}

						return &storage.UserModel{
							ID:        "gen_id",
							FirstName: "first_name",
							LastName:  "last_name",
							Email:     "email",
							Country:   "US",
							Password:  "hashed_password!!!",
						}, nil
					},
				},
			},
		},
		{
			name: "fails hash gen",
			userIn: &service.AddUser{
				FirstName: "first_name",
				LastName:  "last_name",
				Email:     "email",
				Country:   "US",
				Password:  "password",
			},
			isError: true,
			hasherFunc: func(password []byte, cost int) ([]byte, error) {
				return nil, fmt.Errorf("err")
			},
			service: &service.UserServiceImpl{},
		},
		{
			name: "fails user insert",
			userIn: &service.AddUser{
				FirstName: "first_name",
				LastName:  "last_name",
				Email:     "email",
				Country:   "US",
				Password:  "password",
			},
			isError: true,
			hasherFunc: func(password []byte, cost int) ([]byte, error) {
				return []byte("hashed_pw"), nil
			},
			service: &service.UserServiceImpl{
				UserStorage: &storage.MockUser{
					InsertUserFn: func(ctx context.Context, user *storage.InsertUser) (*storage.UserModel, error) {
						return nil, fmt.Errorf("err")
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), testingTCtx{}, t)
			service.GeneratePasswordHash = test.hasherFunc
			u, err := test.service.AddUser(ctx, test.userIn)
			if err != nil {
				if test.isError {
					return
				}
				t.Fatalf("unexpeted error: %s", err)
			}

			matchUsers(t, u, test.userOut)
		})
	}
}

func matchUsers(t *testing.T, u1, u2 *service.User) {
	if u1.ID != u2.ID ||
		u1.FirstName != u2.FirstName ||
		u1.LastName != u2.LastName ||
		u1.Email != u2.Email ||
		u1.Country != u2.Country {
		t.Fatal("users do not match")
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name    string
		isError bool
		idIn    string
		service service.UserService
	}{
		{
			name: "works",
			idIn: "id",
			service: &service.UserServiceImpl{
				&storage.MockUser{
					DeleteUserFn: func(ctx context.Context, id string) error {
						t := testingTFromCtx(ctx)
						if id != "id" {
							t.Fatal("wrong id")
						}

						return nil
					},
				},
			},
		},
		{
			name:    "fails",
			idIn:    "id",
			isError: true,
			service: &service.UserServiceImpl{
				&storage.MockUser{
					DeleteUserFn: func(ctx context.Context, id string) error {
						return fmt.Errorf("err")
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), testingTCtx{}, t)
			if err := test.service.DeleteUser(ctx, test.idIn); err != nil {
				if test.isError {
					return
				}
				t.Fatal("not expecting error")
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name       string
		userIn     *service.UpdateUser
		userOut    *service.User
		service    service.UserService
		hasherFunc func(password []byte, cost int) ([]byte, error)
		isError    bool
	}{
		{
			name: "works",
			userIn: &service.UpdateUser{
				ID:        "some_id",
				FirstName: "first_name",
				LastName:  "last_name",
				Email:     "email",
				Country:   "US",
				Password:  "password",
			},
			userOut: &service.User{
				ID:        "some_id",
				FirstName: "first_name",
				LastName:  "last_name",
				Email:     "email",
				Country:   "US",
			},
			hasherFunc: func(password []byte, cost int) ([]byte, error) {
				return []byte("hashed_password!!!"), nil
			},
			service: &service.UserServiceImpl{
				UserStorage: &storage.MockUser{
					UpdateUserFn: func(ctx context.Context, user *storage.UpdateUser) (*storage.UserModel, error) {
						t := testingTFromCtx(ctx)
						if user.ID != "some_id" || user.FirstName != "first_name" || user.LastName != "last_name" ||
							user.Email != "email" || user.Country != "US" || user.Password != "hashed_password!!!" {
							t.Fatal("user doesn't match update user")
						}

						return &storage.UserModel{
							ID:        "some_id",
							FirstName: "first_name",
							LastName:  "last_name",
							Email:     "email",
							Country:   "US",
							Password:  "hashed_password!!!",
						}, nil
					},
				},
			},
		},
		{
			name: "fails hash gen",
			userIn: &service.UpdateUser{
				FirstName: "first_name",
				LastName:  "last_name",
				Email:     "email",
				Country:   "US",
				Password:  "password",
			},
			isError: true,
			hasherFunc: func(password []byte, cost int) ([]byte, error) {
				return nil, fmt.Errorf("err")
			},
			service: &service.UserServiceImpl{},
		},
		{
			name: "fails user update",
			userIn: &service.UpdateUser{
				ID:        "some_id",
				FirstName: "first_name",
				LastName:  "last_name",
				Email:     "email",
				Country:   "US",
				Password:  "password",
			},
			isError: true,
			hasherFunc: func(password []byte, cost int) ([]byte, error) {
				return []byte("hashed_pw"), nil
			},
			service: &service.UserServiceImpl{
				UserStorage: &storage.MockUser{
					UpdateUserFn: func(ctx context.Context, user *storage.UpdateUser) (*storage.UserModel, error) {
						return nil, fmt.Errorf("err")
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), testingTCtx{}, t)
			service.GeneratePasswordHash = test.hasherFunc
			u, err := test.service.UpdateUser(ctx, test.userIn)
			if err != nil {
				if test.isError {
					return
				}
				t.Fatalf("unexpeted error: %s", err)
			}

			matchUsers(t, u, test.userOut)
		})
	}

}
