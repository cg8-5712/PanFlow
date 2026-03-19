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
	// 通用
	AdminPassword    string `mapstructure:"admin_password"`
	ParsePassword    string `mapstructure:"parse_password"`
	JWTSecret        string `mapstructure:"jwt_secret"`
	JWTExpireHours   int    `mapstructure:"jwt_expire_hours"`
	ShowAnnounce     bool   `mapstructure:"show_announce"`
	Announce         string `mapstructure:"announce"`
	CustomScript     string `mapstructure:"custom_script"`
	CustomButton     string `mapstructure:"custom_button"`
	ShowHero         bool   `mapstructure:"show_hero"`
	Name             string `mapstructure:"name"`
	Logo             string `mapstructure:"logo"`
	Debug            bool   `mapstructure:"debug"`
	DisableCheckRand bool   `mapstructure:"disable_check_rand"`
	SaveHistoriesDay int    `mapstructure:"save_histories_day"`

	// 解析
	ParserServer       string `mapstructure:"parser_server"`
	ParserPassword     string `mapstructure:"parser_password"`
	AllowFolder        bool   `mapstructure:"allow_folder"`
	DdddocrServer      string `mapstructure:"ddddocr_server"`
	TokenParseMode     int    `mapstructure:"token_parse_mode"`
	TokenUserAgent     string `mapstructure:"token_user_agent"`
	GuestParseMode     int    `mapstructure:"guest_parse_mode"`
	GuestUserAgent     string `mapstructure:"guest_user_agent"`
	TokenProxyHost     string `mapstructure:"token_proxy_host"`
	TokenProxyPassword string `mapstructure:"token_proxy_password"`
	GuestProxyHost     string `mapstructure:"guest_proxy_host"`
	GuestProxyPassword string `mapstructure:"guest_proxy_password"`
	MoiuToken          string `mapstructure:"moiu_token"`

	// 限制
	MaxOnce                    float64 `mapstructure:"max_once"`
	MinSingleFilesize          float64 `mapstructure:"min_single_filesize"`
	MaxSingleFilesize          float64 `mapstructure:"max_single_filesize"`
	MaxAllFilesize             float64 `mapstructure:"max_all_filesize"`
	MaxDownloadDailyPreAccount float64 `mapstructure:"max_download_daily_pre_account"`
	LimitCN                    bool    `mapstructure:"limit_cn"`
	LimitProv                  bool    `mapstructure:"limit_prov"`
	RemoveLimit                bool    `mapstructure:"remove_limit"`

	// 全局代理
	ProxyEnable bool   `mapstructure:"proxy_enable"`
	ProxyHTTP   string `mapstructure:"proxy_http"`
	ProxyHTTPS  string `mapstructure:"proxy_https"`

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

	// 伪装
	FakeUserAgent   string `mapstructure:"fake_user_agent"`
	FakeWxUserAgent string `mapstructure:"fake_wx_user_agent"`
	FakeCookie      string `mapstructure:"fake_cookie"`
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
	viper.SetDefault("panflow.name", "PanFlow")
	viper.SetDefault("panflow.logo", "/favicon.ico")
	viper.SetDefault("panflow.jwt_secret", "panflow-change-me")
	viper.SetDefault("panflow.jwt_expire_hours", 24)
	viper.SetDefault("panflow.save_histories_day", 7)
	viper.SetDefault("panflow.max_once", 5)
	viper.SetDefault("panflow.max_single_filesize", 53687091200)
	viper.SetDefault("panflow.max_all_filesize", 10737418240)
	viper.SetDefault("panflow.token_user_agent", "netdisk;P2SP;3.0.20.138")
	viper.SetDefault("panflow.guest_user_agent", "netdisk;P2SP;3.0.20.138")
	viper.SetDefault("panflow.token_proxy_password", "panflow")
	viper.SetDefault("panflow.guest_proxy_password", "panflow")
	viper.SetDefault("panflow.ddddocr_server", "https://ddddocr.huankong.top")
	viper.SetDefault("panflow.fake_user_agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	viper.SetDefault("panflow.fake_cookie", "BAIDUID=A4FDFAE43DDBF7E6956B02F6EF715373:FG=1; BAIDUID_BFESS=A4FDFAE43DDBF7E6956B02F6EF715373:FG=1; newlogin=1")

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
