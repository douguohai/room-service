package service

import (
	hashMap "github.com/douguohai/gods/maps/hashmap"
	"sync"
)

//socketId->用户 关系映射
var socketId2UserTokenMap *hashMap.Map

//用户id->socketId 关系映射
var userToken2SocketIdMap *hashMap.Map

//用于 userId2SocketIdMap、socketId2UserTokenMap 这两个对象初始化
var initOnce sync.Once

//roomId和room 关系映射
var roomId2RoomMap *hashMap.Map

//锁 用于 roomId2RoomMap 这个对象初始化
var roomId2RoomMapOnce sync.Once

//用户唯一标识和用户信息映射
var userTokenUsersMap *hashMap.Map

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

//房间唯一id 和 房间的映射关系
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
