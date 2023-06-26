package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	//https://gitcode.net/mirrors/panjf2000/ants?utm_source=csdn_github_accelerator
	//需要go module管理 , go get -u github.com/panjf2000/ants/v2
	"github.com/panjf2000/ants/v2"
)

var sum int32

type student struct {
	name string
	age int
}

var p  *ants.Pool
var pF *ants.PoolWithFunc

func myFunc(i interface{}) {
	n := i.(int32)
	atomic.AddInt32(&sum, n)
	fmt.Printf("run with %d\n", n)
}

func demoFunc(i interface{}) {
	time.Sleep(10 * time.Millisecond)
	//n := i.(student)
	n := i.(int)
	fmt.Printf("Hello World! => %d\n",n)
}

func init()  {
	//方案1.采用直接new协程池的方式,然后通过Submit调用协程。注:调的是协程，异步执行
	p,  _  = ants.NewPool(10)
	//方案2.采用new协程池绑定协程的方式,通过Invoke执行可带参数,方法参数必须为接口类型
	pF, _  = ants.NewPoolWithFunc(10,myFunc)
}

func main() {
	//关闭协程池
	//defer ants.Release()
	defer p.Release()
	defer pF.Release()

	runTimes := 10

	// Use the common pool.
	var wg sync.WaitGroup

	for i := 0; i < runTimes; i++ {
		//方案1.Submit到对应方法,注:由于异步执行若不加协程等待(wg.Wait)会出现全部传入10的现象,但等待失去异步意义，所以传参建议使用NewPoolWithFunc
		//wg.Add(1)
		_ = p.Submit(func() {
			fmt.Printf("p.Submit i: %d\n", i)
			demoFunc(i)
			//wg.Done()
		})
		//wg.Wait()

		//方案2.直接执行绑定的方法并传参i,等价于 myFunc(int32(i))
		_ = pF.Invoke(int32(i))
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", p.Running())
	fmt.Printf("finish all tasks.\n")

	boolChane := make(chan bool)
	<-boolChane
}
