package main

import (
	"context"
	"ququiz/lintang/quiz-query-service/biz/service"
	pb "ququiz/lintang/quiz-query-service/kitex_gen/quiz-query-service/pb"

	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2/codes"
	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2/status"
	"go.uber.org/zap"
)

// QuizQueryServiceImpl implements the last service interface defined in the IDL.
type QuizQueryServiceImpl struct {
	quizRepo service.QuizRepository
}

// GetQuestionDetail implements the QuizQueryServiceImpl interface.
func (s *QuizQueryServiceImpl) GetQuestionDetail(ctx context.Context, req *pb.GetQuestionReq) (resp *pb.GetQuestionRes, err error) {
	// TODO: Your code here...
	return
}

// GetQuizParticipants implements the QuizQueryServiceImpl interface.
func (s *QuizQueryServiceImpl) GetQuizParticipants(ctx context.Context, req *pb.GetQuizParticipantsReq) (resp *pb.GetQuizParticipantRes, err error) {
	// TODO: Your code here...
	q, err := s.quizRepo.Get(ctx, req.QuizId)
	if err != nil {
		zap.L().Error("s.quizRepo.Get (GetQuizParticipants) (QuizQueryServiceImpl)", zap.Error(err))
		return &pb.GetQuizParticipantRes{}, status.Errorf(codes.Internal, "s.quizRepo.Get  (GetQuizParticipants) %w", err)
	}

	participants := q.Participants
	var userIDs []string = []string{}
	for i := 0; i < len(participants); i++ {
		userIDs = append(userIDs, participants[i].UserID)
	}
	resp.UserIds = userIDs
	return resp, nil
}
