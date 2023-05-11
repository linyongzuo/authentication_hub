package global

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	UserName string
	Password string
}
type ServerConfig struct {
	Port int
}
type Configs struct {
	DatabaseCfg DatabaseConfig
	ServerCfg   ServerConfig
}

var (
	Cfg    *Configs
	DbPool *gorm.DB
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
}
