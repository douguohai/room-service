package service

import (
	hashMap "github.com/douguohai/gods/maps/hashmap"
	"sync"
)

//用于 userId2SocketIdMap、socketId2UserTokenMap 这两个对象初始化
var initOnce sync.Once

//锁 用于 roomId2RoomMap 这个对象初始化
var roomId2RoomMapOnce sync.Once

//锁 用于 userTokenUsersMap 这个对象初始化
var userTokenUsersMapOnce sync.Once

//获取用户唯一id和socket的映射集合
func getUserTokenSocketMapping() (*hashMap.Map, *hashMap.Map) {
	initOnce.Do(func() {
		userToken2SocketIdMap = hashMap.New()
		socketId2UserTokenMap = hashMap.New()
	})
	return userToken2SocketIdMap, socketId2UserTokenMap
}

//获取用户唯一id和socket的映射集合
func getRoomId2RoomMap() *hashMap.Map {
	roomId2RoomMapOnce.Do(func() {
		roomId2RoomMap = hashMap.New()
	})
	return roomId2RoomMap
}

//获取用户唯一id和用户基本信息的映射集合
func getUserTokenUsersMap() *hashMap.Map {
	userTokenUsersMapOnce.Do(func() {
		userTokenUsersMap = hashMap.New()
	})
	return userTokenUsersMap
}
