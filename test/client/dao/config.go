package dao

import "time"

const (
	redisServerAddr     string        = "adx-v8-cache.redis.rds.aliyuncs.com:6379" // redis服务地址
	redisMaxIdle        int           = 50                                         // 最大空闲连接
	redisIdleTimeout    time.Duration = 240                                        // 最大空闲时间
	redisServerIsAuth   bool          = true                                       // redis是否开启权鉴
	redisServerPassword string        = "adx_v8@Passw0rd"                          // redis密码
)
