package service

import (
	"encoding/json"
	"fmt"
	hashSet "github.com/douguohai/gods/sets/hashset"
	"github.com/douguohai/room-service/base"
	"github.com/douguohai/room-service/entity"
	"github.com/gofrs/uuid"
	socketIo "github.com/googollee/go-socket.io"
	"log"
	"strings"
)

// HandleCreateRoom 创建房间，设置房间主人信息，生成房间唯一标识，创建对应关系映射
//http://127.0.0.1:8000/socket.io/?token=1234&EIO=3&transport=polling&t=NzOmN7m
func HandleCreateRoom(s socketIo.Conn, msg string) {
	_, socketId2UserTokenMap := getUserTokenSocketMapping()
	userToken, ok := socketId2UserTokenMap.Get(s.ID())
	if !ok {
		handError(s, entity.Result{
			Code: base.UserTokenNotFindCode,
			Msg:  base.UserTokenNotFindMsg,
			Data: userToken,
		})
		return
	}

	//解析用户参数,进行信息校验
	userInput := entity.User{}
	_ = json.Unmarshal([]byte(msg), &userInput)
	//保存用户唯一标识和用户信息映射
	userSave, ok := getUserTokenUsersMap().Get(userInput.Uid)
	if !ok {
		//代表是新用户，初始化对应的创建房间和加入房间集合
		userInput.CreatedRoom = hashSet.New()
		userInput.JoinedRoom = hashSet.New()
		userSave = userInput
	}
	currentUser := userSave.(entity.User)
	//代表是已经存在的用户
	if !strings.EqualFold(currentUser.Uid, userToken.(string)) {
		handError(s, entity.Result{
			Code: base.UserTokenNotSameCode,
			Msg:  base.UserTokenNotSameMsg,
			Data: currentUser,
		})
		return
	}

	//随即生成uuid作为房间id，并初始化房间对象
	rid := uuid.Must(uuid.NewV4())
	room := entity.Room{
		Rid:       rid.String(),
		UserSet:   hashSet.New(),
		RoomOwner: currentUser.Uid,
	}
	//保存房间id和房间信息的映射
	getRoomId2RoomMap().Put(room.Rid, room)
	//更新当前用户缓存的基本信息
	currentUser.CreatedRoom.Add(room.Rid)
	getUserTokenUsersMap().Put(currentUser.Uid, currentUser)

	log.Printf("用户: %v socket: %v 创建房间成功 房间标识 %v", currentUser.Uid, s.ID(), room.Rid)
	s.Emit("created", room)

}

// HandleRoomJoin 加入房间 目前只支持 1对1
//JoinRoom = {
//      rid: "",
//      user: {
//        uid: "",
//        sdp: "",
//      },
//    };
func HandleRoomJoin(s socketIo.Conn, msg string) {
	joinRoom := entity.JoinRoom{}
	_ = json.Unmarshal([]byte(msg), &joinRoom)
	//根据房间id查找对应房间映射
	room, ok := getRoomId2RoomMap().Get(joinRoom.Rid)
	if !ok {
		handError(s, entity.Result{
			Code: base.RoomInfoNotFindCode,
			Msg:  fmt.Sprintf("根据房间id %v，未查询到房间信息", joinRoom.Rid),
		})
		return
	}
	room2, ok := room.(entity.Room)
	if !ok {
		handError(s, entity.Result{
			Code: base.TypeCaseErrorCode,
			Msg:  "房间类型转换错误",
		})
		return
	}
	//自己加入房间
	room2.UserSet.Add(joinRoom.User)

	base.GetDefaultServer().JoinRoom(s.Namespace(), room2.Rid, s)
	base.GetDefaultServer().BroadcastToRoom(s.Namespace(), room2.Rid, "joined", room2)
	fmt.Printf("用户 %v 成功加入房间 %v\n", joinRoom.User.Uid, joinRoom.Rid)
}

// HandleRoomLeave 离开房间
func HandleRoomLeave(s socketIo.Conn, msg string) {
	s.SetContext(msg)
	fmt.Println("success leave a room hello")
}

//清除空闲无人用的空闲房间
func cleanEmptyRoom(user entity.User) {
	//把参与房间中的自己信息去除掉
	//将空房间全部删除掉
	for _, val := range user.JoinedRoom.Values() {
		obj, ok := getRoomId2RoomMap().Get(val)
		if ok {
			room := obj.(entity.Room)
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
			room := obj.(entity.Room)
			room.UserSet.Remove(user)
			if room.UserSet.Size() == 0 {
				//删除该房间信息
				getRoomId2RoomMap().Remove(val)
			}
		}
	}
}
