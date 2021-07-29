package model

import (
	"github.com/douguohai/gods/sets"
)

//用户信息
type User struct {
	Uid      string `json:"uid"`      //用户唯一标识
	Nickname string `json:"nickname"` //用户昵称
}

func (u *User) ToString() string {
	return u.Nickname + " " + u.Uid
}

//房间定义
type Room struct {
	Rid     string   `json:"rid"`     //房间id
	UserSet sets.Set `json:"userSet"` //房间内用户信息
}
