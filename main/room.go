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

//锁 用于 roomId2RoomMap 这个对象初始化
var roomId2RoomMapOnce sync.Once

//roomId和room 关系映射
var roomId2RoomMap *hashMap.Map

//锁 用于 userTokenUsersMap 这个对象初始化
var userTokenUsersMapOnce sync.Once

//用户唯一标识和用户信息映射
var userTokenUsersMap *hashMap.Map

//用于 userId2SocketIdMap、socketId2UserTokenMap 这两个对象初始化
var initOnce sync.Once

//用户->socketId 关系映射
var userToken2SocketIdMap *hashMap.Map

//socketId->用户 关系映射
var socketId2UserTokenMap *hashMap.Map

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
		//获取当前用户token信息
		userToken, ok := paramMap["token"]
		if ok {
			userToken2SocketIdMap, socketId2UserTokenMap := getUserTokenSocketMapping()
			userToken2SocketIdMap.Put(userToken, s.ID())
			log.Printf("userToken->socket 新增一条记录 用户标识为：%v socket 标识为: %v\n", userToken, s.ID())
			socketId2UserTokenMap.Put(s.ID(), userToken)
			log.Printf("socket->userToken 新增一条记录 用户标识为：%v socket 标识为: %v\n", userToken, s.ID())
			s.SetContext(userToken)
		}
	}
	return nil
}

//创建房间，设置房间主人信息，生成房间唯一标识，创建对应关系映射
//http://127.0.0.1:8000/socket.io/?token=1234&EIO=3&transport=polling&t=NzOmN7m
func handleCreateRoom(s socketIo.Conn, msg string) {
	_, socketId2UserTokenMap := getUserTokenSocketMapping()
	userToken, ok := socketId2UserTokenMap.Get(s.ID())
	if !ok {
		handError(s, model.Result{
			Code: UserTokenNotFindCode,
			Msg:  UserTokenNotFindMsg,
			Data: userToken,
		})
		return
	}

	//解析用户参数,进行信息校验
	userInput := model.User{}
	_ = json.Unmarshal([]byte(msg), &userInput)
	//保存用户唯一标识和用户信息映射
	userSave, ok := getUserTokenUsersMap().Get(userInput.Uid)
	if !ok {
		//代表是新用户，初始化对应的创建房间和加入房间集合
		userInput.CreatedRoom = hashSet.New()
		userInput.JoinedRoom = hashSet.New()
		userSave = userInput
	}
	currentUser := userSave.(model.User)
	//代表是已经存在的用户
	if !strings.EqualFold(currentUser.Uid, userToken.(string)) {
		handError(s, model.Result{
			Code: UserTokenNotSameCode,
			Msg:  UserTokenNotSameMsg,
			Data: currentUser,
		})
		return
	}

	//随即生成uuid作为房间id，并初始化房间对象
	rid := uuid.Must(uuid.NewV4())
	room := model.Room{
		Rid:       rid.String(),
		UserSet:   hashSet.New(),
		RoomOwner: currentUser.Uid,
	}
	//保存房间id和房间信息的映射
	roomId2RoomMap := getRoomId2RoomMap()
	roomId2RoomMap.Put(room.Rid, room)
	//更新当前用户缓存的基本信息
	currentUser.CreatedRoom.Add(room.Rid)
	getUserTokenUsersMap().Put(currentUser.Uid, currentUser)

	log.Printf("用户: %v socket: %v 创建房间成功 房间标识 %v", currentUser.Uid, s.ID(), room.Rid)
	s.Emit("created", room)

}

//加入房间 目前只支持 1对1
//JoinRoom = {
//      rid: "",
//      user: {
//        uid: "",
//        sdp: "",
//      },
//    };
func handleRoomJoin(s socketIo.Conn, msg string) {
	joinRoom := model.JoinRoom{}
	_ = json.Unmarshal([]byte(msg), &joinRoom)
	//根据房间id查找对应房间映射
	roomId2RoomMap := getRoomId2RoomMap()
	room, ok := roomId2RoomMap.Get(joinRoom.Rid)
	if !ok {
		handError(s, model.Result{
			Code: RoomInfoNotFindCode,
			Msg:  fmt.Sprintf("根据房间id %v，未查询到房间信息", joinRoom.Rid),
		})
		return
	}
	room2, ok := room.(model.Room)
	if !ok {
		handError(s, model.Result{
			Code: TypeCaseErrorCode,
			Msg:  "房间类型转换错误",
		})
		return
	}
	//自己加入房间
	room2.UserSet.Add(joinRoom.User)

	GetDefaultServer().JoinRoom(s.Namespace(), room2.Rid, s)
	GetDefaultServer().BroadcastToRoom(s.Namespace(), room2.Rid, "joined", room2)
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
	userToken2SocketIdMap, socketId2UserTokenMap = getUserTokenSocketMapping()
	_, ok := socketId2UserTokenMap.Get(s.ID())
	if ok {
		socketId2UserTokenMap.Remove(s.ID())
	}
	userToken := s.Context().(string)
	_, ok2 := userToken2SocketIdMap.Get(userToken)
	if ok2 {
		userToken2SocketIdMap.Remove(userToken)
	}
	fmt.Printf("已经清除对应人员缓存的信息 socket: %v 用户名: %v\n", s.ID(), userToken)
	//清除对应的空闲房间信息
	user, ok2 := getUserTokenUsersMap().Get(userToken)
	if ok2 {
		go cleanEmptyRoom(user.(model.User))
	}
}

//清除对用的空闲房间
func cleanEmptyRoom(user model.User) {
	//把参与房间中的自己信息去除掉
	//将空房间全部删除掉
	for _, val := range user.JoinedRoom.Values() {
		obj, ok := getRoomId2RoomMap().Get(val)
		if ok {
			room := obj.(model.Room)
			room.UserSet.Remove(user)
			if room.UserSet.Size() == 0 {
				//删除该房间信息
				getRoomId2RoomMap().Remove(val)
			}
		}
	}
	for _, val := range user.CreatedRoom.Values() {
		obj, ok := getRoomId2RoomMap().Get(val)
		if ok {
			room := obj.(model.Room)
			room.UserSet.Remove(user)
			if room.UserSet.Size() == 0 {
				//删除该房间信息
				getRoomId2RoomMap().Remove(val)
			}
		}
	}
}

//异常处理逻辑
func handError(s socketIo.Conn, result model.Result) {
	log.Printf(result.ToString())
	s.Emit("errored", result)
}
