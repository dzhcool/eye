package eye

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

type RedisSvc struct {
	connType string
	host     string
	port     string
	db       int
	auth     string
	pool     *redis.Pool
}

func NewRedis(redis_host string, redis_port string, redis_auth string, redis_db int, enablePool bool) *RedisSvc {
	var redisPool *redis.Pool = nil
	if enablePool == true {
		redisPool = &redis.Pool{
			MaxIdle:     1000,
			MaxActive:   10000,
			IdleTimeout: 240 * time.Second,
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

				_, err = c.Do("SELECT", redis_db)

				return c, err
			},
		}
	}

	return &RedisSvc{
		connType: "tcp",
		host:     redis_host,
		port:     redis_port,
		db:       redis_db,
		auth:     redis_auth,
		pool:     redisPool,
	}
}

func (p *RedisSvc) getRedisPool() redis.Conn {
	if p.pool != nil {
		return p.pool.Get()
	}
	return nil
}

func (p *RedisSvc) ActiveCount() int {
	if p.pool != nil {
		return p.pool.ActiveCount()
	}
	return 0
}

func (p *RedisSvc) DoCmd(cmd string, args ...interface{}) (interface{}, error) {
	var c redis.Conn
	var err error
	if p.pool != nil {
		c = p.getRedisPool()
	} else {
		c, err = redis.Dial("tcp", fmt.Sprintf("%s:%s", p.host, p.port))
		if err != nil {
			return nil, err
		}
	}
	defer c.Close()

	re, err := c.Do(cmd, args...)
	if err != nil {
		return nil, err
	}
	return re, nil
}
