package controller

import (
	"fmt"
	"github.com/douguohai/room-service/base"
	"github.com/douguohai/room-service/service"
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
	base.Service = server
}

func StartSocket() {

	defer func() {
		if err := recover(); err != nil {
			for i := 3; ; i++ {
				pc, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				fmt.Println(pc, file, line)
			}
		}
	}()

	server.OnConnect("/", service.HandleConnected)

	//创建房间
	server.OnEvent("/", "create", service.HandleCreateRoom)

	//进入房间
	server.OnEvent("/", "join", service.HandleRoomJoin)

	//离开房间
	server.OnEvent("/", "leave", service.HandleRoomLeave)

	//拨打电话
	server.OnEvent("/", "handCall", service.HandCall)

	//拨打电话
	server.OnEvent("/", "handAnswer", service.HandAnswer)

	//交换信令
	server.OnEvent("/", "iceCandidate", service.HandIceCandidate)

	//联系出错
	server.OnError("/", func(s socketIo.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	//联系丢失
	server.OnDisconnect("/", service.HandDisconnected)

	go server.Serve()
	defer server.Close()
	http.Handle("/socket.io/", server)
	//http.Handle("/", http.FileServer(http.Dir("/Users/tianwen/Desktop/")))
	http.Handle("/login", service.CrosMiddleware(service.Login))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
