package main

import (
	"fmt"
	"net/rpc"
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
	client, _ := rpc.DialHTTP("tcp", "127.0.0.1:1234")
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
	err := client.Call("Goods.FindById", reqParam, resgoods)
	if err != nil {
		fmt.Println("err : ", err)
	}
	fmt.Println("resgoods : ", resgoods)
	errGoodsName := client.Call("Goods.GetByIdGoodsName", reqParam, &resgoodsName)
	if errGoodsName != nil {
		fmt.Println("err : ", errGoodsName)
	}
	fmt.Println("resgoodsName : ", resgoodsName)
}
