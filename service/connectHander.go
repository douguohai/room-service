package service

import (
	"fmt"
	hashMap "github.com/douguohai/gods/maps/hashmap"
	"github.com/douguohai/room-service/model"
	socketIo "github.com/googollee/go-socket.io"
	"log"
	"strings"
)

//用户->socketId 关系映射
var userToken2SocketIdMap *hashMap.Map

// HandleConnected 服务端建立联系成功,添加用户标识和socket关系的相互映射
func HandleConnected(s socketIo.Conn) error {
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
			userToken2SocketIdMap.Put(userToken, s)
			log.Printf("userToken->socket 新增一条记录 用户标识为：%v socket 标识为: %v\n", userToken, s.ID())
			socketId2UserTokenMap.Put(s.ID(), userToken)
			log.Printf("socket->userToken 新增一条记录 用户标识为：%v socket 标识为: %v\n", userToken, s.ID())
			s.SetContext(userToken)
		}
	}
	return nil
}

// HandDisconnected 处理联系丢失情况或者链接中断的情况
func HandDisconnected(s socketIo.Conn, reason string) {
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

//异常处理逻辑
func handError(s socketIo.Conn, result model.Result) {
	log.Printf(result.ToString())
	s.Emit("errored", result)
}
