package main

import (
	"encoding/json"
	"fmt"
	hashMap "github.com/douguohai/gods/maps/hashmap"
	socketIo "github.com/googollee/go-socket.io"
	"room-service/model"
	"sync"
)

//服务端建立联系成功
func handleConnected(s socketIo.Conn) error {
	fmt.Println("connected:", s.ID())
	return nil
}

//获取用户唯一id和socket的映射集合
func getUser2Socket() hashMap.Map  {
	var once sync.Once
	var myMap hashMap.Map
	once.Do(func() {
		myMap = hashMap.New()
	})

}

//创建房间
func handleCreateRoom(s socketIo.Conn, msg string) {
	fmt.Println("create:", msg, s.ID())
	//解析用户参数
	user := &model.User{}
	json.Unmarshal([]byte(msg), &user)
	//将当前socket和用户身份信息进行绑定
	hashMap := hashMap.New()
	hashMap.Put(s.ID(), &user)
	fmt.Println(user.ToString())
	fmt.Println("success create a room hello")
}

//加入房间
func handleRoomJoin(s socketIo.Conn, msg string) {
	fmt.Println("join:", msg, s.ID())
	fmt.Println("success join a room hello")
}

//离开房间
func handleRoomLeave(s socketIo.Conn, msg string) {
	s.SetContext(msg)
	fmt.Println("success leave a room hello")
}
