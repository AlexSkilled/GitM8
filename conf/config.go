package conf

import (
	"strings"

	"gitlab-tg-bot/internal/interfaces"

	"github.com/sirupsen/logrus"

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

	WebHookUrl = "WebHookUrl"

	DbConnectionString = "Db.ConnectionString"
	DbPort             = "Db.Port"
	DbUser             = "Db.User"
	DbPassword         = "Db.Password"
	DbName             = "Db.Name"

	DefaultLanguage = "DefaultLanguage"
)

func NewConfiguration() interfaces.Configuration {
	conf := viper.New()

	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	conf.SetEnvPrefix("TGBOT")
	conf.AutomaticEnv()

	_ = pflag.String("conf-path", "", "Путь к файлу конфигурации")

	pflag.Parse()

	_ = conf.BindPFlag(ConfPath, pflag.Lookup("conf-path"))
	confPath := conf.GetString(ConfPath)

	if confPath != "" {
		logrus.Println(confPath)
		conf.SetConfigFile(confPath)
		err := conf.ReadInConfig()

		if err != nil {
			panic(err)
		}
	} else {
		logrus.Println("Путь к конфигу пуст!")
	}

	return conf
}

var _ interfaces.Configuration = (*viper.Viper)(nil)
