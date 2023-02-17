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
	msg := "test nian bao 666" // 字节切片
	msgLen := len(msg)
	length := int16(msgLen) // 长度2个字节
	fmt.Println("msgLen : ", msgLen)
	pkg := new(bytes.Buffer)
	binary.Write(pkg, binary.BigEndian, length)
	data := append(pkg.Bytes(), []byte(msg)...)
	// 2. 进行数据的发送&接收数据
	for i := 1; i < 10; i++ {
		conn.Write(data)
	}
	var pack [1024]byte
	n, _ := conn.Read(pack[:])
	fmt.Println("server return : ", string(pack[:n])) // 切片获取信息
	//3. 关闭 // 不关闭不会造成太大 ，如果服务端没有心跳会存在问题
}
