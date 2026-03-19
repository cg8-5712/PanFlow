package config

import (
	"strings"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host string
	Port string
	Mode string
}

type DatabaseConfig struct {
	Driver   string // mysql | postgres | sqlite
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type LogConfig struct {
	Level string
}

type PanflowConfig struct {
	AdminPassword   string `mapstructure:"admin_password"`
	JWTSecret       string `mapstructure:"jwt_secret"`
	JWTExpireHours  int    `mapstructure:"jwt_expire_hours"`
	JWTRefreshDays  int    `mapstructure:"jwt_refresh_days"`
	Debug          bool   `mapstructure:"debug"`
	GuestUserAgent string `mapstructure:"guest_user_agent"`
	ProxyHTTP      string `mapstructure:"proxy_http"`
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Log      LogConfig
	Panflow   PanflowConfig
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// 默认值
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "release")
	viper.SetDefault("database.driver", "mysql")
	viper.SetDefault("database.host", "127.0.0.1")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.name", "panflow")
	viper.SetDefault("redis.host", "127.0.0.1")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("panflow.jwt_secret", "panflow-change-me")
	viper.SetDefault("panflow.jwt_expire_hours", 2)
	viper.SetDefault("panflow.jwt_refresh_days", 7)
	viper.SetDefault("panflow.guest_user_agent", "netdisk;P2SP;3.0.20.138")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
