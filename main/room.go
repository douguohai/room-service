package main

import (
	"encoding/json"
	"fmt"
	hashMap "github.com/douguohai/gods/maps/hashmap"
	hashSet "github.com/douguohai/gods/sets/hashset"
	"github.com/gofrs/uuid"
	socketIo "github.com/googollee/go-socket.io"
	"log"
	"room-service/model"
	"strings"
	"sync"
)

//用于 userSocketMap、socketUserMap 这两个对象初始化
var once sync.Once

//用于 rid2RoomMap 这个对象初始化
var two sync.Once

//用户->socketId 关系映射
var userSocketMap *hashMap.Map

//socketId->用户 关系映射
var socketUserMap *hashMap.Map

//roomId和room 关系映射
var rid2RoomMap *hashMap.Map

//服务端建立联系成功,添加用户标识和socket关系的相互映射
func handleConnected(s socketIo.Conn) error {
	fmt.Println("connected:", s.ID(), s.URL().RawQuery)
	var queryStr = s.URL().RawQuery
	var queryArr = strings.Split(queryStr, "&")
	if len(queryArr) != 0 {
		paramMap := make(map[string]string)
		for _, query := range queryArr {
			//拆分参数为map
			keyVal := strings.Split(query, "=")
			if len(keyVal) == 2 {
				paramMap[keyVal[0]] = keyVal[1]
			}
		}
		userSign, ok := paramMap["token"]
		if ok {
			userSocketMap, socketUserMap := getUser2Socket()
			userSocketMap.Put(userSign, s.ID())
			log.Printf("userSign->socket 新增一条记录 用户标识为：%v socket 标识为: %v\n", userSign, s.ID())
			socketUserMap.Put(s.ID(), userSign)
			log.Printf("socket->userSign 新增一条记录 用户标识为：%v socket 标识为: %v\n", userSign, s.ID())
		}
	}
	return nil
}

//获取用户唯一id和socket的映射集合
func getUser2Socket() (*hashMap.Map, *hashMap.Map) {
	var initMap = func() {
		userSocketMap = hashMap.New()
		socketUserMap = hashMap.New()
	}
	once.Do(initMap)
	return userSocketMap, socketUserMap
}

//获取用户唯一id和socket的映射集合
func getRid2RoomMap() *hashMap.Map {
	var init = func() {
		rid2RoomMap = hashMap.New()
	}
	two.Do(init)
	return rid2RoomMap
}

//创建房间，设置房间主人信息，生成房间唯一标识，创建对应关系映射
func handleCreateRoom(s socketIo.Conn, msg string) {
	fmt.Println("create:", msg, s.ID())
	_, socketUserMap := getUser2Socket()
	userSign, ok := socketUserMap.Get(s.ID())
	if !ok {
		log.Print("非法连接，根据socket 获取不到唯一用户标识")
	}
	//解析用户参数,进行信息校验
	user := &model.User{}
	_ = json.Unmarshal([]byte(msg), &user)
	if !strings.EqualFold(user.Uid, userSign.(string)) {
		log.Print("非法连接，根据socket 获取用户标识和用户传过来的用户标识不一致")
	}
	//初始化set，放入set集合中
	userSet := hashSet.New()
	userSet.Add(&user)
	//随即生成uuid作为房间id
	rid := uuid.Must(uuid.NewV4())
	room := &model.Room{
		Rid:       rid.String(),
		UserSet:   userSet,
		RoomOwner: user.Uid,
	}
	//保存房间号映射
	Rid2RoomMap := getRid2RoomMap()
	Rid2RoomMap.Put(room.Rid, &room)
	log.Printf("用户: %v socket: %v 创建房间成功 房间标识 %v", user.Uid, s.ID(), room.Rid)
	s.Emit("created", room)
}

//加入房间 目前只支持 1对1
func handleRoomJoin(s socketIo.Conn, msg string) {
	joinRoom := model.JoinRoom{}
	_ = json.Unmarshal([]byte(msg), &joinRoom)
	//根据房间id查找对应房间映射
	room, ok := getRid2RoomMap().Get(joinRoom.Rid)
	if !ok {
		log.Printf("根据房间id %v，未查询到房间信息", joinRoom.Rid)
	}
	room2, ok := room.(model.Room)
	if !ok {
		log.Printf("房间类型转换错误")
	}
	//获取房间主人信息，再将自己信息添加进入房间集合中
	for _, item := range room2.UserSet.Values() {
		user, ok := item.(model.User)
		if !ok {
			log.Printf("用户类型转换错误")
		}
		var roomOwner = make(map[string]string)
		roomOwner["sdp"] = user.Sdp
		roomOwner["uid"] = user.Sdp
		s.Emit("joined", roomOwner)
	}
	room2.UserSet.Add(joinRoom.User)
	//加入房间信息，获取对方的sid
	fmt.Printf("用户 %v 成功加入房间 %v\n", joinRoom.User.Uid, joinRoom.Rid)
}

//离开房间
func handleRoomLeave(s socketIo.Conn, msg string) {
	s.SetContext(msg)
	fmt.Println("success leave a room hello")
}
