package main

import (
	"fmt"
	"net/http"
)

func main() {
	// http.Handle("/foo", fooHandler)
	fmt.Println("启动http://127.0.0.1:8888")
	// http://127.0.0.1:8888/bar
	http.HandleFunc("/test", fooHandler) // route
	// http.ListenAndServeTLS("") // https
	http.ListenAndServe(":8888", nil) // http
}

// 要求传递的方法的类型  func(http.ResponseWriter, *http.Request)
// fooHandler(w http.ResponseWriter, r *http.Request)
func fooHandler(w http.ResponseWriter, r *http.Request) {
	// 处理逻辑事项的方法
	fmt.Println("处理逻辑")
	fmt.Println("得到连接", r.RemoteAddr)
	fmt.Println("url", r.URL.Path)
	fmt.Println("method", r.Method)

	//获取客户端的数据，一般用body获取,get请求用下面的Form获取
	var data [1024]byte
	n , _ := r.Body.Read(data[:])
	fmt.Println("body", string(data[:n]))

	//获取post(x-www-form-urlencode请求),get请求数据
	//r.ParseForm()  //预先ParseForm处理一下，然后Form获取
	//fmt.Println("post get info",r.Form)

	w.Write([]byte("zhe shi test 666"))
}
