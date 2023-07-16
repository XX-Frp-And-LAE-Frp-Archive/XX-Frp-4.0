package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func generateToken() string {
	// 生成当前时间戳
	timestamp := time.Now().Unix()
	// 生成随机数
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(1000000) // 随机范围可根据需要更改

	// 将时间戳和随机数拼接成字符串
	tokenString := strconv.FormatInt(timestamp, 10) + strconv.Itoa(random)

	// 计算 MD5 值
	hash := md5.Sum([]byte(tokenString))
	md5Str := hex.EncodeToString(hash[:])

	return md5Str
}

func main() {
	token := generateToken()
	fmt.Println(token)
}
