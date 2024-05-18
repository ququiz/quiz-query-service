package router

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"ququiz.org/lintang/quiz-query-service/biz/domain"
	"ququiz.org/lintang/quiz-query-service/biz/router/middleware"
)

type QuizService interface {
	GetAll(ctx context.Context) ([]domain.BaseQuiz, error)
}

type QuestionService interface {
	GetAllByQuiz(ctx context.Context, quizID string, userID string) ([]domain.Question, error)
	GetUserAnswers(ctx context.Context, quizID string, userID string) ([]domain.QuestionWithUserAnswerAggregate, error)
}

type QuizHandler struct {
	svc         QuizService
	questionSvc QuestionService
}

func QuizRouter(r *server.Hertz, q QuizService, questionSvc QuestionService) {
	handler := &QuizHandler{
		svc:         q,
		questionSvc: questionSvc,
	}

	root := r.Group("/api/v1")
	{
		qH := root.Group("/quiz")
		{
			qH.GET("/", handler.GetAllQuiz)
			qH.GET("/quiz/questions", append(middleware.Protected(), handler.GetQuizQuestion)...)
		}
	}
}

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

type listQuizResp struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	CreatorID string            `json:"creator_id"`
	Passcode  string            `json:"passcode"`
	StartTime time.Time         `json:"start_time"`
	EndTime   time.Time         `json:"end_time"`
	Status    domain.QuizStatus `json:"status"`
}

func (h *QuizHandler) GetAllQuiz(ctx context.Context, c *app.RequestContext) {

	quizs, err := h.svc.GetAll(ctx)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	var resp []listQuizResp
	for _, quiz := range quizs {
		resp = append(resp, listQuizResp{
			ID:        quiz.ID.String(),
			Name:      quiz.Name,
			CreatorID: quiz.CreatorID,
			Passcode:  quiz.Passcode,
			StartTime: quiz.StartTime,
			EndTime:   quiz.EndTime,
			Status:    quiz.Status,
		})
	}
	c.JSON(http.StatusOK, resp)
}

type getQuestionReq struct {
	QuizID string `query:"quizID,required" vd:"regexp('\\w);  msg:'quizID haruslah a-z,A-Z,0-9'"`
}

type getQuestionRes struct {
	ID       string          `json:"id"`
	Question string          `json:"question"`
	Type     string          `json:"type"`
	Choices  []domain.Choice `json:"choices"`
	Weight   int32           `json:"weight"`
}

func (h *QuizHandler) GetQuizQuestion(ctx context.Context, c *app.RequestContext) {
	var req getQuestionReq
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}
	userID, _ := c.Get("userID")
	questions, err := h.questionSvc.GetAllByQuiz(ctx, req.QuizID, userID.(string))
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	var questionsRes []getQuestionRes
	for _, q := range questions {
		questionsRes = append(questionsRes, getQuestionRes{
			ID:       q.ID.String(),
			Question: q.Question,
			Type:     string(q.Type),
			Choices:  q.Choices,
			Weight:   q.Weight,
		})
	}
	c.JSON(http.StatusOK, questionsRes)
}

type getUserAnswerReq struct {
	QuizID string `path:"quizID,required" vd:"regexp('\\w);  msg:'quizID haruslah a-z,A-Z,0-9'" `
}

type userAnswerRes struct {
	UserAnswers []domain.QuestionAndUserAnswer `json:"user_answers"`
}

func (h *QuizHandler) GetUserAnswer(ctx context.Context, c *app.RequestContext) {

	var req getUserAnswerReq
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	userAnswers, err := h.questionSvc.GetUserAnswers(ctx, req.QuizID, userID.(string))
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	var res []domain.QuestionAndUserAnswer
	for _, answer := range userAnswers {
		res = append(res, domain.QuestionAndUserAnswer{
			UserAnswer: answer.UserAnswer.Answer,
			UserChoice: answer.UserAnswer.ChoiceID,
			Weight:     answer.Weight,
			Choices:    answer.Choices,
			Type:       answer.Type,
			Question:   answer.Question,
		})
	}

	c.JSON(http.StatusOK, userAnswerRes{res})
}



func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	var ierr *domain.Error
	if !errors.As(err, &ierr) {
		return http.StatusInternalServerError
	} else {
		switch ierr.Code() {
		case domain.ErrInternalServerError:
			return http.StatusInternalServerError
		case domain.ErrNotFound:
			return http.StatusNotFound
		case domain.ErrConflict:
			return http.StatusConflict
		case domain.ErrBadParamInput:
			return http.StatusBadRequest
		default:
			return http.StatusInternalServerError
		}
	}

}
