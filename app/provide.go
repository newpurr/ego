package app

import (
	"github.com/ebar-go/ego/config"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/utils"
	"github.com/ebar-go/event"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"sync"
	"time"
)

const (
	// mysql connect event
	MySqlConnectEvent = "MYSQL_CONNECT_EVENT"

	// redis connect event
	RedisConnectEvent = "REDIS_CONNECT_EVENT"
)

var initDBOnce, initRedisOnce *sync.Once
func init() {
	initDBOnce = new(sync.Once)
	initRedisOnce = new(sync.Once)

	event.DefaultDispatcher().Register(MySqlConnectEvent, event.Listener{
		Handle: func(ev event.Event) {
			initDBOnce.Do(func() {
				utils.FatalError("ConnectDatabase", connectDatabase())
			})
		},
	})

	event.DefaultDispatcher().Register(RedisConnectEvent, event.Listener{
		Handle: func(ev event.Event) {
			initRedisOnce.Do(func() {
				utils.FatalError("ConnectRedis", connectRedis())
			})
		},
	})

}

// connectRedis
func connectRedis() error {
	connection := redis.NewClient(config.Redis().Options())
	_, err := connection.Ping().Result()
	if err != nil {
		return errors.RedisConnectFailed("%s", err.Error())
	}

	return Container.Provide(func() *redis.Client {
		return connection
	})
}

// connectDatabase
func connectDatabase() error {
	options := config.Mysql()
	connection, err := gorm.Open("mysql", options.Dsn())
	if err != nil {
		return errors.MysqlConnectFailed("%s", err.Error())
	}

	// set log mod
	connection.LogMode(options.LogMode)
	// set pool config
	connection.DB().SetMaxIdleConns(options.MaxIdleConnections)
	connection.DB().SetMaxOpenConns(options.MaxOpenConnections)
	connection.DB().SetConnMaxLifetime(time.Duration(options.MaxLifeTime) * time.Second)

	return Container.Provide(func() (*gorm.DB, error) {
		return connection, nil
	})
}
