package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

var (
	accessKey                 = "your_access_key"
	secretKey                 = "your_secret_key"
	timestampToleranceSeconds = 300
)

func main() {
	r := gin.Default()

	// 使用中间件进行认证
	r.Use(authMiddleware())

	r.POST("/your-api-endpoint", handleAPIRequest)

	r.Run(":8080")
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		timestamp, err := strconv.ParseInt(c.GetHeader("Timestamp"), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid timestamp"})
			c.Abort()
			return
		}

		nonce := c.GetHeader("Nonce")
		clientSign := c.GetHeader("Sign")

		if !isValidTimestamp(timestamp) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Timestamp expired"})
			c.Abort()
			return
		}

		if !isValidNonce(nonce) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Duplicate request"})
			c.Abort()
			return
		}

		serverSign := calculateSign(timestamp, nonce)

		if clientSign != serverSign {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func handleAPIRequest(c *gin.Context) {
	// 处理请求逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Request authenticated and processed successfully"})
}

func isValidTimestamp(timestamp int64) bool {
	now := time.Now().Unix()
	return timestamp >= now-int64(timestampToleranceSeconds) && timestamp <= now+int64(timestampToleranceSeconds)
}

func isValidNonce(nonce string) bool {
	// 检查请求标识是否唯一
	// 这里简单演示直接通过
	return true
}

func calculateSign(timestamp int64, nonce string) string {
	message := strconv.FormatInt(timestamp, 10) + nonce + secretKey
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
