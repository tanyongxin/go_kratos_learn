package rpcTest

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"testing"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"helloworld/api/helloworld/v1"

	"github.com/go-kratos/kratos/v2/log"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"

	_ "go.uber.org/automaxprocs"
)

// pass
func TestNacos1(t *testing.T) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}

	cc := &constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	cli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		//log.Panic(err)
		fmt.Println(err)
	}

	// newClient 可以重用，可以安全的在多个协程中使用
	newClient, err := http.NewClient(context.Background(), http.WithEndpoint("discovery:///helloWorld-server.http"), http.WithDiscovery(nacos.New(cli)), http.WithBlock())

	if err != nil {
		log.Fatal(err)
		//panic(err)

	}
	defer newClient.Close()

	client := v1.NewGreeterHTTPClient(newClient)
	reply, err := client.SayHello(context.Background(), &v1.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
		//panic(err)
	}
	//log.Printf("[grpc] SayHello %+v\n", reply)
	fmt.Printf("[http] SayHello %+v\n", reply)
}

// 通过。grpc client 的调用方式
func TestNacos(t *testing.T) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}

	cc := &constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	cli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		//log.Panic(err)
		fmt.Println(err)
	}

	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///helloWorld-server.grpc"), // 如果服务名称不对会报，Zero endpoint found,refused to write, instances: []
		grpc.WithDiscovery(nacos.New(cli)),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := v1.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &v1.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("[grpc] SayHello %+v\n", reply)
	fmt.Printf("[grpc] SayHello %+v\n", reply)
}

// 测试 rpc 调用
func TestRpc(t *testing.T) {

	conn, err := transgrpc.DialInsecure(
		context.Background(),
		transgrpc.WithEndpoint("127.0.0.1:9000"),
		transgrpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := v1.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &v1.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("[grpc] SayHello %+v\n", reply)
	fmt.Printf("[grpc] SayHello %+v\n", reply)

	// returns error
	_, err = client.SayHello(context.Background(), &v1.HelloRequest{Name: "error"})
	if err != nil {
		//log.Printf("[grpc] SayHello error: %v\n", err)
		fmt.Printf("[grpc] SayHello error: %v\n", err)
	}
	if errors.IsBadRequest(err) {
		//log.Printf("[grpc] SayHello error is invalid argument: %v\n", err)
		fmt.Printf("[grpc] SayHello error is invalid argument: %v\n", err)
	}

}
