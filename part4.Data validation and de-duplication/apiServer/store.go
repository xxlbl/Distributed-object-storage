package objects

import (
	"fmt"
	"io"
	"lib/utils"
	"net/http"
	"net/url"

	"../locate"
)

func storeObject(r io.Reader, hash string, size int64) (int, error) {
	// 如果hash内容已经存在，不用再存，返回ok
	if locate.Exist(url.PathEscape(hash)) {
		return http.StatusOK, nil
	}

	stream, e := putStream(url.PathEscape(hash), size) // 写入tmp
	if e != nil {
		return http.StatusInternalServerError, e
	}

	//io.TeeRcader的功能类似Unix的tee命令
	//它有两个输入参数,分别是作为io.Reader 的r和作为 io. Writer 的 stream,
	//它返回的 reader 也是一个 io.Reader.
	//当 reader 被读取时,其实际的内容读取自r,同时会写入stream.
	reader := io.TeeReader(r, stream)
	//用 utils.CalculateHlash 从 reader 中读取数据的同时也写入了 stream.
	d := utils.CalculateHash(reader)
	if d != hash {
		stream.Commit(false) //delete
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch, calculated=%s, requested=%s", d, hash)
	}
	stream.Commit(true) //put
	return http.StatusOK, nil
}
