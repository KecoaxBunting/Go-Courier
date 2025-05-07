package auth

import (
	"fmt"
	interceptor "go-courier/interceptor"
	authpb "go-courier/proto/auth"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func StartGRPCServer(port int, AuthService authpb.AuthServiceServer) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthInterceptor(os.Getenv("JWT_SECRET_KEY"))),
	)
	authpb.RegisterAuthServiceServer(grpcServer, AuthService)
	log.Printf("Auth gRPC server running on port %d...", port)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
