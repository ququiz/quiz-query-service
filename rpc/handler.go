package rpc

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
	questionRepo service.QuestionRepository
	quizRepo     service.QuizRepository
}

func NewQuizService(qs service.QuestionRepository, quiz service.QuizRepository) *QuizQueryServiceImpl {
	return &QuizQueryServiceImpl{
		questionRepo: qs,
		quizRepo:     quiz,
	}
}

// GetQuestionDetail implements the QuizQueryServiceImpl interface.
func (s *QuizQueryServiceImpl) GetQuestionDetail(ctx context.Context, req *pb.GetQuestionReq) (resp *pb.GetQuestionRes, err error) {
	// TODO: Your code here...
	quiz, err := s.questionRepo.GetQuestionByIDAndQuizID(ctx, req.QuizId, req.QuestionId)
	if err != nil {
		zap.L().Error("s.questionRepo.GetQuestionByIDAndQuizID (Getquiz) (QuizQueryGrpcService)", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "question with id %s not found in quiz %s", req.QuestionId, req.QuizId)
	}

	var correctChoiceID string
	for i := 0; i < len(quiz.Questions.Choices); i++ {
		if quiz.Questions.Choices[i].IsCorrect {
			correctChoiceID = quiz.Questions.Choices[i].ID.Hex()
		}
	}
	resp.CorrectChoiceId = correctChoiceID
	res := &pb.GetQuestionRes{
		Weight:               uint64(quiz.Questions.Weight),
		CorrectEssayAnswerId: quiz.Questions.CorrectAnswer,
		CorrectChoiceId:      correctChoiceID,
	}

	return res, nil
}

// GetQuizParticipants implements the QuizQueryServiceImpl interface.
func (s *QuizQueryServiceImpl) GetQuizParticipants(ctx context.Context, req *pb.GetQuizParticipantsReq) (resp *pb.GetQuizParticipantRes, err error) {
	// TODO: Your code here...
	quiz, err := s.quizRepo.Get(ctx, req.QuizId)
	if err != nil {
		zap.L().Error("s.quizRepo.Get (GetQuizParticipants) (QuizQueryGrpcService)", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "quiz with id %s not found", req.QuizId)
	}
	var participantUserIDs []string = []string{}
	for i := 0; i < len(quiz.Participants); i++ {
		participantUserIDs = append(participantUserIDs, quiz.Participants[i].UserID)
	}
	res := &pb.GetQuizParticipantRes{
		UserIds:  participantUserIDs,
		QuizName: quiz.Name,
	}

	return res, nil
}
