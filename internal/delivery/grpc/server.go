package grpc

import (
	"fmt"
	"golectro-product/internal/delivery/grpc/handler"
	"golectro-product/internal/delivery/grpc/interceptor"
	proto "golectro-product/internal/delivery/grpc/proto/product"
	"golectro-product/internal/usecase"
	"log"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func StartGRPCServer(productUC *usecase.ProductUseCase, port int, viper *viper.Viper) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.UnaryRequestIDInterceptor(),
			interceptor.UnaryLoggingInterceptor(productUC.Log),
		),
	)

	userHandler := &handler.ProductHandler{ProductUseCase: productUC}
	proto.RegisterProductServiceServer(grpcServer, userHandler)

	log.Printf("gRPC server listening at :%d\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
