package config

import (
	"fmt"
	"time"
)

type Database struct {
	Mysql *Mysql
	Redis *Redis
}

type Mysql struct {
	Addr     string `mapstructure:"addr" yaml:"addr"`
	Port     string `mapstructure:"port" yaml:"port"`
	Db       string `mapstructure:"db" yaml:"db"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
	Charset  string `mapstructure:"charset" yaml:"charset"`

	ConnMaxIdleTime string `mapstructure:"connMaxIdleTime" yaml:"ConnMaxIdleTime"`
	ConnMaxLifeTime string `mapstructure:"connMaxLifeTime" yaml:"connMaxLifeTime"`
	MaxIdleConns    int    `mapstructure:"maxIdleConns" yaml:"manIdleConns"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns" yaml:"maxOpenConns"`
}

func (m *Mysql) GetDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		m.Username,
		m.Password,
		m.Addr,
		m.Port,
		m.Db,
		m.Charset)
}

func (m *Mysql) GetConnMaxIdleTime() time.Duration {
	t, _ := time.ParseDuration(m.ConnMaxIdleTime)
	return t
}

func (m *Mysql) GetConnMaxLifeTime() time.Duration {
	t, _ := time.ParseDuration(m.ConnMaxLifeTime)
	return t
}

type Redis struct {
	Addr     string `mapstructure:"addr" yaml:"addr"`
	Port     string `mapstructure:"port" yaml:"port"`
	Password string `mapstructure:"password" yaml:"password"`
	Db       int    `mapstructure:"db" yaml:"db"`
	PoolSize int    `mapstructure:"poolSize" yaml:"poolSize"`
	Network  string `mapstructure:"network" yaml:"network"`
}
