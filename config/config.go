package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bcext/gcash/chaincfg"
	"github.com/copernet/whccommon/model"
	"github.com/spf13/viper"
)

const (
	ConfEnv        = "WHC_CONF"
	ConfTestEnv    = "WHC_TEST_CONF"
	ProjectLastDir = "whcwallet"
	FeeScale       = 0.3
)

var api = map[string]string{
	"testnet": "https://developer-bch-tchain.api.btc.com/appkey-2f7c183e3e9e",
	"mainnet": "https://developer-bch-chain.api.btc.com/appkey-2f7c183e3e9e",
}

var burningAddress = map[string]string{
	"testnet": "bchtest:qqqqqqqqqqqqqqqqqqqqqqqqqqqqqdmwgvnjkt8whc",
	"mainnet": "bitcoincash:qqqqqqqqqqqqqqqqqqqqqqqqqqqqqu08dsyxz98whc",
}

func GetBCHAPI() string {
	conf := GetConf()
	if conf.TestNet {
		return api["testnet"]
	}

	return api["mainnet"]
}

func GetBurningAddress() string {
	conf := GetConf()
	if conf.TestNet {
		return burningAddress["testnet"]
	}

	return burningAddress["mainnet"]
}

func GetChainParam() *chaincfg.Params {
	conf := GetConf()
	if conf.TestNet {
		return &chaincfg.TestNet3Params
	}

	return &chaincfg.MainNetParams
}

var conf *configuration

type configuration struct {
	GoVersion string `mapstructure:"go_version"`
	Version   string `mapstructure:"version"`
	TestNet   bool   `mapstructure:"testnet"`
	Private   struct {
		ServerSecret     string `mapstructure:"server_secret"`
		SessionSecret    string `mapstructure:"session_secret"`
		AesSecret        string `mapstructure:"aes_secret"`
		LoginDifficulty  string `mapstructure:"login_difficulty"`
		RecaptchaPrivate string `mapstructure:"recaptcha_private"`
	}
	DB       *model.DBOption    `mapstructure:"db"`
	WalletDB *model.DBOption    `mapstructure:"walletdb"`
	Redis    *model.RedisOption `mapstructure:"redis"`
	Log      *model.LogOption   `mapstructure:"log"`
	RPC      *model.RPCOption   `mapstructure:"rpc"`
	Tx       struct {
		MiniOutput string `mapstructure:"mini_output"`
	}
	Electron struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}
	Mail struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		UserName string `mapstructure:"user_name"`
		SmtpUser string `mapstructure:"smtp_user"`
		SmtpPass string `mapstructure:"smtp_pass"`
		FromName string `mapstructure:"from_name"`
	}
}

func GetConf() *configuration {
	if conf != nil {
		return conf
	}

	config := &configuration{}
	viper.SetEnvPrefix("whc")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	viper.SetDefault("conf", "./conf.yml")

	// get config file path from environment
	confFile := viper.GetString("conf")

	var realPath string

	// conf.go unit testing
	if viper.GetString("test_conf") != "" {
		realPath = viper.GetString("test_conf")
	} else {
		path, err := filepath.Abs("./")
		if err != nil {
			panic(err)
		}

		lastIndex := strings.Index(path, ProjectLastDir) + len(ProjectLastDir)
		correctPath := path[:lastIndex]
		realPath = filepath.Join(correctPath, confFile)
	}

	// parse config
	file, err := os.Open(realPath)
	if err != nil {
		panic("Open config file error: " + err.Error())
	}
	defer file.Close()

	err = viper.ReadConfig(file)
	if err != nil {
		panic("Read config file error: " + err.Error())
	}

	err = viper.Unmarshal(config)
	if err != nil {
		panic("Parse config file error: " + err.Error())
	}

	// TODO validate configuration
	//helper.Must(nil, config.Validate())

	conf = config
	return config
}
