package heartbeat

import (
	"lib/rabbitmq"
	"os"
	"strconv"
	"sync"
	"time"
)

// url—>接受到的时间
var dataServers = make(map[string]time.Time)
var mutex sync.Mutex //map并发读写 保护

func ListenHeartbeat() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER")) // 生成rabbitmq结构，控制相关api
	defer q.Close()
	q.Bind("apiServers") // 创建消息队列绑定apiservers，发给apiservers的消息都会发给q
	c := q.Consume()
	go removeExpiredDataServer() // 协程超时移除
	for msg := range c {
		dataServer, e := strconv.Unquote(string(msg.Body)) // 解析msg "..."，拿到数据服务地址
		if e != nil {
			panic(e)
		}
		mutex.Lock()
		dataServers[dataServer] = time.Now() // url+当前时间，写入map
		mutex.Unlock()
	}
}

func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second) // 5s检查一次
		mutex.Lock()
		for s, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) { //10s没有更新
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}

// 返回当前所有的数据服务url
func GetDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string, 0)
	for s, _ := range dataServers {
		ds = append(ds, s)
	}
	return ds
}
