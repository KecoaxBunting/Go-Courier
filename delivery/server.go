package delivery

import (
	"fmt"
	interceptor "go-courier/interceptor"
	deliverypb "go-courier/proto/delivery"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func StartGRPCServer(port int, DeliveryService deliverypb.DeliveryServiceServer) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthInterceptor(os.Getenv("JWT_SECRET_KEY"))),
	)
	deliverypb.RegisterDeliveryServiceServer(grpcServer, DeliveryService)
	log.Printf("Delivery gRPC server running on port %d...", port)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
