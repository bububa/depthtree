package common

import (
	"fmt"
	"github.com/bububa/depthtree"
	"github.com/garyburd/redigo/redis"
	"github.com/ziutek/mymysql/autorc"
	_ "github.com/ziutek/mymysql/thrsafe"
	"strconv"
	"time"
)

type RedisConf struct {
	Host           string
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

type RedisConn struct {
	Master *redis.Pool
	Slave  *redis.Pool
}

// Service struct provide i/o resources for application
type Service struct {
	Db        *autorc.Conn        `json:"-"`
	TreeBase  *depthtree.Database `json:"-"`
	Redis     *RedisConn          `json:"-"`
	redisConf *RedisConf          `json:"-"`
}

func NewService(config Config) *Service {

	mdb := autorc.New("tcp", "", config.MySQL.Host, config.MySQL.User, config.MySQL.Passwd, config.MySQL.DB)
	mdb.Register("set names utf8mb4")

	service := &Service{
		Db:       mdb,
		TreeBase: depthtree.NewDatabase(config.DBPath),
	}

	service.RedisPool(config.Redis.Master, config.Redis.Slave, 10, 120)
	return service
}

func (this *Service) Close() {
	this.CloseRedisPool()
}

func newRedisPool(server string, maxIdle int, idleTime time.Duration) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTime * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			/*if _, err := c.Do("AUTH", password); err != nil {
			    c.Close()
			    return nil, err
			}*/
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func (this *Service) RedisPool(master string, slave string, maxIdle int, idleTime time.Duration) {
	masterPool := newRedisPool(master, maxIdle, idleTime)
	slavePool := newRedisPool(slave, maxIdle, idleTime)
	this.Redis = &RedisConn{
		Master: masterPool,
		Slave:  slavePool,
	}
}

func (this *Service) CloseRedisPool() error {
	if this.Redis == nil {
		return nil
	}
	var err error
	if this.Redis.Master != nil {
		err = this.Redis.Master.Close()
	}
	if this.Redis.Slave != nil {
		err = this.Redis.Slave.Close()
	}
	return err
}

// JsonTime marshal time.Time compatible with PHP json_encode.
type JsonTime time.Time

func (t JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", time.Time(t).Unix())), nil
}

func (t *JsonTime) UnmarshalJSON(s []byte) (err error) {
	sec, err := strconv.ParseInt(string(s), 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(sec, 0)
	return
}

func (t JsonTime) String() string { return time.Time(t).String() }
