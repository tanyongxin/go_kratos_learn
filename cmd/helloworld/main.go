package main

import (
	_ "context"
	"flag"
	"fmt"
	_ "fmt"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	_ "github.com/go-kratos/kratos/v2/errors"
	_ "github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"helloworld/internal/biz"
	"helloworld/internal/data"
	"helloworld/internal/server"
	"helloworld/internal/service"
	"os"

	_ "helloworld/api/helloworld/v1"

	"helloworld/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	_ "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	_ "github.com/go-sql-driver/mysql"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func newApp1(logger log.Logger, gs *grpc.Server, hs *http.Server, r *nacos.Registry) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
		kratos.Registrar(r),
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	// 读取配置
	if err := c.Load(); err != nil {
		panic(err)
	}

	// conf.Bootstrap 是一个结构体
	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)

	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	// app.Run() 方法中会调用注册中心的注册服务的方法，服务名为 si.Name + "." + u.Scheme
	// si.Name 就是 Name 属性值，u.Scheme 就是服务暴露的方式，本例中有两种，一种是 http，另一种是 grpc
	// 因此在 nacos 控制台上看到的服务名称就是 helloWorld-server.grpc 和 helloWorld-server.http
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func wireApp1(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	// greeterRepo 类似于 java 中的 Repo 接口。该接口提供了操作 internal/biz/greeter.go:18 实体类的增删查改方法
	greeterRepo := data.NewGreeterRepo(dataData, logger)
	// greeterUsecase 类似于 java 中的 Manager 层，Manager 层中注入了 greeterRepo
	greeterUsecase := biz.NewGreeterUsecase(greeterRepo, logger)
	// greeterService 类似于 java 中的 service 层
	greeterService := service.NewGreeterService(greeterUsecase)
	// grpcServer 创建用于服务间调用的 server
	grpcServer := server.NewGRPCServer(confServer, greeterService, logger)
	// httpServer 创建 http 服务
	httpServer := server.NewHTTPServer(confServer, greeterService, logger)

	// 配置 nacos 节点信息
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}

	// 配置服务提供者的相关信息，名称空间、服务名等
	cc := constant.ClientConfig{
		NamespaceId:         "public",
		AppName:             "helloWorld",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	r := nacos.New(client)

	Name = "helloWorld-server"

	app := newApp1(logger, grpcServer, httpServer, r)

	fmt.Println("app", app)
	fmt.Println("Name", Name)

	return app, func() {
		cleanup()
	}, nil
}
