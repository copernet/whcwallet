package model

import "time"

type DBOption struct {
	User                  string `mapstructure:"user"`
	Passwd                string `mapstructure:"passwd"`
	Host                  string `mapstructure:"host"`
	Port                  int    `mapstructure:"port"`
	Database              string `mapstructure:"database"`
	MaxOpenConnection     int    `mapstructure:"maxopenconnection"`
	MaxIdleConnection     int    `mapstructure:"maxidleconnection"`
	MaxConnectionLifeTime int    `mapstructure:"maxconnectionlifetime"`
	Log                   bool   `mapstructure:"sql_log"`
}

type RedisOption struct {
	Passwd           string `mapstructure:"passwd"`
	Host             string `mapstructure:"host"`
	Port             int    `mapstructure:"port"`
	DbNum            int    `mapstructure:"db_num"`
	Timeout          int    `mapstructure:"timeout"`
	MaxIdleConns     int    `mapstructure:"max_idle_conns"`
	MaxOpenConns     int    `mapstructure:"max_open_conns"`
	InitialOpenConns int    `mapstructure:"initial_open_conns"`
}

type LogOption struct {
	Filename string        `mapstructure:"filename"`
	Level    string        `mapstructure:"level"`
	MaxAge   time.Duration `mapstructure:"maxage"`
}

type RPCOption struct {
	User   string `mapstructure:"user"`
	Passwd string `mapstructure:"passwd"`
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
}
