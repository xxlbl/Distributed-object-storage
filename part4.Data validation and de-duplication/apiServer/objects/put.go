package objects

import (
	"lib/es"
	"lib/utils"
	"log"
	"net/http"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	size := utils.GetSizeFromHeader(r.Header)
	// 内容+hash+大小（需要确定tmp临时对象的大小）
	c, e := storeObject(r.Body, hash, size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(c)
		return
	}
	if c != http.StatusOK {
		w.WriteHeader(c)
		return
	}

	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	e = es.AddVersion(name, hash, size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
