package system

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

var SystemConfig *Config

const (
	SERVER_ADDRESS     = "server.address"
	SERVER_PORT        = "server.port"
	APP_NODE           = "app.node"
	ACCESS_TOKEN_TIME  = "auth.accessTokenExpireTime"
	REFRESH_TOEKN_TIME = "auth.refreshTokenExpireTime"
)

type Config struct {
	Db     DbConfig            `mapstructure:"db"`
	Server ServerConfig        `mapstructure:"server"`
	Log    LogConfig           `mapstructure:"log"`
	App    AppConfig           `mapstructure:"app"`
	JwtKey string              `mapstructure:"jwt_key"`
	Auth   AuthConfig          `mapstructure:"auth"`
	Binary BinaryStorageConfig `mapstructure:"bin"`
}

type DbConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	TimeZone string `mapstructure:"timezone"`
	Mirgate  bool   `mapstructure:"migrate"`
}

type ServerConfig struct {
	Address      string   `mapstructure:"address"`
	Port         string   `mapstructure:"port"`
	AllowOrigins []string `mapstructure:"allowOrigins"`
	AllowMethods []string `mapstructure:"allowMethods"`
}

type LogConfig struct {
	Path  string `mapstructure:"path"`
	Level string `mapstructure:"level"`
}

type AppConfig struct {
	Node string `mapstructure:"node"`
}

type AuthConfig struct {
	RefreshTokenExpireTime uint64 `mapstructure:"refreshTokenExpireTime"`
	AccessTokenExpireTime  uint64 `mapstructure:"accessTokenExpireTime"`
}

type BinaryStorageConfig struct {
	ServerAddress string `mapstructure:"serverAddress"`
	Username      string `mapstructure:"username"`
	Password      string `mapstructure:"password"`
	Bucket        string `mapstructure:"bucket"`
}

func InitSystemConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Cannot read config file", err)
	}
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config changed")
		err := viper.Unmarshal(&SystemConfig)
		if err != nil {
			log.Fatal("Cannot read config file", err)
		}
	})
	setDefaultConfig()
	err = viper.Unmarshal(&SystemConfig)
	if err != nil {
		fmt.Println("Cannot read config file", err)
	}
}

func setDefaultConfig() {
	viper.SetDefault(SERVER_PORT, 7777)
	viper.SetDefault(SERVER_ADDRESS, "localhost")
	viper.SetDefault(ACCESS_TOKEN_TIME, 1800000)
	viper.SetDefault(REFRESH_TOEKN_TIME, 2592000000)
	viper.Set(APP_NODE, "1")
}
