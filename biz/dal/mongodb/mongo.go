package mongodb

import (
	"context"
	"fmt"
	"time"

	"ququiz/lintang/quiz-query-service/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Mongodb struct {
	Conn      *mongo.Database
	FakerConn *mongo.Database
}

func NewMongo(cfg *config.Config) *Mongodb {

	zap.L().Info(fmt.Sprintf("url mongo: %s", cfg.Mongodb.MongoURL))
	ctxTimeoutRead, cancelRead := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelRead()
	client, err := mongo.Connect(ctxTimeoutRead, options.Client().ApplyURI(cfg.Mongodb.MongoURL))
	if err != nil {
		zap.L().Fatal("mongo.Connect()", zap.Error(err))
	}

	db := client.Database(cfg.Mongodb.Database)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	writeClient, err := mongo.Connect(ctxTimeout, options.Client().ApplyURI(cfg.Mongodb.MongoWriteURL))
	if err != nil {
		zap.L().Fatal("mongo.Connect() (write db)", zap.Error(err))
	}
	writeDb := writeClient.Database(cfg.Mongodb.Database)

	return &Mongodb{db, writeDb}
}

func (db *Mongodb) Close(ctx context.Context) {
	db.Conn.Client().Disconnect(ctx)
}


