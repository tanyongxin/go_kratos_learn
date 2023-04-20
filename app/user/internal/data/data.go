package data

import (
	"helloworld/app/user/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.通过 NewData 和 NewGreeterRepo 都是 go 中的函数，函数的参数声明了依赖关系
// 比如 NewData 函数需要 conf.Data，因此生成的代码中会先创建出 conf.Data,其返回值作为提供者被其他 ProviderSet 所依赖
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	// TODO wrapped database client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{}, cleanup, nil
}
