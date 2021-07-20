package internal

import (
	"strings"

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

	WorkAsPublicService = "WorkAsPublicService"
	Debug               = "IsDebugging"

	DbConnectionString = "Db.ConnectionString"
	DbUser             = "Db.User"
	DbPassword         = "Db.Password"
	DbName             = "Db.Name"
)

func NewConfiguration() (Configuration, error) {
	conf := viper.New()
	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	conf.SetEnvPrefix("bot")
	conf.AutomaticEnv()
	_ = pflag.String("conf-path", "", "Path to configuration file")
	pflag.Parse()
	_ = conf.BindPFlag(ConfPath, pflag.Lookup("conf-path"))
	confPath := conf.GetString(ConfPath)
	conf.SetDefault(Token, "")
	conf.SetDefault(ChatId, "")
	conf.SetDefault(ServerUrl, "::")
	conf.SetDefault(ServerCertPath, "")
	conf.SetDefault(ServerKeyPath, "")
	conf.SetDefault(SecretKey, "")
	conf.SetDefault(NoAuth, true)
	conf.SetDefault(WorkAsPublicService, false)
	conf.SetDefault(DbConnectionString, "")
	conf.SetDefault(DbUser, "")
	conf.SetDefault(DbPassword, "")
	conf.SetDefault(DbName, "")

	if confPath != "" {
		conf.SetConfigFile(confPath)
		err := conf.ReadInConfig()
		return conf, err
	}
	return conf, nil
}

var _ Configuration = (*viper.Viper)(nil)
