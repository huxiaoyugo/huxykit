package utils

import (
	"fmt"
	"net"
)

// 检测端口是否可用
func CheckPorts(port string) bool {
	var err error

	tcpAddress, err := net.ResolveTCPAddr("tcp4", ":"+port)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	listener, err := net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		fmt.Println(port, " 被占用")
		return false
	} else {
		fmt.Println(port, " 可以使用")
		listener.Close()
	}

	return true
}
