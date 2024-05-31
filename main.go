// Code generated by hertz generator.

package main

import (
	"fmt"
	"net"
	"time"

	"github.com/bytedance/gopkg/util/gopool"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/hertz-contrib/pprof"
	_ "go.uber.org/automaxprocs"
	"ququiz.org/lintang/quiz-query-service/biz/dal/mongodb"
	rediscache "ququiz.org/lintang/quiz-query-service/biz/dal/redisCache"
	"ququiz.org/lintang/quiz-query-service/biz/router"
	"ququiz.org/lintang/quiz-query-service/biz/service"
	"ququiz.org/lintang/quiz-query-service/config"
	"ququiz.org/lintang/quiz-query-service/pkg"

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

	h := server.Default(
		server.WithHostPorts(fmt.Sprintf(`0.0.0.0:%s`, cfg.HTTP.Port)),
		server.WithExitWaitTime(4*time.Second),
	)

	pprof.Register(h)
	var callback []route.CtxCallback

	// service & router
	mongo := mongodb.NewMongo(cfg)
	rds := rediscache.NewRedis(cfg)

	// repository
	quizRepo := mongodb.NewQuizRepository(mongo.Conn)
	questionRepo := mongodb.NewQuestionRepository(mongo.Conn)

	cacheRepo := rediscache.NewRedisCache(rds.Client)

	// service
	questionService := service.NewQuestionService(questionRepo, cacheRepo, quizRepo)
	quizService := service.NewQuizService(quizRepo)

	// router
	router.QuizRouter(h, quizService, questionService)

	// insert data to mongodb pake faker
	// util.InsertQuizData(cfg, mongo)

	callback = append(callback, mongo.Close, rds.Close)
	h.Engine.OnShutdown = append(h.Engine.OnShutdown, callback...) /// graceful shutdown

	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf(`127.0.0.1:%s`, cfg.GRPC.URLGrpc)) // grpc address
	var opts []kitexServer.Option
	opts = append(opts, kitexServer.WithMetaHandler(transmeta.ServerHTTP2Handler))
	opts = append(opts, kitexServer.WithServiceAddr(addr))
	opts = append(opts, kitexServer.WithExitWaitTime(5*time.Second))

	// srv := helloservice.NewServer(new(rpc.HelloServiceImpl), opts...) //grpc server

	// go func() {
	// 	// start kitex rpc server (grpc)
	// 	err := srv.Run()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	h.Spin()
}
