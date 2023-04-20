package data

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"fmt"

	"helloworld/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// greeterRepo 实现了 internal/biz/greeter.go:23 接口
type greeterRepo struct {
	data *Data
	log  *log.Helper
}

type user struct {
	id       int
	userName string
	userPass string
}

// NewGreeterRepo .
func NewGreeterRepo(data *Data, logger log.Logger) biz.GreeterRepo {
	return &greeterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// 增
func (r *greeterRepo) Save(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {

	driver := r.data.Driver
	//rows := &sql.Rows{}
	//var args interface{} = make([]interface{},0)
	//err := driver.Query(ctx,"select userName from `user` limit 1",args,rows) // 不带条件的查询
	//
	//// userName 的类型是指针类型，因此后面取值时要用 *
	//res := new(string)
	//args := make([]interface{},2)
	//args[0] = 1
	//args[1] = "admin"
	//err := driver.Query(ctx,"select userPass from `user` where id = ? and userName = ?",args,rows)
	//rows.ColumnScanner.Next() // 调用 rows.ColumnScanner.Scan 需要先调用 rows.ColumnScanner.Next() 注释掉该行代码会报：Scan called without calling Next
	//err = rows.ColumnScanner.Scan(res)
	//
	//if err != nil {
	//	fmt.Println("query error",err)
	//}else {
	//	fmt.Println("userName is ",*res)
	//}

	selector := &sql.Selector{}
	selector.WithContext(ctx)
	table := sql.Table("user")
	selector.From(table)
	//
	users := []*user{}
	//
	//
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   "user",
			Columns: []string{"id", "userName", "userPass"},
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: "id",
			},
		},
		From:   selector,
		Unique: true,
	}
	//
	// 将数据库中的列和实体字段的类型做个映射
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		values := make([]interface{}, len(columns))
		for i := range columns {
			switch columns[i] {
			case "id":
				values[i] = &sql.NullInt64{}
			case "userName", "userPass":
				values[i] = &sql.NullString{}
			default:
				return nil, fmt.Errorf("unexpected column %q for type Article", columns[i])
			}
		}
		return values, nil
	}
	// _spec.Assign 会将查询返回的列设置到实体中
	_spec.Assign = func(columns []string, values []interface{}) error {
		user := &user{}
		for i := range columns {
			switch columns[i] {
			case "id":
				value, ok := values[i].(*sql.NullInt64)
				if !ok {
					return fmt.Errorf("unexpected type %T for field id", value)
				}
				user.id = int(int64(value.Int64))
			case "userName":
				if value, ok := values[i].(*sql.NullString); !ok {
					return fmt.Errorf("unexpected type %T for field title", values[i])
				} else if value.Valid {
					user.userName = value.String
				}
			case "userPass":
				if value, ok := values[i].(*sql.NullString); !ok {
					return fmt.Errorf("unexpected type %T for field content", values[i])
				} else if value.Valid {
					user.userPass = value.String
				}
			}
		}
		users = append(users, user)
		return nil
	}

	err := sqlgraph.QueryNodes(ctx, driver, _spec)

	if err == nil {
		for _, v := range users {
			fmt.Println(*v)
		}
		//res := new(string)
		//rows.ColumnScanner.Next() // 调用 rows.ColumnScanner.Scan 需要先调用 rows.ColumnScanner.Next() 注释掉该行代码会报：Scan called without calling Next
		//err = rows.ColumnScanner.Scan(res)
		//if err == nil {
		//	fmt.Println(*res)
		//}else {
		//	fmt.Println(err)
		//}

	} else {
		fmt.Println(err)
	}

	return g, nil
}

// 改
func (r *greeterRepo) Update(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	return g, nil
}

// 根据 id 查找
func (r *greeterRepo) FindByID(context.Context, int64) (*biz.Greeter, error) {
	return nil, nil
}

func (r *greeterRepo) ListByHello(context.Context, string) ([]*biz.Greeter, error) {
	return nil, nil
}

// 查询列表
func (r *greeterRepo) ListAll(context.Context) ([]*biz.Greeter, error) {
	return nil, nil
}
