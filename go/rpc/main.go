package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	//"net/rpc/jsonrpc"
	"log"
)

func main() {
	//注册对象
	rpc.Register(new(Goods))

	//固定写法
	rpc.HandleHTTP()
	//监听rpc服务地址
	listener, e := net.Listen("tcp", "localhost:1234")
	if e != nil {
		log.Fatal("Starting RPC-server -listen error:", e)
	}
	go http.Serve(listener, nil)

	// 利用通道读取的阻塞来执行上面协程
	done := make(chan bool)
	<-done
}

type Goods struct {
	CheckId   int
	Name string
}
type Params struct {
	Id   int
	Name string
}

/**
argReq 客户端请求的参数
reply  返回给客户端的参数
error  是方法的默认返回值，方法必须得有个返回值
 */
func (gg *Goods) FindById(argReq *Params, reply *Goods) error {
	fmt.Println("接收到请求信息 ：", argReq)
	reply.CheckId = argReq.Id + 1
	reply.Name    = argReq.Name
	return nil
}

func (g *Goods) GetByIdGoodsName(argReq *Params, reply *string) error {
	fmt.Println("接收到请求信息 ：", argReq)
	*reply = argReq.Name + "aaa"
	return nil
}
