package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	pb "github.com/toncek345/userservice/proto"
	"github.com/toncek345/userservice/server/health"
	"github.com/toncek345/userservice/server/users"
	"github.com/toncek345/userservice/service"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	server       *grpc.Server
	grpcListener net.Listener
	httpPort     int
	mux          http.Handler
	closeProxy   context.CancelFunc
}

func (s *Server) Start() error {
	return s.server.Serve(s.grpcListener)
}

func (s *Server) StartHTTP() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", s.httpPort), s.mux)
}

func (s *Server) Stop() {
	s.closeProxy()
	s.server.GracefulStop()
}

func NewServer(grpcPort, httpPort int, userService service.UserService) (*Server, error) {
	grpcHost := fmt.Sprintf("localhost:%d", grpcPort)
	lis, err := net.Listen("tcp", grpcHost)
	if err != nil {
		return nil, fmt.Errorf("net listen: %w", err)
	}

	server := grpc.NewServer()
	pb.RegisterUsersServer(server, &users.UserServer{UserService: userService})
	pb.RegisterHealthServer(server, &health.HealthServer{})

	ctx, cancel := context.WithCancel(context.Background())
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterUsersHandlerFromEndpoint(ctx, mux, grpcHost, opts); err != nil {
		defer cancel()
		return nil, fmt.Errorf("register user service: %w", err)
	}
	if err := pb.RegisterHealthHandlerFromEndpoint(ctx, mux, grpcHost, opts); err != nil {
		defer cancel()
		return nil, fmt.Errorf("register user service: %w", err)
	}

	return &Server{server, lis, httpPort, mux, cancel}, nil
}
