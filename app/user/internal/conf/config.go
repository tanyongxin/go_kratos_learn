package conf

import (
	"github.com/google/wire"
	"google.golang.org/protobuf/runtime/protoimpl"
	"helloworld/app/user/internal/data"
)

var ProviderSet = wire.NewSet(NewData)

type confData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Database *Data_Database `protobuf:"bytes,1,opt,name=database,proto3" json:"database,omitempty"`
	Redis    *Data_Redis    `protobuf:"bytes,2,opt,name=redis,proto3" json:"redis,omitempty"`
}

func NewData(d *data.Data) *confData {

	d1 := Data{}

	d2 := confData(d1)

	return &d2

}
