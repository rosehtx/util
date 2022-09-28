package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Student struct {
	Class string
	School string
	Love string
	StudentBase
}

type StudentBase struct {
	Name string `json:"Name"`
	Age  int 	`json:"Age"`
	Sex  int 	`json:"Sex"`
}

func main()  {
	path   := "http://192.168.44.127/test/test.php"
	//method := "GET"

	paramsDict := map[string]interface{}{
		"Name": "test rose",
		"Age":  1,
	}
	args := make([]string, 0)
	for k, v := range paramsDict {
		args = append(args, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(fmt.Sprintf("%v", v))))
	}

	var response *http.Response
	var body []byte

	//方案一直接用map操作，请求application/json
	reqMap,_      := json.Marshal(paramsDict)
	response,_     = http.Post(path,"application/json",bytes.NewBuffer(reqMap))
	fmt.Println("=====方案一=====")
	body, _        = ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	u2 	 := new(StudentBase)
	err2 := json.Unmarshal(body, &u2)
	if err2 != nil {
		return
	}
	fmt.Println(*u2)
	fmt.Println("=====方案一=====\n")

	//方案二用结构体，请求application/json
	reqJson,_     := json.Marshal(StudentBase{
		"this is a json struct",
		2,
		2,
	})
	response,_     = http.Post(path,"application/json",bytes.NewBuffer(reqJson))
	fmt.Println("=====方案二=====")
	body, _        = ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	_ = json.Unmarshal(body, &u2)
	fmt.Println(*u2)
	fmt.Println("=====方案二=====\n")

	//方案三，请求x-www-form-urlencoded
	args = append(args,"Sex=1")
	response,_     = http.Post(path,"application/x-www-form-urlencoded",strings.NewReader(strings.Join(args,"&")))
	fmt.Println("=====方案三=====")
	body, _        = ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	_ = json.Unmarshal(body, &u2)
	fmt.Println(*u2)
	fmt.Println("=====方案三=====\n")

	response.Body.Close()

	return
}
