package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"testing"
)

// tcp客户端
func TestTcpClient(t *testing.T) {
	// 1. 创建建立连接
	conn, _ := net.Dial("tcp", "127.0.0.1:9501")
	fmt.Println("与tcp://127.0.0.1:9501建立连接")
	defer conn.Close()

	// int16 -> 2 个字节
	msg := "test 粘包 666" // 字节切片
	msgLen := len(msg)
	length := int16(msgLen) // 长度2个字节
	//fmt.Println("msgLen : ", msgLen)
	pkg := new(bytes.Buffer)
	binary.Write(pkg, binary.BigEndian, length)
	data := append(pkg.Bytes(), []byte(msg)...)

	// 2. 进行数据的发送&接收数据
	var pack [1024]byte
	for i := 1; i < 10; i++ {
		conn.Write(data)
		//这样等待请求结果没有粘包现象
		//n, _ := conn.Read(pack[:])
		//fmt.Println("server return : ", string(pack[:n])) // 切片获取信息
	}
	//上面不读取这边只读取一次，会出现粘包现象
	n, _ := conn.Read(pack[:])
	fmt.Println("server return : ", string(pack[:n])) // 切片获取信息

	//3. 关闭 // 不关闭不会造成太大 ，如果服务端没有心跳会存在问题

	// 利用通道读取的阻塞来执行上面协程
	//done := make(chan bool)
	//<-done
}
