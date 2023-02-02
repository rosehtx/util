package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

/**
这边是一个注册方法调用的示例
 */
func main() {
	// 注册服务   =》》 hash表
	// RegisterName （name,recv）
	// （name 服务标识    recv 具体的服务
	rpc.RegisterName("goods", new(Goods))

	// 开启rpc监听 -》端口和ip
	listen, _ := net.Listen("tcp", "127.0.0.1:9500")
	defer listen.Close()

	for  {
		// 建立连接
		conn, e := listen.Accept()
		if e != nil {
			continue
		}
		go jsonrpc.ServeConn(conn)
	}

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
