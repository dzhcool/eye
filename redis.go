package eye

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

func NewRedisPool(redis_host string, redis_port string, redis_auth string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     1000,
		MaxActive:   10000,
		IdleTimeout: 120 * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			sn := fmt.Sprintf("%s:%s", redis_host, redis_port)
			c, err := redis.Dial("tcp", sn)
			if err != nil {
				return nil, err
			}
			if redis_auth != "" {
				if _, err := c.Do("AUTH", redis_auth); err != nil {
					c.Close()
					return nil, err
				}
			}

			_, err = c.Do("SELECT", 0)

			return c, err
		},
	}
}
