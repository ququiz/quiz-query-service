package main

import (
	"context"
	pb "ququiz.org/lintang/quiz-query-service/kitex_gen/quiz-query-service/pb"
)

// QuizQueryServiceImpl implements the last service interface defined in the IDL.
type QuizQueryServiceImpl struct{}

// GetQuestionDetail implements the QuizQueryServiceImpl interface.
func (s *QuizQueryServiceImpl) GetQuestionDetail(ctx context.Context, req *pb.GetQuestionReq) (resp *pb.GetQuestionRes, err error) {
	// TODO: Your code here...
	return
}

// GetQuizParticipants implements the QuizQueryServiceImpl interface.
func (s *QuizQueryServiceImpl) GetQuizParticipants(ctx context.Context, req *pb.GetQuizParticipantsReq) (resp *pb.GetQuizParticipantRes, err error) {
	// TODO: Your code here...
	return
}
