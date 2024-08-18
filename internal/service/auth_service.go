package service

import (
	"context"
	"grpcJwt/pb"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer

	accountStore    AccountStore
	jwtManager      *JWTManager
	sessionStore    SessionStore
	refreshTokenExp time.Duration
}

func NewAuthService(accountStore AccountStore, jwtManager *JWTManager, sessionStore SessionStore, refreshTokenExp time.Duration) *AuthService {
	return &AuthService{
		accountStore:    accountStore,
		jwtManager:      jwtManager,
		sessionStore:    sessionStore,
		refreshTokenExp: refreshTokenExp,
	}
}

func (auth *AuthService) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	if req.Account.Email == "" || req.Account.Password == "" {
		return nil, status.Errorf(codes.FailedPrecondition, "parametrs email/password is empty")
	}

	account, err := auth.accountStore.Find(req.Account.Email)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot find account with email %s", req.Account.Email)
	}

	if account == nil || !account.IsCorrectPassword(req.Account.Password) {
		return nil, status.Errorf(codes.NotFound, "incorrect username/password")
	}

	accessToken, err := auth.jwtManager.GenerateAccessToken(account)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	refreshToken, err := auth.jwtManager.GenerateRefreshToken()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate refresh token")
	}

	session := NewSession(refreshToken, account.Id, auth.refreshTokenExp)
	if err := auth.sessionStore.Save(session); err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create a session")
	}

	res := &pb.SignInResponse{
		PairTokens: &pb.PairTokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}

	return res, nil
}

func (auth *AuthService) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	if req.Account.Email == "" || req.Account.Password == "" || req.Account.Age == 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "some parametrs is empty")
	}

	newAccount, err := NewAccount(req.Account.Email, req.Account.Password, "user", req.Account.Age)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create account: %v", err)
	}

	if err := auth.accountStore.Save(newAccount); err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "Account with email %s already exist", newAccount.Email)
	}

	accessToken, err := auth.jwtManager.GenerateAccessToken(newAccount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	refreshToken, err := auth.jwtManager.GenerateRefreshToken()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate refresh token")
	}

	session := NewSession(refreshToken, newAccount.Id, auth.refreshTokenExp)
	if err := auth.sessionStore.Save(session); err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create a session")
	}

	res := &pb.SignUpResponse{
		PairTokens: &pb.PairTokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}

	return res, nil
}

func (auth *AuthService) SignOut(ctx context.Context, req *pb.SignOutRequest) (*pb.SignOutResponse, error) {

	return nil, status.Errorf(codes.Unimplemented, "Not implemented method")
}

func (auth *AuthService) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Not implemented")
}
