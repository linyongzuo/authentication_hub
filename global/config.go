package global

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host         string
	Port         int
	UserName     string
	Password     string
	DatabaseName string
}
type RedisConfig struct {
	Host     string
	Port     int
	UserName string
	Password string
	Db       int
	PoolSize int
}
type ServerConfig struct {
	Port int
}
type Configs struct {
	DatabaseCfg DatabaseConfig
	ServerCfg   ServerConfig
	RedisCfg    RedisConfig
}

var (
	Cfg    *Configs
	DbPool *gorm.DB
	Rdb    *redis.Client
)
var Upgrader = websocket.Upgrader{ReadBufferSize: 1024,
	WriteBufferSize: 1024}

func init() {
	Cfg = &Configs{}
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./conf")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(Cfg)
	if err != nil {
		panic(err)
	}
	logrus.Info("配置文件:", Cfg)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=UTC&interpolateParams=true",
		Cfg.DatabaseCfg.UserName,
		Cfg.DatabaseCfg.Password,
		Cfg.DatabaseCfg.Host,
		Cfg.DatabaseCfg.Port,
		Cfg.DatabaseCfg.DatabaseName,
	)
	DbPool, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         uint(256),
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		panic(err)
	}
	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Cfg.RedisCfg.Host, Cfg.RedisCfg.Port),
		Password: Cfg.RedisCfg.Password, // 密码
		DB:       Cfg.RedisCfg.Db,       // 数据库
		PoolSize: Cfg.RedisCfg.PoolSize, // 连接池大小
	})
}
