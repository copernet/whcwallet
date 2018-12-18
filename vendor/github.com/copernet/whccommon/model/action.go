package model

import (
	"fmt"
	"strconv"
	"time"

	"github.com/copernet/whc.go/rpcclient"
	"github.com/gomodule/redigo/redis"

	"github.com/jinzhu/gorm"
)

func ConnectDatabase(conf *DBOption) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.User, conf.Passwd, conf.Host, conf.Port, conf.Database)

	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.LogMode(conf.Log)
	db.DB().SetMaxIdleConns(conf.MaxIdleConnection)
	db.DB().SetMaxOpenConns(conf.MaxOpenConnection)
	db.DB().SetConnMaxLifetime(time.Duration(conf.MaxConnectionLifeTime) * time.Minute)

	return
}

func ConnectRedis(option *RedisOption) (*RedisStorage, error) {
	r := &RedisStorage{
		config: option,
		Pool: &redis.Pool{
			MaxIdle:     option.MaxIdleConns,
			MaxActive:   option.MaxOpenConns,
			IdleTimeout: time.Duration(option.Timeout) * time.Second,
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}

	r.Pool.Dial = r.createRedisConnect

	if option.InitialOpenConns > option.MaxIdleConns {
		option.InitialOpenConns = option.MaxIdleConns
	} else if option.InitialOpenConns == 0 {
		option.InitialOpenConns = 1
	}

	return r, r.initRedisConnect()
}

func (r *RedisStorage) createRedisConnect() (redis.Conn, error) {
	return redis.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", r.config.Host, r.config.Port),
		redis.DialDatabase(r.config.DbNum),
		redis.DialPassword(r.config.Passwd),
	)
}

func (r *RedisStorage) initRedisConnect() error {
	cons := make([]redis.Conn, r.config.InitialOpenConns)
	defer func() {
		for _, c := range cons {
			if c != nil {
				c.Close()
			}
		}
	}()

	for i := 0; i < r.config.InitialOpenConns; i++ {
		cons[i] = r.Pool.Get()
		if _, err := cons[i].Do("PING"); err != nil {
			return err
		}
	}

	return nil
}

type RedisStorage struct {
	Pool   *redis.Pool
	config *RedisOption
}

var client *rpcclient.Client

func ConnRpc(conf *RPCOption) *rpcclient.Client {
	if client != nil {
		return client
	}

	// rpc client instance
	connCfg := &rpcclient.ConnConfig{
		Host:         conf.Host + ":" + strconv.Itoa(conf.Port),
		User:         conf.User,
		Pass:         conf.Passwd,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	c, err := rpcclient.New(connCfg, nil)
	if err != nil {
		panic(err)
	}

	client = c
	return c
}
