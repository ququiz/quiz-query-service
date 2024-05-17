package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"ququiz.org/lintang/quiz-query-service/config"
)

type Mongodb struct {
	Conn *mongo.Database
}

func NewMongo(cfg *config.Config) *Mongodb {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.Mongodb.MongoURL))
	if err != nil {
		zap.L().Fatal("mongo.Connect()", zap.Error(err))
	}

	db := client.Database(cfg.Mongodb.Database)

	return &Mongodb{db}
}

func (db *Mongodb) Close(ctx context.Context) {
	db.Conn.Client().Disconnect(ctx)
}
