package objects

import (
	"fmt"
	"lib/objectstream"

	"../heartbeat"
)

//objeetstream 包是对 Go 语言 http 包的一个封装，用来把一些 http 函数的调用转换成读写流的形式方便处理
func putStream(object string) (*objectstream.PutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}

	return objectstream.NewPutStream(server, object), nil
}
