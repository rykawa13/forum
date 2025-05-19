package controller

import (
	"context"

	"github.com/forum-backend/auth-service/internal/entity"
	"github.com/forum-backend/auth-service/internal/usecase"
	pb "github.com/forum-backend/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	authUC usecase.AuthUseCase
}

func NewAuthServer(authUC usecase.AuthUseCase) *AuthServer {
	return &AuthServer{authUC: authUC}
}

func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	input := entity.UserCreate{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	user, err := s.authUC.Register(ctx, input)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.RegisterResponse{
		Id:       int32(user.ID),
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	input := entity.UserLogin{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	tokens, err := s.authUC.Login(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &pb.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
