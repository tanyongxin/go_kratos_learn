package data

import (
	"entgo.io/ent/dialect/sql"
	"helloworld/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	// TODO wrapped database client。数据库客户端的包装
	Driver *sql.Driver
}

// NewData .
func NewData(conf *conf.Data, logger log.Logger) (*Data, func(), error) {
	// drv 的类型为 Driver，Driver 继承了 H:/GoProject/pkg/mod/entgo.io/ent@v0.11.8/dialect/sql/driver.go:91 Conn 结构体
	//
	drv, err := sql.Open(
		conf.Database.Driver,
		conf.Database.Source,
	)

	if err != nil {
		log.Errorf("failed opening connection to sqlite: %v", err)
		return nil, nil, err
	}

	// cleanup 是关闭数据源时的回调函数
	cleanup := func() {
		// 关闭链接
		e := drv.Close()
		if e != nil {
			log.Errorf("failed opening connection to sqlite: %v", err)
		}
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{drv}, cleanup, nil
}
