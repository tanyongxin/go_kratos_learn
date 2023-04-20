package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	v1 "helloworld/api/helloworld/v1"
	"helloworld/internal/conf"
	"helloworld/internal/service"
	"runtime"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *http.Server {

	myHandlerFunc := func(ctx context.Context, req, err interface{}) error {

		e, ok := err.(*errors.Error)

		if ok {
			//fmt.Println(e.Code)
			//fmt.Println(e.Unwrap())

			buf := make([]byte, 64<<10) //nolint:gomnd
			// 调用栈的信息都存放在 buf 中
			n := runtime.Stack(buf, false)
			log.Context(ctx).Infof("堆栈信息大小：%s\n", n)
			buf = buf[:n]
			log.Context(ctx).Errorf("%v: %+v\n%s\n", e, req, buf)
			return errors.Clone(e)
		}

		return recovery.ErrUnknownRequest

	}

	option := recovery.WithHandler(myHandlerFunc)

	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(option),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
