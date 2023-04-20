//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

/*
 	"helloworld/app/user/internal/biz"
	"helloworld/app/user/internal/data"
	"helloworld/app/user/internal/conf"
	"helloworld/app/user/internal/server"
	"helloworld/app/user/internal/service"
	"github.com/google/wire"

*/
import (
	"github.com/google/wire"
	"helloworld/app/user/internal/biz"
	"helloworld/app/user/internal/conf"
	"helloworld/app/user/internal/data"
	"helloworld/app/user/internal/server"
	"helloworld/app/user/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	//"helloworld/app/user/internal/conf"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	// service.ProviderSet 声明了需要 biz.GreeterUsecase ，biz.ProviderSet 可以提供 biz.GreeterUsecase
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
	//panic(wire.Build( data.ProviderSet, server.ProviderSet,biz.ProviderSet, service.ProviderSet, newApp)) // 可以正常生成 wire_gen.go 文件

	// conf.ProviderSet 声明了依赖 data 包下的 Data 结构体，而 data.ProviderSet 声明了依赖 conf 包下的 Data 结构体，报  multiple bindings
	//panic(wire.Build(data.ProviderSet, server.ProviderSet,service.ProviderSet,biz.ProviderSet,  newApp))

}
