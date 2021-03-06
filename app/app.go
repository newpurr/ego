package app

import (
	"github.com/ebar-go/ego/component/mns"
	"github.com/ebar-go/ego/config"
	"github.com/ebar-go/ws"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/dig"
	"net"
	"net/http"
	"time"
)

var (
	Container = NewContainer()
)

// NewContainer return an empty container
func NewContainer() *dig.Container {
	return dig.New()
}

// WebSocket return ws manager
func WebSocket() (manager ws.Manager) {
	if err := Container.Invoke(func(m ws.Manager) {
		manager = m
	}); err != nil {
		manager = ws.New()
		_ = Container.Provide(func() ws.Manager{
			return manager
		})
	}
	return
}

// Redis get redis connection
func Redis() (connection *redis.Client) {
	_ = Container.Invoke(func(conn *redis.Client) {
		connection = conn
	})
	return
}

// Mysql return mysql connection
func Mysql() (connection *gorm.DB) {
	_ = Container.Invoke(func(conn *gorm.DB) {
		connection = conn
	})
	return
}

// Mns return ali yun mns client
func Mns() (client mns.Client) {
	if err :=  Container.Invoke(func(cli mns.Client) {
		client = cli
	}); err != nil {
		mnsConfig := config.Mns()
		client = mns.NewClient(
			mnsConfig.Url,
			mnsConfig.AccessKeyId,
			mnsConfig.AccessKeySecret,
			)
		_ = Container.Provide(func() (mns.Client) {
			return client
		})
	}
	return
}

// Http client
func Http() (client *http.Client)  {
	if err :=  Container.Invoke(func(cli *http.Client) {
		client = cli
	}); err != nil {
		client = &http.Client{
			Transport: &http.Transport{ // 配置连接池
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				IdleConnTimeout: time.Duration(config.Server().HttpRequestTimeOut) * time.Second,
			},
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       time.Duration(config.Server().HttpRequestTimeOut) * time.Second,
		}
	}
	return
}