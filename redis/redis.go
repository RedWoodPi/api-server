package redis

import (
    "github.com/garyburd/redigo/redis"
    "time"
    "api/config"
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
            c, err := redis.Dial("tcp", config.Host)
            if err != nil {
                return nil, err
            }
            if _, err := c.Do("AUTH", config.Password); err != nil {
                c.Close()
                return nil, err
            }
            return c, nil
        },
    }
}