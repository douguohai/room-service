package main

const (
	UserTokenNotSameCode int = 9000 + iota
	UserTokenNotFindCode
	RoomInfoNotFindCode
	TypeCaseErrorCode
)

const (
	UserTokenNotSameMsg string = "非法连接，根据socket 获取用户标识和用户传过来的用户标识不一致"
	UserTokenNotFindMsg string = "非法连接，根据socketId 未获取到用户信息"
)
