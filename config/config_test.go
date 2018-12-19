package config

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/copernet/whccommon/model"
	. "github.com/smartystreets/goconvey/convey"
)

var confData = []byte(`
go_version: 1.10.0
version: 1.0.0

testnet: true

private:
  server_secret: 1D1E6531B32D52D
  session_secret: 970D65CC268C9CFFFA4CA929

db:
  user: root
  passwd: 3A836BF12C34D4FEA4C600151A016
  host: 127.0.0.1
  port: 3306
  database: wallet
  maxopenconnection: 50
  maxidleconnection: 30
  maxconnectionlifetime: 10
  sql_log: true
redis:
  user: root
  passwd: A2BF528202F71857ED5CB86DC1
  host: 127.0.0.1
  port: 6379
  db_num: 1
  timeout: 5
  max_idle_conns: 100
  max_open_conns: 300
  initial_open_conns: 10
log:
  filename: app.log
  level: debug
tx:
  mini_output: "0.00000546"
electron:
  host: 127.0.0.1
  port: 50002
`)

func TestInitConfig(t *testing.T) {
	Convey("Given config file", t, func() {
		filename := fmt.Sprintf("conf_test%04d.yml", rand.Intn(9999))
		os.Setenv(ConfEnv, filename)
		path, err := filepath.Abs("./")
		if err != nil {
			panic(err)
		}
		correctPath := filepath.Join(path, filename)
		os.Setenv(ConfTestEnv, correctPath)

		ioutil.WriteFile(filename, confData, 0664)

		Convey("When init configuration", func() {
			config := GetConf()

			Convey("Configuration should resemble default configuration", func() {
				expected := &configuration{}
				expected.GoVersion = "1.10.0"
				expected.Version = "1.0.0"
				expected.TestNet = true
				expected.Private.ServerSecret = "1D1E6531B32D52D"
				expected.Private.SessionSecret = "970D65CC268C9CFFFA4CA929"

				var db model.DBOption
				db.User = "root"
				db.Passwd = "3A836BF12C34D4FEA4C600151A016"
				db.Host = "127.0.0.1"
				db.Port = 3306
				db.Database = "wallet"
				db.MaxOpenConnection = 50
				db.MaxIdleConnection = 30
				db.MaxConnectionLifeTime = 10
				db.Log = true
				expected.DB = &db

				var redis model.RedisOption
				redis.Passwd = "A2BF528202F71857ED5CB86DC1"
				redis.Host = "127.0.0.1"
				redis.Port = 6379
				redis.DbNum = 1
				redis.Timeout = 5
				redis.MaxIdleConns = 100
				redis.MaxOpenConns = 300
				redis.InitialOpenConns = 10
				expected.Redis = &redis

				var log model.LogOption
				log.Filename = "app.log"
				log.Level = "debug"
				expected.Log = &log

				expected.Tx.MiniOutput = "0.00000546"
				expected.Electron.Host = "127.0.0.1"
				expected.Electron.Port = "50002"

				So(config, ShouldResemble, expected)
			})
		})

		Reset(func() {
			os.Unsetenv(ConfEnv)
			os.Remove(filename)
		})
	})
}
