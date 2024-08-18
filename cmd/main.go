package main

import (
	"context"
	"fmt"
	"grpcJwt/initializers"
	"grpcJwt/internal/db"
	"grpcJwt/internal/service"
	"grpcJwt/pb"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	err := initializers.LoadEnvVariables()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config file: %v", err))
	}

	initializers.SetupLogger()

	accessTokenExpDur, err := time.ParseDuration(os.Getenv("TOKEN_EXP"))
	if err != nil {
		panic(fmt.Sprintf("Failed to get duration access token: %v", err))
	}

	jwtManager := service.NewJWTManager(os.Getenv("SECRET_KEY"), accessTokenExpDur)
	accountStore := db.NewInMemoryAccountStore()
	sessionStore := db.NewInMemorySessionStore()

	// if err := seedAccounts(accountStore); err != nil {
	// 	panic("cannot seedUsers")
	// }

	refreshTokenExpDur, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXP"))
	if err != nil {
		panic(fmt.Sprintf("Failed to get duration refresh token: %v", err))
	}
	authService := service.NewAuthService(accountStore, jwtManager, sessionStore, refreshTokenExpDur)

	server := grpc.NewServer(grpc.UnaryInterceptor(unary()))

	pb.RegisterAuthServiceServer(server, authService)

	reflection.Register(server)

	serverUrl := fmt.Sprintf("%s:%s", os.Getenv("SERVER_ADDRESS"), os.Getenv("SERVER_PORT"))
	listener, err := net.Listen("tcp", serverUrl)
	if err != nil {
		panic(fmt.Sprintf("Failed to start listen on %s", serverUrl))
	}

	go func() {
		logrus.Infof("Starting server on %s", serverUrl)
		err = server.Serve(listener)
		if err != nil {
			logrus.Errorf("Failed to serve: %v", err)
		}
	}()

	<-quit
	logrus.Info("Received termination signal, stopping server...")
	server.GracefulStop()
	logrus.Info("Server stopped, exiting")
}

// func seedAccounts(accountStore service.AccountStore) error {
// 	user1, err := service.NewAccount("user@gmail.com", "pwd", "user", 15)
// 	if err != nil {
// 		return err
// 	}

// 	return accountStore.Save(user1)
// }

func unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// logrus.Println("--> unary interceptor: ", info.FullMethod)

		// Получение имя метода
		// strings := strings.Split(info.FullMethod, "/")
		// methodName := strings[len(strings)-1]
		// logrus.Println(methodName)

		return handler(ctx, req)
	}
}
