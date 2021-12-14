package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
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

var rwlock sync.RWMutex
var count int

// 读写锁：比较适合于读多、写少的情况
func main() {
	file := "E:\\phpStudy\\PHPTutorial\\WWW\\go\\src\\test\\1.txt"

	for i := 0; i < 10; i++ {
		go func() {
			rwlock.RLock() // 加读锁
			d, _ := ReadFile(file)
			fmt.Println("read data", d)
			rwlock.RUnlock() // 释放读锁
		}()
	}

	for i := 0; i < 3; i++ {
		go func() {
			rwlock.Lock() // 加写锁
			count++
			fmt.Println("写入数据:", strconv.Itoa(count) + ":test;")
			// Oper(file, strconv.Itoa(i))
			WriteFile(file, strconv.Itoa(count) + ":test;")
			rwlock.Unlock() // 释放写锁
		}()
	}

	time.Sleep(1e9)
}

// 先查询再写入数据 =》并不是以读为主
func Oper(path, data string) {
	d, _ := ReadFile(path)
	if d == "" {
		WriteFile(path, data)
		fmt.Println("写信息成功 : ", data)
	} else {
		fmt.Println("存在信息", d)
	}
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