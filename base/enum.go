package base

const (
	SUCCESSCode          int = 0
	FailCode             int = -1
	UserTokenNotSameCode int = 9000 + iota
	UserTokenNotFindCode
	RoomInfoNotFindCode
	TypeCaseErrorCode
	UserSocketNotFindCode
)

const (
	SuccessMsg               string = "操作成功"
	FailMsg                  string = "操作失败"
	UserTokenNotSameMsg      string = "非法连接，根据socket 获取用户标识和用户传过来的用户标识不一致"
	UserTokenNotFindMsg      string = "非法连接，根据socketId 未获取到用户信息"
	ParameterParsingErrorMsg string = "参数解析异常"
)
