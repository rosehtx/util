package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("启动服务端 ： tcp://127.0.0.1:9501")
	// 1. 监听端口 tcp://0.0.0.0:9501  监听的网络主要以本机可用ip为主
	listen, err := net.Listen("tcp", "127.0.0.1:9501")
	if err != nil {
		fmt.Println("err : ", err)
		return // return 表示程序结束
	}
	for {
		// 2. 接收客户向服务端建立的连接
		conn, err := listen.Accept() // 可以与客户端建立连接 ， 如果没有连接挂起阻塞状态
		if err != nil {
			fmt.Println("err : ", err)
			return // return 表示程序结束
		}
		// 设置短连接(10秒)
		conn.SetReadDeadline(time.Now().Add(time.Duration(10)*time.Second))
		// 3. 处理用户的连接信息
		go handler(conn)
	}
}

// 处理用户的连接信息
func handler(c net.Conn) {
	defer c.Close() // 一定要写 ，关闭连接
	//方案一:采用listen.Accept()返回的链接直接Read读取
	buffer := make([]byte, 1024)
	//reader := bufio.NewReader(c)

	//心跳计时
	message := make(chan []byte)
	go HeartBeating(c, &message, 10)

	for {
		fmt.Println("============start read")
		//读取到字节数组的长度
		bufferLen, err := c.Read(buffer)
		if err != nil {
			fmt.Println("err : ", err)
			break
		}

		unpackConRead(buffer,bufferLen,&message)

		//var data [1024]byte // 数组 - 》定义每一次数据读取的量
		// Read(p []byte) 需要采用切片接收
		// 数组用 : 处理完之后会变为切片
		//n, err := bufio.NewReader(c).Read(data[:]) //n代表切片数据读取的位置

		//
		//msg, err := unpack(reader)
		//if err != nil {
		//	fmt.Println("err : ", err)
		//	break
		//}
		//fmt.Println("client data :", msg)
		// Write(b []byte) (n int, err error)
		c.Write([]byte("this is server"))
	}
}

//listen.Accept()数据处理
func unpackConRead(buffer []byte,byteLen int,messageChan *chan []byte){
	//获取包头长度
	var headStart int
	var dataStart int
	headStart    = 0//包头开始的位置
	dataStart 	 = 2//包体开始的位置，包头为2个字节所以永远是加2
	for {
		headLenByte := buffer[headStart:dataStart]//这是个字节数组
		headLen     := binary.BigEndian.Uint16(headLenByte)//包体的长度
		if headLen == 0{
			break
		}
		//具体包的内容
		dataEnd     := int(headLen) + dataStart
		bufferData  := buffer[dataStart:dataEnd]
		*messageChan <-bufferData

		headStart = dataEnd
		dataStart = headStart + 2
	}
}

func unpack(reader *bufio.Reader) (string, error) {
	// 包头解析
	lenByte, _ := reader.Peek(2) // 读取前 固定的几个字节的信息

	lengthBuff := bytes.NewBuffer(lenByte)
	var length int16
	err := binary.Read(lengthBuff, binary.BigEndian, &length)
	fmt.Println("length : ", length)
	if err != nil {
		return "", err
	}

	// 获取信息
	// 包头 + 数据长度 = length
	if int16(reader.Buffered()) < length + 2 {
		return "", errors.New("数据异常")
	}

	// 真正的读取
	pack  := make([]byte, int(2+length))
	_, err = reader.Read(pack)
	if err != nil {
		return "", err
	}
	return string(pack[2:]), nil
}

func HeartBeating(conn net.Conn,message *chan []byte, timeout int)  {
	heart:
	for  {
		select {
		case msg := <- *message:
			fmt.Println("包体:",string(msg))
			conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))

		case <- time.After(5 * time.Second):
			fmt.Println("conn dead")
			conn.Close()
			break heart
		}
	}
}





 