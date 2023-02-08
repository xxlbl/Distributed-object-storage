package locate

import (
	"lib/rabbitmq"
	"os"
	"strconv"
)

// 查文件属性
func Locate(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func StartLocate() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()
	for msg := range c {
		object, e := strconv.Unquote(string(msg.Body)) //查询的对象名称
		if e != nil {
			panic(e)
		}
		if Locate(os.Getenv("STORAGE_ROOT") + "/objects/" + object) {
			q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
			// 文件存在：将本存储服务地址发送给发查询消息的api服务
		}
	}
}
