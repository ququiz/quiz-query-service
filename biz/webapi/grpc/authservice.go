package grpc

import (
	"context"
	"ququiz/lintang/quiz-query-service/biz/domain"
	"ququiz/lintang/quiz-query-service/pb"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type AuthClient struct {
	service pb.UsersServiceClient
}

func NewAuthClient(cc *grpc.ClientConn) *AuthClient {
	svc := pb.NewUsersServiceClient(cc)
	return &AuthClient{service: svc}
}

func (a *AuthClient) GetUsersByIds(ctx context.Context, userIDs []string) ([]domain.User, error) {
	grpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.GetUserRequestByIds{
		Ids: userIDs,
	}

	res, err := a.service.GetUserByIds(grpcCtx, req)
	if err != nil {
		zap.L().Error("m.service.GetUserByIds  (GetUsersByIds) (UserGRPClient)", zap.Error(err))
		return []domain.User{}, err
	}

	var usernames []domain.User = []domain.User{}
	for i := 0; i < len(res.Users); i++ {
		usernames = append(usernames, domain.User{
			ID:       res.Users[i].Id,
			Username: res.Users[i].Username,
		})
	}

	return usernames, nil
}

func (a *AuthClient) GetUserByID(ctx context.Context, userID string) (domain.User, error) {
	grpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.GetUserRequest{
		Id: userID,
	}

	res, err := a.service.GetUserById(grpcCtx, req)
	if err != nil {
		zap.L().Error("m.service.GetUserByID  (GetUserByID) (UserGRPClient)", zap.Error(err))
		return domain.User{}, err
	}

	user := domain.User{
		ID:       userID,
		Username: res.Username,
	}

	return user, nil
}
