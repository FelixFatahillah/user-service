package cmd

import (
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"user-service/internal/config"
	"user-service/internal/delivery"
	"user-service/pkg/logger"
)

var (
	addr = flag.String("addr", fmt.Sprintf(":%d", config.Viper().GetInt("PORT_SERVER")), "the address to connect to")
	port = flag.Int("port", config.Viper().GetInt("PORT_GRPC"), "gRPC server port")
)

func Execute() {
	flag.Parse()

	numCPU := runtime.NumCPU()
	logger.Info(fmt.Sprintf("Number of CPU cores: %d", numCPU), zap.Field{
		Key:    "context",
		String: "server",
		Type:   zapcore.StringType,
	})

	// Set context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Open Redis connection
	redisClient := config.NewRedis()
	logger.Info(redisClient.Ping(ctx).String(), zap.Field{
		Key:    "context",
		String: "redis",
		Type:   zapcore.StringType,
	})

	// Open DB connection
	db, err := config.NewClient()
	if err != nil {
		logger.Error("error", zap.Error(err))
	}

	// Repositories
	repository := delivery.NewRepository(db)
	if err != nil {
		logger.Error("error", zap.Error(err))
	}

	// RPC connection
	//catalogConn, err := config.NewGrpcDial(ctx, config.Viper().GetString("GRPC_CATALOG_ADDR"))
	//if err != nil {
	//	logger.Fatal("Failed to connect to catalog server", zap.Error(err))
	//}

	//paymentConn, err := config.NewGrpcDial(ctx, config.Viper().GetString("GRPC_PAYMENT_ADDR"))
	//if err != nil {
	//	logger.Fatal("Failed to connect to catalog server", zap.Error(err))
	//}
	//
	//vocaCreditConn, err := config.NewGrpcDial(ctx, config.Viper().GetString("GRPC_VOCA_CREDIT_ADDR"))
	//if err != nil {
	//	logger.Fatal("Failed to connect to catalog server", zap.Error(err))
	//}
	//
	//rpcConnection := delivery.NewGRPC(
	//	catalogConn,
	//	paymentConn,
	//	vocaCreditConn,
	//)

	// Services
	services := delivery.NewService(delivery.Deps{
		Repository: repository,
		Redis:      *redisClient,
	})

	// Register GRPC
	rpcServer := grpc.NewServer()
	tcpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	go func() {
		if err := rpcServer.Serve(tcpListener); err != nil {
			log.Fatalf("failed to serve : %v", err)
		}
	}()
	logger.Info("Grpc Server running ...")

	// Consumer

	handler := delivery.NewHandler(
		services.UserService,
	)

	app := handler.Init()

	go func() {
		if err := app.Listen(*addr); err != nil {
			logger.Error("server error", zap.Error(err))
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("Shutting down the server...")

	cancel()

	logger.Info("Server has been shut down gracefully")
}
