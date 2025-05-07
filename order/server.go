package order

import (
	"fmt"
	interceptor "go-courier/interceptor"
	orderpb "go-courier/proto/order"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func StartGRPCServer(port int, OrderService orderpb.OrderServiceServer) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthInterceptor(os.Getenv("JWT_SECRET_KEY"))),
	)
	orderpb.RegisterOrderServiceServer(grpcServer, OrderService)
	log.Printf("Order gRPC server running on port %d...", port)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
