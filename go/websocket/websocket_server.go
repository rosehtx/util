package main

import (
	"fmt"
	"net/http"
	//注：要go get golang.org/x/net下载下来才能使用websocket
	//要用go mod模式引用外部包
	"golang.org/x/net/websocket"
)

func main() {
	fmt.Println("启动websocket://:7777")
	// 设定websocket的服务信息处理
	http.Handle("/", websocket.Handler(server))
	// 设定监听
	err :=http.ListenAndServe(":7777", websocket.Server{
		websocket.Config{},nil,server})
	if err != nil{
		fmt.Println(err)
	}
}

// 当有连接进来的时候底层会自动
func server(ws *websocket.Conn) {
	fmt.Println("new connection")
	data := make([]byte, 1024)
	for {
		// 读取信息
		d, err := ws.Read(data)
		if err != nil {
			fmt.Println("err : ", err)
			break
		}
		fmt.Println("读取到的信息 ", string(data[:d]))
		ws.Write([]byte("this is webscoket server"))
	}
}
