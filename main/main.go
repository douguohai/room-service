package main

import (
	"fmt"
	socketIo "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"log"
	"net/http"
	"runtime"
	"time"
)

// Easier to get running with CORS. Thanks for help @Vindexus and @erkie
var allowOriginFunc = func(r *http.Request) bool {
	return true
}

var server *socketIo.Server = nil

func init() {
	//定义服务，处理跨域问题
	server = socketIo.NewServer(&engineio.Options{
		PingTimeout:  time.Second * 2,
		PingInterval: time.Millisecond * 20,
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})
}

func main() {

	defer func() {
		if err := recover(); err != nil {
			for i := 3; ; i++ {
				pc, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				fmt.Println(pc, file, line)
			}
			return
		}
	}()

	server.OnConnect("/", handleConnected)

	//创建房间
	server.OnEvent("/", "create", handleCreateRoom)

	//进入房间
	server.OnEvent("/", "join", handleRoomJoin)

	//离开房间
	server.OnEvent("/", "leave", handleRoomLeave)

	//联系出错
	server.OnError("/", func(s socketIo.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	//联系丢失
	server.OnDisconnect("/", handDisconnected)

	go server.Serve()
	defer server.Close()
	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("/Users/tianwen/Desktop/")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

// GetDefaultServer 获取socketIo 的server
func GetDefaultServer() *socketIo.Server {
	return server
}
