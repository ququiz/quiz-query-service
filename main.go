// Code generated by hertz generator.

package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"ququiz/lintang/quiz-query-service/biz/dal/mongodb"
	"ququiz/lintang/quiz-query-service/biz/dal/rabbitmq"
	rediscache "ququiz/lintang/quiz-query-service/biz/dal/redisCache"
	"ququiz/lintang/quiz-query-service/biz/router"
	"ququiz/lintang/quiz-query-service/biz/service"
	"ququiz/lintang/quiz-query-service/biz/webapi/grpc"
	"ququiz/lintang/quiz-query-service/config"
	"ququiz/lintang/quiz-query-service/kitex_gen/quiz-query-service/pb/quizqueryservice"
	"ququiz/lintang/quiz-query-service/pkg"
	"ququiz/lintang/quiz-query-service/rpc"

	"github.com/bytedance/gopkg/util/gopool"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/logger/accesslog"
	"github.com/hertz-contrib/pprof"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
	grpcGo "google.golang.org/grpc"

	kitexServer "github.com/cloudwego/kitex/server"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		hlog.Fatalf("Config error: %s", err)
	}
	gopool.SetCap(500000) // naikin goroutine netpool for high performance

	logsCores := pkg.InitZapLogger(cfg)
	defer logsCores.Sync()
	hlog.SetLogger(logsCores)

	zap.L().Debug(fmt.Sprint(`0.0.0.0:%s`, cfg.HTTP.Port))
	h := server.Default(
		server.WithHostPorts(fmt.Sprintf(`0.0.0.0:%s`, cfg.HTTP.Port)),
		server.WithExitWaitTime(4*time.Second),
	)

	h.Use(accesslog.New()) // jangan pake acess log zap (bikin latency makin tinggi)

	// setup cors
	corsHandler := cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "content-type", "authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"},
		ExposeHeaders:    []string{"Origin", "content-type", "authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"},
		AllowCredentials: true,

		MaxAge: 12 * time.Hour,
	})

	h.Use(corsHandler)

	pprof.Register(h)
	var callback []route.CtxCallback

	// service & router
	mongo := mongodb.NewMongo(cfg)
	rds := rediscache.NewRedis(cfg)

	// repository
	quizRepo := mongodb.NewQuizRepository(mongo.Conn)
	questionRepo := mongodb.NewQuestionRepository(mongo.Conn)

	cacheRepo := rediscache.NewRedisCache(rds.Client)

	// rabbitmq
	rmq := rabbitmq.NewRabbitMQ(cfg)
	scoringSvcConsumer := rabbitmq.NewScoringSvcConsumer(rmq)
	err = scoringSvcConsumer.ListenAndServe()

	// rabbtimq consumer producer
	quizCommandProd := rabbitmq.NewQuizCommandServiceProducerMQ(rmq)
	scoringProd := rabbitmq.NewScoringServiceProducerMQ(rmq)

	//grpc
	cc, err := grpcGo.NewClient(cfg.GRPC.AuthClient)
	if err != nil {
		zap.L().Fatal("Newclient gprc (main)", zap.Error(err))
	}

	authClient := grpc.NewAuthClient(cc)

	// service
	questionService := service.NewQuestionService(questionRepo, cacheRepo, quizRepo, scoringProd, quizCommandProd)
	quizService := service.NewQuizService(quizRepo, authClient)

	// router
	router.QuizRouter(h, quizService, questionService)

	// insert data to mongodb pake faker

	callback = append(callback, mongo.Close, rds.Close)
	h.Engine.OnShutdown = append(h.Engine.OnShutdown, callback...) /// graceful shutdown

	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf(`127.0.0.1:%s`, cfg.GRPC.URLGrpc)) // grpc address
	var opts []kitexServer.Option
	opts = append(opts, kitexServer.WithMetaHandler(transmeta.ServerHTTP2Handler))
	opts = append(opts, kitexServer.WithServiceAddr(addr))
	opts = append(opts, kitexServer.WithExitWaitTime(5*time.Second))

	quizGRPCSvc := rpc.NewQuizService(questionRepo, quizRepo)
	srv := quizqueryservice.NewServer(quizGRPCSvc, opts...)
	go func() {
		// start kitex rpc server (grpc)
		err := srv.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()

	h.Spin()
}
