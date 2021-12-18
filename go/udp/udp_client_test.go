package main

import (
	"fmt"
	"net"
	"testing"
)

func TestUdpClient(t *testing.T) {
	socket, _ := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 8501,
	})
	defer socket.Close()
	// 发送信息
	socket.Write([]byte("hello upd server"))
	// 读取数据
	var data [1024]byte
	n, addr, _ := socket.ReadFromUDP(data[:]) // 接收数据
	fmt.Println("data : ", string(data[:n]), addr)
}
