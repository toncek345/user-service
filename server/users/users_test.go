package users_test

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/toncek345/userservice/proto"
	"github.com/toncek345/userservice/server/users"
	"github.com/toncek345/userservice/service"
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
		userIn     *pb.AddUserMessage
		userOut    *pb.User
		isError    bool
		userServer users.UserServer
	}{
		{
			name: "works",
			userIn: &pb.AddUserMessage{
				FirstName: "first_name",
				LastName:  "last_name",
				Email:     "email",
				Country:   "US",
				Password:  "password",
			},
			userOut: &pb.User{
				Id:        "gen_id",
				FirstName: "first_name",
				LastName:  "last_name",
				Email:     "email",
				Country:   "US",
			},
			userServer: users.UserServer{
				UserService: &service.UsersMock{
					AddUserFn: func(ctx context.Context, user *service.AddUser) (*service.User, error) {
						t := testingTFromCtx(ctx)
						if user.FirstName != "first_name" || user.LastName != "last_name" ||
							user.Email != "email" || user.Country != "US" || user.Password != "password" {
							t.Fatal("user doesn't match")
						}

						return &service.User{
							ID:        "gen_id",
							FirstName: "first_name",
							LastName:  "last_name",
							Email:     "email",
							Country:   "US",
						}, nil
					},
				},
			},
		},
		{
			name:    "fails",
			isError: true,
			userIn:  &pb.AddUserMessage{},
			userServer: users.UserServer{
				UserService: &service.UsersMock{
					AddUserFn: func(ctx context.Context, user *service.AddUser) (*service.User, error) {
						return nil, fmt.Errorf("error")
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := testingTToCtx(context.Background(), t)

			u, err := test.userServer.AddUser(ctx, test.userIn)
			if err != nil {
				if test.isError {
					return
				}

				t.Fatalf("unexpected error: %s", err)
			}

			matchUsers(t, u, test.userOut)
		})
	}
}

func matchUsers(t *testing.T, u1, u2 *pb.User) {
	if u1.Id != u2.Id ||
		u1.FirstName != u2.FirstName ||
		u1.LastName != u2.LastName ||
		u1.Email != u2.Email ||
		u1.Country != u2.Country {
		t.Fatal("users do not match")
	}
}

func TestDeleteUserr(t *testing.T) {
	tests := []struct {
		name       string
		isError    bool
		idIn       string
		userServer users.UserServer
	}{
		{
			name: "works",
			idIn: "idasdf",
			userServer: users.UserServer{
				UserService: &service.UsersMock{
					DeleteUserFn: func(ctx context.Context, id string) error {
						t := testingTFromCtx(ctx)
						if id != "idasdf" {
							t.Fatal("id doesn't match")
						}
						return nil
					},
				},
			},
		},
		{
			name:    "fails",
			idIn:    "idasdf",
			isError: true,
			userServer: users.UserServer{
				UserService: &service.UsersMock{
					DeleteUserFn: func(ctx context.Context, id string) error {
						t := testingTFromCtx(ctx)
						if id != "idasdf" {
							t.Fatal("id doesn't match")
						}
						return fmt.Errorf("error")
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := testingTToCtx(context.Background(), t)
			if _, err := test.userServer.DeleteUser(ctx, &pb.DeleteUserMessage{Id: test.idIn}); err != nil {
				if test.isError {
					return
				}
				t.Fatal("unexpected error")
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name       string
		userIn     *pb.UpdateUserMessage
		userOut    *pb.User
		isError    bool
		userServer users.UserServer
	}{
		{
			name: "works",
			userIn: &pb.UpdateUserMessage{
				Id:        "some_id",
				FirstName: "first_name",
				LastName:  "last_name",
				Email:     "email",
				Country:   "US",
				Password:  "password",
			},
			userOut: &pb.User{
				Id:        "some_id",
				FirstName: "first_name",
				LastName:  "last_name",
				Email:     "email",
				Country:   "US",
			},
			userServer: users.UserServer{
				UserService: &service.UsersMock{
					UpdateUserFn: func(ctx context.Context, user *service.UpdateUser) (*service.User, error) {
						t := testingTFromCtx(ctx)
						if user.ID != "some_id" || user.FirstName != "first_name" || user.LastName != "last_name" ||
							user.Email != "email" || user.Country != "US" || user.Password != "password" {
							t.Fatal("user doesn't match")
						}

						return &service.User{
							ID:        "some_id",
							FirstName: "first_name",
							LastName:  "last_name",
							Email:     "email",
							Country:   "US",
						}, nil
					},
				},
			},
		},
		{
			name:    "fails",
			isError: true,
			userIn:  &pb.UpdateUserMessage{},
			userServer: users.UserServer{
				UserService: &service.UsersMock{
					UpdateUserFn: func(ctx context.Context, user *service.UpdateUser) (*service.User, error) {
						return nil, fmt.Errorf("error")
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := testingTToCtx(context.Background(), t)

			u, err := test.userServer.UpdateUser(ctx, test.userIn)
			if err != nil {
				if test.isError {
					return
				}

				t.Fatalf("unexpected error: %s", err)
			}

			matchUsers(t, u, test.userOut)
		})
	}
}
