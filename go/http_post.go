package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	ak = "fsferioldfsdj34r0j"
	sk = "ksdfs490fjdslfjlsd"
)

func main() {
	path   := "http://192.168.69.127/test/test.php"
	method := "POST"

	paramsDict := map[string]interface{}{
		"app_id":       1234,
		"access_token": "q3fafa33sHFU+V9h32h0v8weVEH/04hgsrHFHOHNNQOBC9fnwejasubw==",
		"ts":           1555912969,
		"server_id":    "123",
	}

	// 拼接 application/x-www-form-urlencoded content-type 的参数
	args := make([]string, 0)
	for k, v := range paramsDict {
		args = append(args, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(fmt.Sprintf("%v", v))))
	}
	formUrlencodedArgs := strings.Join(args, "&")

	client := &http.Client{}
	req, err := http.NewRequest(method, path, strings.NewReader(formUrlencodedArgs))
	if err != nil {
		fmt.Println(err)
	}

	// 认证字符串，参与签名的数据必须要和放入 http body 中的一致
	signature := Sign(ak, sk, []byte(formUrlencodedArgs))

	// 放入 HTTP 请求 Header 中
	req.Header.Set("Agw-Auth", signature)

	// 登录验证接口要求的 content-type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}

func Sign(ak string, sk string, data []byte) string {
	expiration := 1800 // 过期时间 单位 s

	// 认证字符串 auth-v2/$access_key/$timestamp/$expiretime
	signKeyInfo := fmt.Sprintf("auth-v2/%s/%d/%d", ak, time.Now().Unix(), expiration)

	// 认证字符串生成的 sign_key
	signKey := sha256HMAC([]byte(sk), []byte(signKeyInfo))

	// 结合请求参数，生成 HTTP 请求参数的摘要 signature
	signature := sha256HMAC(signKey, data)

	// 最终拼接出来的认证字符串 auth-v2/$access_key/$timestamp/$expiretime/$signature
	return fmt.Sprintf("%v/%v", signKeyInfo, string(signature))
}

func sha256HMAC(key []byte, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return []byte(fmt.Sprintf("%x", mac.Sum(nil)))
}