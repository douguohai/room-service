package base

import socketIo "github.com/googollee/go-socket.io"

var Service *socketIo.Server = nil

// GetDefaultServer 获取socketIo 的server
func GetDefaultServer() *socketIo.Server {
	return Service
}
