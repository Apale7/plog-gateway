package db

import (
	"context"
	"log"
	"plog_gateway/config"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	conf     config.Conf
	mgoCli   *mongo.Client
	redisCli *redis.Client
	err      error
)

const (
	Redis   = "redis"
	MongoDb = "mongo"
)

func DialMongo() *mongo.Client {
	conf = config.Init()
	clientOptions := options.Client().ApplyURI(conf.MongoDb.Uri)
	mgoCli, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	err = mgoCli.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return mgoCli
}

func DialRedis() *redis.Client {
	conf = config.Init()
	redisCli = redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		DB:       conf.Redis.Db,
		Password: conf.Redis.Password,
	})
	_, err = redisCli.Ping().Result()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return redisCli
}
