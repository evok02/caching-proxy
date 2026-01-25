package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Proxy        ProxyConfig
	ReverseProxy ReverseProxyConfig
}

type DBConfig struct {
	Port string
	Name string
}

type ProxyConfig struct {
	Port   string
	Host   string
	Logs   LogsConfig
	DB     DBConfig
	Target string
}

type LogsConfig struct {
	ErrLogsPath  string `mapstructure:"err_log"`
	InfoLogsPath string `mapstructure:"info_log"`
}

type ReverseProxyConfig struct {
	Expose string
}

func Init() (*Config, error) {
	if err := setUpViper(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("proxy.log", &cfg.Proxy.Logs); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("reverse_proxy", &cfg.ReverseProxy); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("proxy.db", &cfg.Proxy.DB); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("proxy", &cfg.Proxy); err != nil {
		return err
	}

	return nil
}

func setUpViper() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	return viper.ReadInConfig()
}
