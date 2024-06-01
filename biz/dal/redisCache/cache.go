package rediscache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"ququiz/lintang/quiz-query-service/biz/domain"
)

type RedisCache struct {
	cli *redis.Client
}

func NewRedisCache(cli *redis.Client) *RedisCache {
	return &RedisCache{cli}
}

func (c *RedisCache) GetCachedQuestion(ctx context.Context, quizID string) ([]domain.Question, error) {
	var questions []domain.Question
	var jsonQuestion []byte
	err := c.cli.Get(ctx, fmt.Sprintf("questions:%s", quizID)).Scan(&jsonQuestion)
	if err != nil {
		zap.L().Debug(fmt.Sprintf("Question from Quiz %s not cached", quizID))
		return []domain.Question{}, err
	}

	err = json.Unmarshal(jsonQuestion, &questions)
	if err != nil {
		zap.L().Error("json.Unmarshal (GetCachaedQuestion) (Cache)", zap.Error(err))
		return []domain.Question{}, domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	return questions, nil
}

func (c *RedisCache) SetCachedQuestion(ctx context.Context, quizID string, qs []domain.Question) error {

	// err := c.cli.HSet(ctx, fmt.Sprintf("questions:%s", quizID), qs)
	jsonQuestion, err := json.Marshal(qs)
	if err != nil {
		zap.L().Error("json.Marshal (SetCachedQuestion) (RedisCache)", zap.Error(err))
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}
	err = c.cli.Set(ctx, fmt.Sprintf("questions:%s", quizID), jsonQuestion, 1*time.Hour).Err()
	if err != nil {
		zap.L().Error("c.cli.Hset (SetCachedQuestion) (RedisCache)", zap.Error(err))
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	// set expiration for 1h
	c.cli.Expire(ctx, fmt.Sprintf("questions:%s", quizID), 1*time.Hour)
	return nil
}
