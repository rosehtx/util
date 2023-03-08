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
			fmt.Println("listen err : ", err)
			continue // 这里只是跳出本次循环
		}
		// 设置短连接(10秒)
		//conn.SetReadDeadline(time.Now().Add(time.Duration(10)*time.Second))
		// 3. 处理用户的连接信息
		go handler(conn)
	}
}

// 处理用户的连接信息
func handler(c net.Conn) {
	defer c.Close() // 一定要写 ，关闭连接
	//方案二用到reader 这边提前建一个reader防止重复建
	reader  := bufio.NewReader(c)
	//心跳计时
	message := make(chan []byte)
	go HeartBeating(c, &message, 10)

	for {
		buffer := make([]byte, 1024)
		//方案一:采用listen.Accept()返回的链接直接Read读取
		//循环读取客户端消息 读取到字节数组的长度
		//bufferLen, err := c.Read(buffer)
		//if err != nil {
		//	fmt.Println("Read err : ", err)
		//	break
		//}
		//fmt.Println("read bufferLen ",bufferLen)
		//unpackConRead(buffer,&message)

		//方案二:采用io的reader自带Read方法读取
		bufferLen, err := reader.Read(buffer)
		if err != nil {
			fmt.Println("bufio NewReader err : ", err)
			break
		}
		fmt.Println("bufio NewReader bufferLen ",bufferLen)
		unpackConRead(buffer,&message)

		//涉及到一点buffer处理
		//msg, err := unpackIoRead(reader)
	}
}

//listen.Accept()数据处理
func unpackConRead(buffer []byte,messageChan *chan []byte){
	//获取包头长度
	var headStart int
	var dataStart int
	headStart    = 0//包头开始的位置
	dataStart 	 = 2//包体开始的位置，包头为2个字节所以永远是加2
	for {
		headLenByte := buffer[headStart:dataStart]//这是个字节数组
		bodyLen     := binary.BigEndian.Uint16(headLenByte)//包体的长度
		fmt.Println("bodyLen :",bodyLen)
		if bodyLen == 0{
			break
		}
		//具体包的内容
		dataEnd     := int(bodyLen) + dataStart
		bufferData  := buffer[dataStart:dataEnd]
		*messageChan <-bufferData

		headStart = dataEnd
		dataStart = headStart + 2
	}
}

//reader来处理
func unpackIoRead(reader *bufio.Reader) (string, error) {
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
			//conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
			clintMsg := string(msg)
			go writeClientMsg(clintMsg,conn)
		//链接过期时间自己调节
		case <- time.After(5 * time.Second):
			fmt.Println("conn dead")
			conn.Close()
			break heart
		}
	}
}

//这边要做一个类似分发的功能
func writeClientMsg(msg string,conn net.Conn)  {
	conn.Write([]byte("server get "+msg))
}





 