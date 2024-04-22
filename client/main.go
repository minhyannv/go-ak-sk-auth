package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	accessKey = "your_access_key"
	secretKey = "your_secret_key1"
)

func main() {
	// 模拟发送请求
	sendRequest()
}

func sendRequest() {
	// 生成当前时间戳
	timestamp := time.Now().Unix()

	// 生成随机的请求唯一标识
	nonce := generateNonce()

	// 计算签名
	sign := calculateSign(timestamp, nonce)

	// 构造请求头
	headers := map[string]string{
		"Timestamp": strconv.FormatInt(timestamp, 10),
		"Nonce":     nonce,
		"Sign":      sign,
	}

	// 发送请求
	// 注意替换为实际的请求地址和参数
	req, err := http.NewRequest("POST", "http://localhost:8080/your-api-endpoint", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	// 这里简单演示直接打印响应
	fmt.Println("Response:", resp.Status)
}

func generateNonce() string {
	// 生成随机的请求标识
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf)
}

func calculateSign(timestamp int64, nonce string) string {
	// 计算签名，使用 HMAC-SHA256 算法
	message := fmt.Sprintf("%d%s%s", timestamp, nonce, secretKey)
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
