package main

import (
	"fmt"
	"sync"
)

var s string
var by []byte
var wg sync.WaitGroup


func main()  {

	ch   := make(chan int)
	done := make(chan bool)

	// 会打印 1~ 3
	go send(1, 4, ch)
	go recv(ch, done) // 因为采用循环 =》从通道中获取信息 =》 后面的done通道无法设置参数

	// 利用通道读取的阻塞来执行上面协程
	<-done
}

func send(start, end int, ch chan<- int) {
	for i := start; i < end; i++ {
		ch <- i
	}
	//由于下面recv方法循环读取通道，需要关闭通道来执行接下来的done通道写入，避免死锁
	close(ch)
}

func recv(in <-chan int, done chan<- bool) {
	for num := range in {
		fmt.Println("num : ", num)
	}
	done <- true
}

