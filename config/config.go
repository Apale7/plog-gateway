package config

import (
	"log"

	"github.com/spf13/viper"
)

var (
	conf Conf
)

func Init() Conf {
	viper.AddConfigPath("../conf")
	conf = GetConfig()
	return conf
}

func GetConfig() Conf {
	viper.SetConfigName("db_conf")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
		panic(err)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal(err)
		panic(err)
	}
	return conf
}
