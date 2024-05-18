package rpc

import (
	"context"
	pb "ququiz.org/lintang/quiz-query-service/kitex_gen/go_hertz_template_lintang/pb"
)

// HelloServiceImpl implements the last service interface defined in the IDL.
type HelloServiceImpl struct{}

// Hello implements the HelloServiceImpl interface.
func (s *HelloServiceImpl) Hello(ctx context.Context, req *pb.HelloReq) (resp *pb.HelloResp, err error) {
	// TODO: Your code here...
	return
}
