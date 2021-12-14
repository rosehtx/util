package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var s string
var by []byte
var wg sync.WaitGroup

type User struct {
	Name string
	Age int
}

func init()  {
	s  = "ssss"
	by = []byte("woowowo")
}

var mutex sync.Mutex

// 互斥锁：只适合读和写相差不多的情况下运用
func main() {
	file := "E:\\phpStudy\\PHPTutorial\\WWW\\go\\src\\test\\1.txt"
	data := "ceshi Mutex"
	for i := 0; i < 5; i++ {
		go Oper(file, data)
	}
	time.Sleep(1e9)
}

func Oper(path, data string) {
	// 加锁 => 一定要成对存在
	mutex.Lock()
	d, _ := ReadFile(path)
	if d == "" {
		WriteFile(path, data)
		fmt.Println("写信息成功 : ", data)
	} else {
		fmt.Println("存在信息", d)
	}

	mutex.Unlock() // 释放锁
}
func ReadFile(path string) (string, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("读取json文件失败", err)
		return "", err
	}
	return string(bytes), nil
}
func WriteFile(path, data string) (bool, error) {
	// 打开文件
	// 0666 是文件的写入权限
	outputFile, outputError := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)

	if outputError != nil {
		return false, outputError
	}
	defer outputFile.Close()
	// 创建写的缓冲区
	outputWriter := bufio.NewWriter(outputFile)
	// 写入信息
	outputWriter.WriteString(data)
	// 刷新
	outputWriter.Flush()
	return true, nil
}


