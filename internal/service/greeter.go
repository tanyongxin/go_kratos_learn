package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"helloworld/api/helloworld/v1"
	"helloworld/internal/biz"
)

// GreeterService is a greeter service.。相当于 java 中的 service 层
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	//if err != nil {
	//	return nil, err
	//}
	if in.Name == "a" {
		//panic(errors.New(403,"invalid param", "参数有误"))
		e := errors.New(403, "invalid param", "参数有误")
		e.Metadata = make(map[string]string, 1)
		e.Metadata["Name"] = in.Name

		return nil, e
	}

	return &v1.HelloReply{Message: "Hello CNM" /*+ g.Hello*/}, nil
}

//func (s *GreeterService) Test(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
//	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
//	if err != nil {
//		return nil, err
//	}
//	return &v1.HelloReply{Message: "Test CNM" + g.Hello}, nil
//}
