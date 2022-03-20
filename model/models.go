package model

import (
	"fmt"
	"github.com/douguohai/gods/sets"
	hashSet "github.com/douguohai/gods/sets/hashset"
)

// User 用户信息
type User struct {
	Uid         string   `json:"uid"`         //用户唯一标识
	Sdp         string   `json:"sdp"`         //webrtc 核心
	CreatedRoom sets.Set `json:"createdRoom"` //自己创建的房间集合
	JoinedRoom  sets.Set `json:"joinedRoom"`  //加入的房间集合
}

// ToString User
func (u *User) ToString() string {
	return fmt.Sprintf("[%v,%v,%v]", u.Uid, u.Sdp)
}

// Room 房间定义
type Room struct {
	Rid       string       `json:"rid"`       //房间id
	UserSet   *hashSet.Set `json:"userSet"`   //房间内用户信息
	RoomOwner string       `json:"roomOwner"` //房屋创建人员
}

// JoinRoom 申请加入人员的信息
type JoinRoom struct {
	Rid  string `json:"rid"`  //要加入的房间id
	User User   `json:"user"` //加入人员的信息
}

// Result 统一返回结果
type Result struct {
	Code int         `json:"code"` //状态码
	Msg  string      `json:"msg"`  //状态信息
	Data interface{} `json:"data"` //具体业务数据
}

// ToString Result
func (result *Result) ToString() string {
	return fmt.Sprintf("[%v,%v,%v]", result.Code, result.Msg, result.Data)
}

// CallSomeone 拨打电话id
type CallSomeone struct {
	FromUid   string `json:"fromUid"`   //来自用户id
	ToUid     string `json:"toUid"`     //目标用户id
	Offer     string `json:"offer"`     //发起的offer
	OfferType string `json:"offerType"` //发起的offer的类型
}

// AnswerSomeone 回复拨打电话信息
type AnswerSomeone struct {
	FromUid    string `json:"fromUid"`    //来自用户id
	ToUid      string `json:"toUid"`      //目标用户id
	Answer     string `json:"answer"`     //回复offer的 answer
	AnswerType string `json:"answerType"` //发起的offer的类型
}

// IceCandidate 交换
type IceCandidate struct {
	FromUid      string      `json:"fromUid"`      //来自用户id
	ToUid        string      `json:"toUid"`        //目标用户id
	IceCandidate interface{} `json:"iceCandidate"` //发起的offer的类型
}
