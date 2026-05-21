package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var C Config

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	MySQL     MySQLConfig     `mapstructure:"mysql"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Snowflake SnowflakeConfig `mapstructure:"snowflake"`
	Consul    ConsulConfig    `mapstructure:"consul"`
}

type ConsulConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	Address     string `mapstructure:"address"`
	ServiceAddr string `mapstructure:"service-addr"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Charset  string `mapstructure:"charset"`
}

func (m MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		m.Username, m.Password, m.Host, m.Port, m.Database, m.Charset)
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

type SnowflakeConfig struct {
	Node int64 `mapstructure:"node"`
}

func Load(cfgPath string) error {
	v := viper.New()
	v.SetConfigFile(cfgPath)
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvPrefix("JOB")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置失败: %w", err)
	}
	return v.Unmarshal(&C)
}
