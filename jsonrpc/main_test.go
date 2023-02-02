package main

import (
	"fmt"
	"net/rpc/jsonrpc"
	"testing"
)

type GoodsRes struct {
	CheckId   int
	Name string
}

type ParamsReq struct {
	Id   int
	Name string
}

func TestRpc(t *testing.T) {
	conn, _ 	 := jsonrpc.Dial("tcp", "127.0.0.1:9500")
	defer conn.Close()
	resgoods 	 := &GoodsRes{}
	resgoodsName := ""
	reqParam := &ParamsReq{
		Id:   1,
		Name: "请求",
	}
	/**
	发起请求，调用服务端的方法
	reqParam 请求参数
	resgoods 接受参数
	*/
	err := conn.Call("goods.FindById", reqParam, resgoods)
	if err != nil {
		fmt.Println("err : ", err)
	}
	fmt.Println("resgoods : ", resgoods)
	errGoodsName := conn.Call("goods.GetByIdGoodsName", reqParam, &resgoodsName)
	if errGoodsName != nil {
		fmt.Println("err : ", errGoodsName)
	}
	fmt.Println("resgoodsName : ", resgoodsName)
}
