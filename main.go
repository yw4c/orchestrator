package main

import (
	"fmt"
	"net"
	"orchestrator/config"
	"orchestrator/facade"
	"orchestrator/handler"
	"orchestrator/pb"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("container initializing ... ")

	// register facades
	facade.RegisterAsyncBookingFacade()
	facade.RegisterSyncBookingFacade()

	// gRPC Connection
	lis, err := net.Listen("tcp", ":"+config.GetConfigInstance().Server.Grpc.Port)
	if err != nil {
		log.Panic().Msg(err.Error())
	}

	// recover gRPC panic
	var customFunc grpc_recovery.RecoveryHandlerFunc
	customFunc = func(p interface{}) (err error) {
		log.Error().Interface("message", p).
			Str("trace", string(debug.Stack())).
			Msg("GRPC Recover")

		return status.Errorf(codes.Unknown, "panic recovered: %v", p)
	}
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}

	// gRPC middleware
	grpcSvc := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(opts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(opts...),
		),
	)

	// Bind gRPC endpoints under here
	pb.RegisterBookingServiceServer(grpcSvc, &handler.BookingService{})

	reflection.Register(grpcSvc)

	go func() {
		if err := grpcSvc.Serve(lis); err != nil {
			log.Panic().Msg(err.Error())
		}
	}()

	log.Info().Msg("gRPC Server Started ")
	fmt.Println("container is ready ")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	log.Info().Msg("Shutting down server...")
}
