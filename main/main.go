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
	"time"
)

// Easier to get running with CORS. Thanks for help @Vindexus and @erkie
var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func main() {

	//定义服务，处理跨域问题
	server := socketIo.NewServer(&engineio.Options{
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
	server.OnDisconnect("/", func(s socketIo.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go server.Serve()
	defer server.Close()
	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("/Users/tianwen/Desktop/")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}


