package conf

import (
	"gitlab-tg-bot/internal/interfaces"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	ConfPath       = "ConfPath"
	Token          = "Token"
	ChatId         = "ChatId"
	ServerUrl      = "ServerUrl"
	ServerCertPath = "ServerCertPath"
	ServerKeyPath  = "ServerKeyPath"
	SecretKey      = "SecretKey"
	NoAuth         = "NoAuth"

	Debug      = "IsDebugging"
	WebHookUrl = "WebHookUrl"

	DbConnectionString = "Db.ConnectionString"
	DbPort             = "Db.Port"
	DbUser             = "Db.User"
	DbPassword         = "Db.Password"
	DbName             = "Db.Name"
)

func NewConfiguration() interfaces.Configuration {
	conf := viper.New()

	pflag.Parse()

	_ = conf.BindPFlag(ConfPath, pflag.Lookup("conf-path"))
	confPath := conf.GetString(ConfPath)

	if confPath != "" {
		conf.SetConfigFile(confPath)
		err := conf.ReadInConfig()

		if err != nil {
			panic(err)
		}
	}

	return conf
}

var _ interfaces.Configuration = (*viper.Viper)(nil)
