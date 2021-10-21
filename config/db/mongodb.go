package db

import (
	"context"
	"log"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mgoCli *mongo.Client
	err    error
)

type conf struct {
	MongoDb mongodbConf `json:"mongodb"`
}

type mongodbConf struct {
	Uri string `json:"uri"`
	Db  string `json:"db"`
}

func init() {
	mongodbConf := getMongoDbConf()
	clientOptions := options.Client().ApplyURI(mongodbConf.Uri)
	mgoCli, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = mgoCli.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getMongoDbConf() mongodbConf {
	viper.SetConfigName("db_conf")
	viper.AddConfigPath("./conf")
	if err = viper.ReadInConfig(); err != nil {
		panic(err)
	}
	var dbconf conf
	if err = viper.Unmarshal(&dbconf); err != nil {
		panic(err)
	}
	return dbconf.MongoDb
}
