package service

import (
	"context"
	"grpcJwt/pb"
	"time"

	"github.com/sirupsen/logrus"
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
	// exp, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXP"))
	// if err != nil {
	// 	logrus.Debugf("failed to parse token exp: %v", err)
	// }

	// logrus.Println(time.Duration(exp.Seconds()))
	logrus.Debugf(refreshTokenExp.String())

	return &AuthService{
		accountStore:    accountStore,
		jwtManager:      jwtManager,
		sessionStore:    sessionStore,
		refreshTokenExp: refreshTokenExp,
	}
}

func (auth *AuthService) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	account, err := auth.accountStore.Find(req.GetAccount().Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find account: %v", err)
	}

	if account == nil || !account.IsCorrectPassword(req.Account.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, "incorrect username/password")
	}

	accessToken, err := auth.jwtManager.GenerateAccessToken(account)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	refreshToken, err := auth.jwtManager.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	session := NewSession(refreshToken, account.Id, auth.refreshTokenExp)
	if err := auth.sessionStore.Save(session); err != nil {
		return nil, err
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
	return nil, status.Errorf(codes.Unimplemented, "Not implemented")
}

func (auth *AuthService) SignOut(ctx context.Context, req *pb.SignOutRequest) (*pb.SignOutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Not implemented method")
}

func (auth *AuthService) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Not implemented")
}
