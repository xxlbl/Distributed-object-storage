package main

import (
	"log"
	"net/http"
	"os"

	"test/objects"
)

func main() {
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}

//测试

//LISTEN_ADDRESS=:12345 STORAGE_ROOT=./tmp go run server.go

// $ curl -v http://localhost:12345/objects/test1 -X PUT -d "this is a object"
// *   Trying 127.0.0.1:12345...
// * Connected to localhost (127.0.0.1) port 12345 (#0)
// > PUT /objects/test1 HTTP/1.1
// > Host: localhost:12345
// > User-Agent: curl/7.79.1
// > Accept: */*
// > Content-Length: 16
// > Content-Type: application/x-www-form-urlencoded
// >
// * Mark bundle as not supporting multiuse
// < HTTP/1.1 200 OK
// < Date: Wed, 08 Feb 2023 04:31:55 GMT
// < Content-Length: 0
// <
// * Connection #0 to host localhost left intact

// $ curl -v http://localhost:12345/objects/test1
// *   Trying 127.0.0.1:12345...
// * Connected to localhost (127.0.0.1) port 12345 (#0)
// > GET /objects/test1 HTTP/1.1
// > Host: localhost:12345
// > User-Agent: curl/7.79.1
// > Accept: */*
// >
// * Mark bundle as not supporting multiuse
// < HTTP/1.1 200 OK
// < Date: Wed, 08 Feb 2023 04:32:09 GMT
// < Content-Length: 16
// < Content-Type: text/plain; charset=utf-8
// <
// * Connection #0 to host localhost left intact
// this is a object%
