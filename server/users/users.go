package users

import (
	"context"
	"errors"
	"log"

	pb "github.com/toncek345/userservice/proto"
	"github.com/toncek345/userservice/service"
	"github.com/toncek345/userservice/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserServer struct {
	UserService service.UserService
	pb.UnimplementedUsersServer
}

func serviceUserToPUser(u *service.User) *pb.User {
	return &pb.User{
		Id:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Country:   u.Country,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}

func (u *UserServer) AddUser(ctx context.Context, msg *pb.AddUserMessage) (*pb.User, error) {
	// TODO: some form of validation

	user, err := u.UserService.AddUser(ctx, &service.AddUser{
		FirstName: msg.FirstName,
		LastName:  msg.LastName,
		Email:     msg.Email,
		Country:   msg.Country,
		Password:  msg.Password,
	})
	if err != nil {
		log.Printf("adding user failed: %s\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return serviceUserToPUser(user), nil
}

func (u *UserServer) DeleteUser(ctx context.Context, msg *pb.DeleteUserMessage) (*emptypb.Empty, error) {
	if err := u.UserService.DeleteUser(ctx, msg.Id); err != nil {
		log.Printf("deleting user failed: %s\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &emptypb.Empty{}, nil
}

func (u *UserServer) UpdateUser(ctx context.Context, msg *pb.UpdateUserMessage) (*pb.User, error) {
	// TODO: some form of validation

	user, err := u.UserService.UpdateUser(ctx, &service.UpdateUser{
		ID:        msg.Id,
		FirstName: msg.FirstName,
		LastName:  msg.LastName,
		Email:     msg.Email,
		Country:   msg.Country,
		Password:  msg.Password,
	})
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		log.Printf("updating user failed: %s\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return serviceUserToPUser(user), nil
}

func (u *UserServer) SearchUser(ctx context.Context, msg *pb.SearchUserMessage) (*pb.SearchUserResponse, error) {
	// TODO: extract pagination handling
	pageSize := 5
	page := 1

	if msg.Page != 0 {
		page = int(msg.Page)
	}

	if msg.PageSize != 0 {
		pageSize = int(msg.PageSize)
	}

	countryFilter := ""
	if msg.Filters != nil && msg.Filters.Country != "" {
		countryFilter = msg.Filters.Country
	}

	users, err := u.UserService.SearchUser(ctx, int64(page), int64(pageSize), countryFilter)
	if err != nil {
		log.Printf("searching users failed: %s\n", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	up := make([]*pb.User, 0, len(users))
	for _, v := range users {
		up = append(up, serviceUserToPUser(v))
	}

	return &pb.SearchUserResponse{
		Users: up,
	}, nil
}
