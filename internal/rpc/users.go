package rpc

import (
	"context"
	"log"

	pb "github.com/dom/user/api/dom/user/v1"
	"github.com/dom/user/internal/users"
)

type UserSvcParams struct {
	Query users.UserQuery
	Cmd   users.UserCommand
}

type server struct {
	pb.UnimplementedUserServiceServer
}

type UserCommands struct {
	Command users.UserCommand
}

type userSvc struct {
	query users.UserQuery
	cmd   users.UserCommand
}

func NewUserService(p *UserSvcParams) *userSvc {
	return &userSvc{
		query: p.Query,
		cmd:   p.Cmd,
	}
}

func (s *userSvc) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", req.GetName())
	return &pb.HelloReply{Message: "Hello " + req.GetName()}, nil
}

func (s *userSvc) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Printf("Received: %v", req.UserName)

	userId, err := s.cmd.SaveUser(ctx, users.User{
		UserName: req.UserName,
		LastName: req.LastName,
		Email:    req.Email,
		Role:     req.Role,
	})
	if err != nil {
		log.Printf("UserService: %w", err)
		return nil, err
	}

	return &pb.CreateUserResponse{UserId: int32(userId)}, nil
}

// DeleteUser implements v1.UserServiceServer
func (s *userSvc) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	rowsAffected, err := s.cmd.DeleteUser(ctx, int(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteUserResponse{RowsAffected: rowsAffected}, nil
}

// GetUsers implements v1.UserServiceServer
func (s *userSvc) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	users, err := s.query.GetUsers(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.GetUsersResponse{Users: users}, nil
}

// UpdateUser implements v1.UserServiceServer
func (s *userSvc) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	rowsAffected, err := s.cmd.UpdateUser(ctx, users.User{
		UserId:   int(req.UserId),
		UserName: req.UserName,
		LastName: req.LastName,
		Email:    req.Email,
		Role:     req.Role,
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpdateUserResponse{RowsAffected: rowsAffected}, nil
}
