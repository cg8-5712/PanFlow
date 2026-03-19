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
	AdminPassword    string `mapstructure:"admin_password"`
	ParsePassword    string `mapstructure:"parse_password"`
	JWTSecret        string `mapstructure:"jwt_secret"`
	JWTExpireHours   int    `mapstructure:"jwt_expire_hours"`
	Debug            bool   `mapstructure:"debug"`
	SaveHistoriesDay int    `mapstructure:"save_histories_day"`

	// 解析
	GuestUserAgent string `mapstructure:"guest_user_agent"`

	// 全局代理
	ProxyHTTP string `mapstructure:"proxy_http"`

	// 邮件
	MailSwitch      bool   `mapstructure:"mail_switch"`
	MailHost        string `mapstructure:"mail_host"`
	MailPort        int    `mapstructure:"mail_port"`
	MailUsername    string `mapstructure:"mail_username"`
	MailPassword    string `mapstructure:"mail_password"`
	MailFromAddress string `mapstructure:"mail_from_address"`
	MailFromName    string `mapstructure:"mail_from_name"`
	MailToAddress   string `mapstructure:"mail_to_address"`
	MailToName      string `mapstructure:"mail_to_name"`
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
	viper.SetDefault("database.host", "127.0.0.1")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.name", "panflow")
	viper.SetDefault("redis.host", "127.0.0.1")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("panflow.jwt_secret", "panflow-change-me")
	viper.SetDefault("panflow.jwt_expire_hours", 24)
	viper.SetDefault("panflow.save_histories_day", 7)
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

	// 保证 save_histories_day 最小为 7
	if cfg.Panflow.SaveHistoriesDay < 7 {
		cfg.Panflow.SaveHistoriesDay = 7
	}

	return &cfg, nil
}

// SaveConfig writes key/value pairs back to the .env-style config file.
// PanFlow uses YAML config instead of .env; callers update the viper store
// and then call this to persist.
func UpdateAndSave(updates map[string]interface{}) error {
	for k, v := range updates {
		viper.Set(k, v)
	}
	return viper.WriteConfig()
}
