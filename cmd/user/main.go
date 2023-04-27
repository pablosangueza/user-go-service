package main

import (
	"log"
	"net"

	pb "github.com/dom/user/api/dom/user/v1"

	"github.com/dom/user/internal/config"
	"github.com/dom/user/internal/database"
	"github.com/dom/user/internal/rpc"
	"github.com/dom/user/internal/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.OpenDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, rpc.NewUserService(&rpc.UserSvcParams{
		Querier: users.NewQuerier(db),
		Cmd: &rpc.UserCommands{
			CreateUser: users.NewUserCommand(db),
		},
	}))

	reflection.Register(s)

	log.Printf("Server listening on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
