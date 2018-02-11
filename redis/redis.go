package redis

import (
    "github.com/garyburd/redigo/redis"
    "time"
)

var(
    RedisClient *redis.Pool
)

func init()  {
    RedisClient = &redis.Pool{
        MaxIdle: 20,
        MaxActive: 100,
        IdleTimeout: 120 * time.Second,
        Dial: func() (redis.Conn, error) {
            c, err := redis.Dial("tcp", "r-uf67bbcf0d9ea284.redis.rds.aliyuncs.com:6379")
            if err != nil {
                return nil, err
            }
            return c, nil
        },
    }
}