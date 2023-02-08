package objects

import (
	"fmt"
	"io"
	"lib/objectstream"

	"../locate"
)

func getStream(object string) (io.Reader, error) {
	server := locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("object %s locate fail", object)
	}
	// 创建一个服务器流
	return objectstream.NewGetStream(server, object)
}
