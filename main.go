// Code generated by hertz generator.

package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"ququiz/lintang/quiz-query-service/biz/dal/mongodb"
	"ququiz/lintang/quiz-query-service/biz/dal/rabbitmq"
	rediscache "ququiz/lintang/quiz-query-service/biz/dal/redisCache"
	"ququiz/lintang/quiz-query-service/biz/router"
	"ququiz/lintang/quiz-query-service/biz/service"
	"ququiz/lintang/quiz-query-service/biz/webapi/grpc"
	"ququiz/lintang/quiz-query-service/config"
	"ququiz/lintang/quiz-query-service/pkg"
	"ququiz/lintang/quiz-query-service/rpc"

	"github.com/bytedance/gopkg/util/gopool"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/gzip"
	grpcClient "google.golang.org/grpc"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/logger/accesslog"
	"github.com/hertz-contrib/pprof"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
	grpcGo "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	h.Use(gzip.Gzip(gzip.DefaultCompression)) // gzip compress

	h.Use(accesslog.New(accesslog.WithLogConditionFunc(func(ctx context.Context, c *app.RequestContext) bool {
		if c.FullPath() == "/healthz" {
			return false
		}
		return true
	}))) // jangan pake acess log zap (bikin latency makin tinggi)

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
	scoringSvcConsumer := rabbitmq.NewScoringSvcConsumer(rmq, cacheRepo)
	err = scoringSvcConsumer.ListenAndServe()
	if err != nil {
		zap.L().Error("ScoringSvcConsumer.ListenAndServe()", zap.Error(err))
	}

	// rabbtimq consumer producer
	quizCommandProd := rabbitmq.NewQuizCommandServiceProducerMQ(rmq)
	scoringProd := rabbitmq.NewScoringServiceProducerMQ(rmq)

	//grpc
	cc, err := grpcGo.NewClient(cfg.GRPC.AuthClient+"?wait=30s", grpcGo.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Fatal("Newclient gprc (main)", zap.Error(err))
	}
	defer cc.Close() // close auth grpc connection when graceful shutdown to avoid memory leaks

	authClient := grpc.NewAuthClient(cc)

	// service
	questionService := service.NewQuestionService(questionRepo, cacheRepo, quizRepo, scoringProd, quizCommandProd)
	quizService := service.NewQuizService(quizRepo, authClient)

	// router
	h.GET("/healthz", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, "service is healthy")
	}) // health probes
	router.QuizRouter(h, quizService, questionService)

	callback = append(callback, mongo.Close, rds.Close, rmq.Close)
	h.Engine.OnShutdown = append(h.Engine.OnShutdown, callback...) /// graceful shutdown

	// addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf(`0.0.0.0:%s`, cfg.GRPC.URLGrpc)) // grpc address
	// var opts []kitexServer.Option
	// opts = append(opts, kitexServer.WithMetaHandler(transmeta.ServerHTTP2Handler))
	// opts = append(opts, kitexServer.WithServiceAddr(addr))
	// opts = append(opts, kitexServer.WithExitWaitTime(5*time.Second))
	// opts = append(opts, kitexServer.WithGRPCWriteBufferSize(1000*1000*100))
	// opts = append(opts, kitexServer.WithGRPCReadBufferSize(1000*1000*100))
	// opts = append(opts, kitexServer.WithGRPCInitialConnWindowSize(1000*1000*100))
	// opts = append(opts, kitexServer.WithGRPCInitialWindowSize(1000*1000*100))

	// quizGRPCSvc := rpc.NewQuizService(questionRepo, quizRepo)
	// srv := quizqueryservice.NewServer(quizGRPCSvc, opts...)
	// klog.SetLogger(kitexlogrus.NewLogger())
	// klog.SetLevel(klog.LevelDebug)
	// go func() {
	// 	// start kitex rpc server (grpc)

	// 	err := srv.Run()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	listener, err := net.Listen("tcp", cfg.GRPC.URLGrpc)
	if err != nil {
		zap.L().Fatal("net.Listen(, cfg.GRPC.URLGrpc)", zap.Error(err))
	}

	quizServiceGRPC := rpc.NewQuizService(questionRepo, quizRepo)
	grpcServerChan := make(chan *grpcClient.Server)

	go func() {
		zap.L().Info("grpc server run on port: " + cfg.GRPC.URLGrpc)
		err := service.RunGRPCServer(quizServiceGRPC, listener, grpcServerChan)
		if err != nil {
			zap.L().Fatal("cannot start GRPC  Server", zap.Error(err))
		}
	}()
	var grpcServer = <-grpcServerChan
	fmt.Println(grpcServer)
	h.Spin()
}
