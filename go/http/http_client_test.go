package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHttpClient(t *testing.T) {
	resp, _ := http.Get("http://127.0.0.1:8888/test")
	fmt.Println("resp", resp)

	var data [1024]byte // [:]  data := make([]byte, 1024)
	n, _ := resp.Body.Read(data[:])
	fmt.Println("信息 : ", string(data[:n]))
}
