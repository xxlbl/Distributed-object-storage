package objects

import (
	"io"
	"net/http"
)

func storeObject(r io.Reader, object string) (int, error) {
	// stream是一个io.write接口
	stream, e := putStream(object) // 创建一个写入文件流
	if e != nil {
		return http.StatusServiceUnavailable, e
	}
	// 将对象内容复制到文件写入流中
	io.Copy(stream, r)
	e = stream.Close()
	if e != nil {
		return http.StatusInternalServerError, e
	}
	return http.StatusOK, nil
}
