package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Type string `mapstructure:"type"`
		URI  string `mapstructure:"uri"`
	} `mapstructure:"database"`

	Compression struct {
		Algorithm string `mapstructure:"algorithm"`
		Level     int    `mapstructure:"level"`
	} `mapstructure:"compression"`

	Gather struct {
		XDGConfigHome   string `mapstructure:"xdg_config_home"`
		XDGConfigDirs   string `mapstructure:"xdg_config_dirs"`
		SystemConfigDir string `mapstructure:"system_config_dir"`
	} `mapstructure:"gather"`

	Logging struct {
		Level  string `mapstructure:"level"`
		Format string `mapstructure:"format"`
	} `mapstructure:"logging"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)
	v.SetEnvPrefix("SLAPPER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate config
	if cfg.Database.Type == "" {
		return nil, fmt.Errorf("database.type is required")
	}
	if cfg.Database.URI == "" {
		return nil, fmt.Errorf("database.uri is required")
	}

	// Set defaults if empty
	if cfg.Gather.SystemConfigDir == "" {
		cfg.Gather.SystemConfigDir = "/etc"
	}

	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "info"
	}
	if cfg.Logging.Format == "" {
		cfg.Logging.Format = "text"
	}

	return &cfg, nil
}
