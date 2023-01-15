package config

type Logger struct {
	LogLevel   string `mapstructure:"logLevel" yaml:"logLevel"`
	MaxSize    int    `mapstructure:"maxSize" yaml:"maxSize"`
	MaxAge     int    `mapstructure:"maxAge" yaml:"maxAge"`
	MaxBackups int    `mapstructure:"maxBackups" yaml:"maxBackups"`
	Compress   bool   `mapstructure:"compress" yaml:"compress"`
}
