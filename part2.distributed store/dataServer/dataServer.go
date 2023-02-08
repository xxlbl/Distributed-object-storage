package main

import (
	"./heartbeat"
	"./locate"
	"./objects"
	"log"
	"net/http"
	"os"
)

func main() {
	go heartbeat.StartHeartbeat()
	// 监听定位查询
	go locate.StartLocate()
	// 存取
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
