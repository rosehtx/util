package main

import (
	"fmt"
	"net"
)

func main() {
	// 监听地址 ： 注意，传递net.UDPAddr对象
	fmt.Println("启动udp://127.0.0.1:8501")
	listen, _ := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 8501,
	})
	defer listen.Close()
	for {
		// 接收信息
		var data [1024]byte
		n, addr, err := listen.ReadFromUDP(data[:]) // 接收数据
		fmt.Println("data : ", string(data[:n]), addr)
		if err != nil {
			continue
		}
		// 发送信息
		listen.WriteToUDP([]byte("this is udp server"), addr)
	}
}
