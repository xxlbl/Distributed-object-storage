package objects

import (
	"log"
	"net/http"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	c, e := storeObject(r.Body, object) // 存储对象的内容和名字
	if e != nil {
		log.Println(e)
	}
	// 返回http code
	w.WriteHeader(c)
}
