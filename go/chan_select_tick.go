package main

import (
	"fmt"
	"sync"
	"time"
)

var s string
var by []byte
var wg sync.WaitGroup


func main()  {
	goodsCh := make(chan string) // 存储数据
	orderCh := make(chan string) // 存储数据

	done := make(chan bool)

	go goods(goodsCh)
	go order(orderCh) //

	go doneFuc(done)

	//select里加break可直接跳出for循环
	PRODUCT_CHECK:
	//因为select每次只执行一个case,加for循环可以多次检查
	for {
		// 检查 =》 选择可以运行的通道进行执行，每次只执行其中的一个 （非阻塞的）
		select {
		case order := <-orderCh:
			//这里处理order的业务逻辑
			fmt.Println("order", order)
		case goods := <-goodsCh:
			//这里处理goods的业务逻辑
			fmt.Println("goods", goods)
		case <-done:
			//return // 结束当前的程序
			break PRODUCT_CHECK
		case <-time.After(2e9):
			fmt.Println("超时")
			break PRODUCT_CHECK
			//default:
			//	fmt.Println("default")
		}
	}
	fmt.Println("kk")

	//这是一个定时器示例，每秒执行curl方法
	ticker := time.Tick(1e9)
	tickN  := 0
	for {
		<-ticker
		go curl()
		tickN ++
		if tickN >= 3 {
			break
		}
	}

}

func curl() {
	fmt.Println("请求某一个接口")
}

func goods(ch chan<- string) {
	time.Sleep(1e9)
	ch <- "goods"
}
func order(ch chan<- string) {
	time.Sleep(3e9)
	ch <- "order"
}
func doneFuc(done chan bool) {
	time.Sleep(2e9)
	done <- true
}

