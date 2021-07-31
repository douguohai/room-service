package model

import (
	"fmt"
	"github.com/douguohai/gods/sets"
)

//用户信息
type User struct {
	Uid      string `json:"uid"`      //用户唯一标识
	Nickname string `json:"nickname"` //用户昵称
	Sdp      string `json:"sdp"`      //webrtc 核心
}

func (u *User) ToString() string {
	return fmt.Sprintf("[%v,%v,%v]", u.Uid, u.Nickname, u.Sdp)
}

//房间定义
type Room struct {
	Rid       string   `json:"rid"`       //房间id
	UserSet   sets.Set `json:"userSet"`   //房间内用户信息
	RoomOwner string   `json:"roomOwner"` //房屋创建人员
}

//申请加入人员的信息
type JoinRoom struct {
	Rid  string `json:"rid"`  //要加入的房间id
	User User   `json:"user"` //加入人员的信息
}
