package x_net

import (
	"github.com/golang/protobuf/proto"
	"reflect"
)

type Handler func(data interface{}, c Conner)
type HandlerModel struct {
	T reflect.Type
	H Handler
}

var handlerMap map[uint64]*HandlerModel //k:msgId v:HandlerModel

func init() {
	handlerMap = make(map[uint64]*HandlerModel)
}

func RegisterHandler(id uint64, msg proto.Message, h Handler) {
	t := reflect.TypeOf(msg)
	_, ok := handlerMap[id]
	if ok {
		return
	}
	handlerMap[id] = &HandlerModel{T: t, H: h}
}

func GetHandlerModel(msgId uint64) (hm *HandlerModel, ok bool) {
	hm, ok = handlerMap[msgId]
	return
}
