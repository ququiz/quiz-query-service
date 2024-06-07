package router

import (
	"context"
	"errors"
	"net/http"
	"time"

	"ququiz/lintang/quiz-query-service/biz/domain"
	"ququiz/lintang/quiz-query-service/biz/router/middleware"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuizService interface {
	GetAll(ctx context.Context) ([]domain.BaseQuiz, error)
	Get(ctx context.Context, quizID string) (domain.BaseQuiz, error)
	GetQuizByCreatorID(ctx context.Context, creatorID string) ([]domain.BaseQuiz, error)
	GetQuizHistory(ctx context.Context, participantID string) ([]domain.BaseQuizIsParticipant, error)
}

type QuestionService interface {
	GetAllByQuiz(ctx context.Context, quizID string, userID string) ([]domain.Question, error)
	GetUserAnswers(ctx context.Context, quizID string, userID string) ([]domain.QuestionUserAnswer, error)
	GetAllByQuizNotCached(ctx context.Context, quizID string, userID string) ([]domain.Question, error)
	UserAnswerAQuestion(ctx context.Context, quizID string, questionID string,
		userChoiceID string, userEssayAnswer string, userID string, username string) (bool, error)
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
			qH.GET("", handler.GetAllQuiz)

			qH.GET("/:quizID", handler.GetQuizDetail)

			qH.GET("/:quizID/questions", append(middleware.Protected(), handler.GetQuizQuestion)...)                   //append(middleware.Protected(),
			qH.GET("/:quizID/questionsNotCached", append(middleware.Protected(), handler.GetQuizQuestionNotCached)...) //append(middleware.Protected(),

			qH.GET("/:quizID/result", append(middleware.Protected(), handler.GetUserAnswer)...)
			qH.POST("/:quizID/questions/:questionID/answer", append(middleware.Protected(), handler.UserAnswerAQuestion)...)

			qH.GET("/mine", append(middleware.Protected(), handler.GetCreatedQuiz)...)
			qH.GET("/history", append(middleware.Protected(), handler.GetQuizHistory)...)

		}
	}
}

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

type listQuizResp struct {
	ID          primitive.ObjectID   `json:"id"`
	Name        string               `json:"name"`
	CreatorID   string               `json:"creator_id"`
	Passcode    string               `json:"passcode"`
	StartTime   time.Time            `json:"start_time"`
	EndTime     time.Time            `json:"end_time"`
	Status      domain.QuizStatus    `json:"status"`
	Participant []domain.Participant `json:"participants"`
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
			ID:          quiz.ID,
			Name:        quiz.Name,
			CreatorID:   quiz.CreatorID,
			Passcode:    quiz.Passcode,
			StartTime:   quiz.StartTime,
			EndTime:     quiz.EndTime,
			Status:      quiz.Status,
			Participant: quiz.Participants,
		})
	}
	c.JSON(http.StatusOK, resp)
}

type getQuestionReq struct {
	QuizID string `path:"quizID,required" vd:"regexp('^\\w') && len($) == 24;  msg:'quizID haruslah a-z,A-Z,0-9'"`
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
			ID:       q.ID.Hex(),
			Question: q.Question,
			Type:     string(q.Type),
			Choices:  q.Choices,
			Weight:   q.Weight,
		})
	}
	c.JSON(http.StatusOK, questionsRes)
}

func (h *QuizHandler) GetQuizQuestionNotCached(ctx context.Context, c *app.RequestContext) {
	var req getQuestionReq
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}
	userID, _ := c.Get("userID")
	questions, err := h.questionSvc.GetAllByQuizNotCached(ctx, req.QuizID, userID.(string))
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	var questionsRes []getQuestionRes
	for _, q := range questions {
		questionsRes = append(questionsRes, getQuestionRes{
			ID:       q.ID.Hex(),
			Question: q.Question,
			Type:     string(q.Type),
			Choices:  q.Choices,
			Weight:   q.Weight,
		})
	}
	c.JSON(http.StatusOK, questionsRes)
}

type getUserAnswerReq struct {
	QuizID string `path:"quizID,required" vd:"regexp('^\\w') && len($) == 24;  msg:'quizID haruslah a-z,A-Z,0-9'" `
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
			UserAnswer: answer.UserAnswers.Answer,
			UserChoice: answer.UserAnswers.ChoiceID,
			Weight:     answer.Weight,
			Choices:    answer.Choices,
			Type:       answer.Type,
			Question:   answer.Question,
		})
	}

	c.JSON(http.StatusOK, userAnswerRes{res})
}

type getQuizDetailReq struct {
	QuizID string `path:"quizID,required" vd:"regexp('^\\w') && len($) == 24;  msg:'quizID haruslah a-z,A-Z,0-9'" `
}

type quizRes struct {
	Quiz domain.BaseQuiz `json:"quiz"`
}

func (h *QuizHandler) GetQuizDetail(ctx context.Context, c *app.RequestContext) {
	var req getQuizDetailReq
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	// userID, err := c.Get("userID")
	quizDetail, err := h.svc.Get(ctx, req.QuizID)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}
	quizDetail.Questions = []domain.Question{}
	c.JSON(http.StatusOK, quizRes{quizDetail})
}

type userAnswerAQuestionReq struct {
	QuizID      string `path:"quizID,required" vd:"regexp('^\\w') && len($) == 24;  msg:'quizID haruslah a-z,A-Z,0-9 dan panjang haruslah 24'"`
	QuestionID  string `path:"questionID,required" vd:"regexp('^\\w') && len($) == 24;  msg:'questionID haruslah a-z,A-Z,0-9 dan panjang haruslah 24'"`
	ChoiceID    string `json:"choiceID" vd:"regexp('^\\w') && len($) == 24;  msg:'questionID haruslah a-z,A-Z,0-9 dan panjang haruslah 24'" `
	EssayAnswer string `json:"essayAnswer" vd:" ;msg:'essay answer harus anda isi gan'"`
}

type userAnswerAQuestionRes struct {
	Message   string `json:"message"`
	IsCorrect bool   `json:"isCorrect"`
}

func (h *QuizHandler) UserAnswerAQuestion(ctx context.Context, c *app.RequestContext) {
	var req userAnswerAQuestionReq
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")
	isCorrect, err := h.questionSvc.UserAnswerAQuestion(ctx, req.QuizID, req.QuestionID, req.ChoiceID, req.EssayAnswer, userID.(string), username.(string))
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}
	var resMessage string
	if isCorrect {
		resMessage = "congrats ma bro, your answer is correct "
	} else {
		resMessage = "Sorry ma bro, your answer is wrong "

	}
	c.JSON(http.StatusOK, userAnswerAQuestionRes{Message: resMessage, IsCorrect: isCorrect})
}

func (h *QuizHandler) GetCreatedQuiz(ctx context.Context, c *app.RequestContext) {
	userID, _ := c.Get("userID")
	quizs, err := h.svc.GetQuizByCreatorID(ctx, userID.(string))
	var resp []listQuizResp
	for _, quiz := range quizs {
		resp = append(resp, listQuizResp{
			ID:        quiz.ID,
			Name:      quiz.Name,
			CreatorID: quiz.CreatorID,
			Passcode:  quiz.Passcode,
			StartTime: quiz.StartTime,
			EndTime:   quiz.EndTime,
			Status:    quiz.Status,
		})
	}
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *QuizHandler) GetQuizHistory(ctx context.Context, c *app.RequestContext) {
	userID, _ := c.Get("userID")
	quizs, err := h.svc.GetQuizHistory(ctx, userID.(string))
	var resp []listQuizResp
	for _, quiz := range quizs {
		resp = append(resp, listQuizResp{
			ID:        quiz.ID,
			Name:      quiz.Name,
			CreatorID: quiz.CreatorID,
			Passcode:  quiz.Passcode,
			StartTime: quiz.StartTime,
			EndTime:   quiz.EndTime,
			Status:    domain.QuizStatus(quiz.Status),
		})
	}
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)

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
