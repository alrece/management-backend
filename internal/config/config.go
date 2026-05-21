package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// C 全局配置实例
var C Config

type Config struct {
	Server        ServerConfig        `mapstructure:"server"`
	MySQL         MySQLConfig         `mapstructure:"mysql"`
	Redis         RedisConfig         `mapstructure:"redis"`
	JWT           JWTConfig           `mapstructure:"jwt"`
	Log           LogConfig           `mapstructure:"log"`
	Snowflake     SnowflakeConfig     `mapstructure:"snowflake"`
	Tenant        TenantConfig        `mapstructure:"tenant"`
	CORS          CORSConfig          `mapstructure:"cors"`
	OSS           OSSConfig           `mapstructure:"oss"`
	LoginSecurity LoginSecurityConfig `mapstructure:"login-security"`
	Consul        ConsulConfig        `mapstructure:"consul"`
	GRPC          GRPCConfig          `mapstructure:"grpc"`
	Telemetry     TelemetryConfig     `mapstructure:"telemetry"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"max-idle-conns"`
	MaxOpenConns int    `mapstructure:"max-open-conns"`
	LogLevel     string `mapstructure:"log-level"`
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
	Sentinel RedisSentinelConfig `mapstructure:"sentinel"`
}

type RedisSentinelConfig struct {
	Enabled    bool     `mapstructure:"enabled"`
	MasterName string   `mapstructure:"master-name"`
	Addrs      []string `mapstructure:"addrs"`
	Password   string   `mapstructure:"password"`
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type JWTConfig struct {
	Secret        string `mapstructure:"secret"`
	AccessExpire  int64  `mapstructure:"access-expire"`
	RefreshExpire int64  `mapstructure:"refresh-expire"`
	Issuer        string `mapstructure:"issuer"`
}

type LogConfig struct {
	Level  string  `mapstructure:"level"`
	Format string  `mapstructure:"format"`
	Output string  `mapstructure:"output"`
	File   LogFile `mapstructure:"file"`
}

type LogFile struct {
	Enable     bool `mapstructure:"enable"`
	MaxSize    int  `mapstructure:"max-size"`
	MaxBackups int  `mapstructure:"max-backups"`
	MaxAge     int  `mapstructure:"max-age"`
}

type SnowflakeConfig struct {
	Node int64 `mapstructure:"node"`
}

type TenantConfig struct {
	Enable            bool     `mapstructure:"enable"`
	IgnoreTables      []string `mapstructure:"ignore-tables"`
	MaxPoolSize       int      `mapstructure:"max-pool-size"`
	IdleTimeout       int      `mapstructure:"idle-timeout"`
	MaxConnsPerTenant int      `mapstructure:"max-conns-per-tenant"`
	DBNamePrefix      string   `mapstructure:"db-name-prefix"`
}

func (t TenantConfig) IsIgnoreTable(table string) bool {
	for _, v := range t.IgnoreTables {
		if v == table {
			return true
		}
	}
	return false
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed-origins"`
}

type OSSConfig struct {
	Type  string      `mapstructure:"type"`
	Minio MinioConfig `mapstructure:"minio"`
}

type MinioConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access-key"`
	SecretKey string `mapstructure:"secret-key"`
	Bucket    string `mapstructure:"bucket"`
	UseSSL    bool   `mapstructure:"use-ssl"`
}

type LoginSecurityConfig struct {
	MaxFailCount    int `mapstructure:"max-fail-count"`
	RateLimitPerIP  int `mapstructure:"rate-limit-per-ip"`
	RateLimitPerAcct int `mapstructure:"rate-limit-per-acct"`
	PasswordMinLen  int `mapstructure:"password-min-len"`
}

type ConsulConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	Address     string `mapstructure:"address"`
	ServiceAddr string `mapstructure:"service-addr"`
}

type GRPCConfig struct {
	Port int `mapstructure:"port"`
}

type TelemetryConfig struct {
	Enabled   bool    `mapstructure:"enabled"`
	Endpoint  string  `mapstructure:"endpoint"`
	SampleRate float64 `mapstructure:"sample-rate"`
}

// Validate 启动时校验敏感配置
func (c *Config) Validate() error {
	if c.JWT.Secret == "" || c.JWT.Secret == "your-secret-key" {
		return errors.New("JWT 密钥必须通过环境变量 MB_JWT_SECRET 设置非默认值")
	}
	if c.MySQL.Username == "root" && c.MySQL.Password == "123456" {
		return errors.New("MySQL 禁止使用 root/123456，请通过环境变量设置安全密码")
	}
	return nil
}

// Load 加载配置文件，支持 MB_ 前缀的环境变量覆盖
func Load(cfgPath string) error {
	v := viper.New()
	v.SetConfigFile(cfgPath)
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvPrefix("MB")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}
	if err := v.Unmarshal(&C); err != nil {
		return fmt.Errorf("解析配置失败: %w", err)
	}
	return nil
}
