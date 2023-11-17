package pb

import (
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
)

var (
	_ gomock.Matcher = (*SaveRequest)(nil)
	_ gomock.Matcher = (*GetRequest)(nil)
	_ gomock.Matcher = (*DeleteRequest)(nil)
)

func (x *SaveRequest) Matches(y interface{}) bool   { return match(x, y) }
func (x *GetRequest) Matches(y interface{}) bool    { return match(x, y) }
func (x *DeleteRequest) Matches(y interface{}) bool { return match(x, y) }

func match(x, y interface{}) bool {
	p1, ok1 := x.(proto.Message)
	p2, ok2 := y.(proto.Message)
	if !ok1 || !ok2 {
		return false
	}

	return proto.Equal(p1, p2)
}
