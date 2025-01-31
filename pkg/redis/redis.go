package redis

import (
	"github.com/ahmr-bot/ME-Frp/pkg/config"
	"github.com/go-redis/redis"
)

func ConnectRedis() *redis.Client {
	// 初始化配置
	conf := config.GetConfig()

	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Host + ":" + conf.Redis.Port,
		Password: conf.Redis.Password, // no password set
		DB:       conf.Redis.DB,       // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	return client
}
