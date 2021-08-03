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
var getRid2RoomOnce sync.Once

//用于 id2UserMap 这个对象初始化
var id2UserMapOnce sync.Once

//用户->socketId 关系映射
var userSocketMap *hashMap.Map

//socketId->用户 关系映射
var socketUserMap *hashMap.Map

//roomId和room 关系映射
var rid2RoomMap *hashMap.Map

//用户唯一标识和用户信息映射
var id2UserMap *hashMap.Map

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
			s.SetContext(userSign)
		}
	}
	return nil
}

//获取用户唯一id和socket的映射集合
func getUser2Socket() (*hashMap.Map, *hashMap.Map) {
	once.Do(func() {
		userSocketMap = hashMap.New()
		socketUserMap = hashMap.New()
	})
	return userSocketMap, socketUserMap
}

//获取用户唯一id和socket的映射集合
func getRid2RoomMap() *hashMap.Map {
	getRid2RoomOnce.Do(func() {
		rid2RoomMap = hashMap.New()
	})
	return rid2RoomMap
}

//获取用户唯一id和用户基本信息的映射集合
func getId2RoomMap() *hashMap.Map {
	id2UserMapOnce.Do(func() {
		id2UserMap = hashMap.New()
	})
	return id2UserMap
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
	userInput := model.User{}
	_ = json.Unmarshal([]byte(msg), &userInput)
	//保存用户唯一标识和用户信息映射
	userSave, ok := getId2RoomMap().Get(userInput.Uid)
	if !ok {
		//代表是新用户，初始化对应的创建房间和加入房间集合
		userInput.CreatedRoom = hashSet.New()
		userInput.JoinedRoom = hashSet.New()
		userSave = userInput
	}
	user := userSave.(model.User)
	//代表是已经存在的用户
	if !strings.EqualFold(user.Uid, userSign.(string)) {
		log.Print("非法连接，根据socket 获取用户标识和用户传过来的用户标识不一致")
	}

	userSet := hashSet.New()
	//初始化set，放入set集合中
	userSet.Add(user)
	//随即生成uuid作为房间id
	rid := uuid.Must(uuid.NewV4())
	room := model.Room{
		Rid:       rid.String(),
		UserSet:   userSet,
		RoomOwner: user.Uid,
	}
	//保存房间号映射
	Rid2RoomMap := getRid2RoomMap()
	Rid2RoomMap.Put(room.Rid, room)
	//重新更新用户缓存的基本信息
	user.CreatedRoom.Add(room.Rid)
	getId2RoomMap().Put(user.Uid, user)

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

//处理呼叫
func handCall(s socketIo.Conn, msg string) {
	//根据msg中对方的人员信息判断对方是否在线
	//在线的话，给对方发送通通，将当前人员的基本信息发送过去
	//不在线的话，直接给申请方发送联系人不在线提示
}

//处理拒绝
func handReject(msg string) {
	//根据msg中对方的人员信息判断对方是否在线
	//在线的话，直接给申请方发送联系人不在线提示
}

//处理联系丢失情况或者链接中断的情况
func handDisconnected(s socketIo.Conn, reason string) {
	//清除对应socket和用户之间的关联信息
	userSocketMap, socketUserMap = getUser2Socket()
	_, ok := socketUserMap.Get(s.ID())
	if ok {
		socketUserMap.Remove(s.ID())
	}
	userSign := s.Context().(string)
	_, ok2 := userSocketMap.Get(userSign)
	if ok2 {
		userSocketMap.Remove(userSign)
	}
	fmt.Printf("已经清除对应人员缓存的信息 socket: %v 用户名: %v\n", s.ID(), userSign)
	//清除对应的空闲房间信息
	user, ok2 := getId2RoomMap().Get(userSign)
	if ok2 {
		go cleanEmptyRoom(user.(model.User))
	}
}

//清除对用的空闲房间
func cleanEmptyRoom(user model.User) {
	//把参与房间中的自己信息去除掉
	//将空房间全部删除掉
	for _, val := range user.JoinedRoom.Values() {
		obj, ok := getRid2RoomMap().Get(val)
		if ok {
			room := obj.(model.Room)
			room.UserSet.Remove(user)
			if room.UserSet.Size() == 0 {
				//删除该房间信息
				getRid2RoomMap().Remove(val)
			}
		}
	}
	for _, val := range user.CreatedRoom.Values() {
		obj, ok := getRid2RoomMap().Get(val)
		if ok {
			room := obj.(model.Room)
			room.UserSet.Remove(user)
			if room.UserSet.Size() == 0 {
				//删除该房间信息
				getRid2RoomMap().Remove(val)
			}
		}
	}
}
