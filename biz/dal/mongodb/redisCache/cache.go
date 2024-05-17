package rediscache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"ququiz.org/lintang/quiz-query-service/biz/domain"
)

type RedisCache struct {
	cli *redis.Client
}


func NewRedisCache(cli *redis.Client) *RedisCache {
	return &RedisCache{cli}
}

func (c *RedisCache) GetCachedQuestion(ctx context.Context, quizID string) ([]domain.Question, error) {
	var questions []domain.Question
	err := c.cli.HGetAll(ctx, fmt.Sprintf("questions:%s", quizID)).Scan(&questions)

	if err != nil {
		zap.L().Debug(fmt.Sprintf("Question from Quiz %s not cached", quizID))
		return []domain.Question{}, err
	}

	return questions, nil
}

func (c *RedisCache) SetCachedQuestion(ctx context.Context, quizID string, qs []domain.Question) error {
	
	err := c.cli.HSet(ctx, fmt.Sprintf("questions:%s", quizID), qs)
	if err != nil {
		zap.L().Error("c.cli.Hset (SetCachedQuestion) (RedisCache)", zap.Error(err.Err()))

	}

	// set expiration for 1h 
	c.cli.Expire(ctx, fmt.Sprintf("questions:%s", quizID), 1 * time.Hour )
	return nil
}
