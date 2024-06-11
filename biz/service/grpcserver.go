package service

import (
	"net"
	"ququiz/lintang/quiz-query-service/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunGRPCServer(
	quizGRPCServer pb.QuizQueryServiceServer,
	listener net.Listener,
	ch chan *grpc.Server,
) error {
	// GRPC Server
	grpcServer := grpc.NewServer()
	pb.RegisterQuizQueryServiceServer(grpcServer, quizGRPCServer)
	reflection.Register(grpcServer)

	ch <- grpcServer

	return grpcServer.Serve(listener)
}
