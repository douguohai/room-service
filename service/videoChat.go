package service

import (
	"encoding/json"
	"fmt"
	"github.com/douguohai/room-service/base"
	"github.com/douguohai/room-service/entity"
	socketIo "github.com/googollee/go-socket.io"
)

// HandCall 视频呼叫
func HandCall(s socketIo.Conn, msg string) {
	callSomeone := entity.CallSomeone{}
	_ = json.Unmarshal([]byte(msg), &callSomeone)
	//根据msg中对方的人员信息判断对方是否在线
	userToken2SocketIdMap, _ := getUserTokenSocketMapping()
	toSocket, ok := userToken2SocketIdMap.Get(callSomeone.ToUid)
	//不在线的话，直接给申请方发送联系人不在线提示
	if !ok {
		handError(s, entity.Result{
			Code: base.UserSocketNotFindCode,
			Msg:  "未发现拨号用户信息或者对方此时未在线",
		})
		return
	}
	callSomeone.FromUid = s.Context().(string)
	//在线的话，给对方发送通通，将offer信息同步过去
	toSocket2, ok := toSocket.(socketIo.Conn)
	toSocket2.Emit("callYou", callSomeone)
}

// HandAnswer 处理呼叫
func HandAnswer(s socketIo.Conn, msg string) {
	answerSomeone := entity.AnswerSomeone{}
	_ = json.Unmarshal([]byte(msg), &answerSomeone)
	fmt.Printf(answerSomeone.FromUid)
	//根据msg中对方的人员信息判断对方是否在线
	userToken2SocketIdMap, _ := getUserTokenSocketMapping()
	toSocket, ok := userToken2SocketIdMap.Get(answerSomeone.ToUid)
	//不在线的话，直接给申请方发送联系人不在线提示
	if !ok {
		handError(s, entity.Result{
			Code: base.UserSocketNotFindCode,
			Msg:  "未发现拨号用户信息或者对方此时未在线",
		})
		return
	}
	//在线的话，给对方发送通通，将offer信息同步过去
	toSocket2, ok := toSocket.(socketIo.Conn)
	toSocket2.Emit("answerYou", answerSomeone)
}

// HandIceCandidate 交换凭证
func HandIceCandidate(s socketIo.Conn, msg string) {
	iceCandidate := entity.IceCandidate{}
	_ = json.Unmarshal([]byte(msg), &iceCandidate)
	//根据msg中对方的人员信息判断对方是否在线
	userToken2SocketIdMap, _ := getUserTokenSocketMapping()
	toSocket, ok := userToken2SocketIdMap.Get(iceCandidate.ToUid)
	//不在线的话，直接给申请方发送联系人不在线提示
	if !ok {
		handError(s, entity.Result{
			Code: base.UserSocketNotFindCode,
			Msg:  "未发现拨号用户信息或者对方此时未在线",
		})
		return
	}
	//在线的话，给对方发送通通，将offer信息同步过去
	toSocket2, ok := toSocket.(socketIo.Conn)
	toSocket2.Emit("iceCandidate", iceCandidate)
}

// HandReject 处理拒绝
func HandReject(msg string) {
	//根据msg中对方的人员信息判断对方是否在线
	//在线的话，直接给申请方发送联系人不在线提示
}
