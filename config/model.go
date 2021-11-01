package config

type Conf struct {
	MongoDb MongodbConf `json:"mongodb"`
	Redis   RedisConf   `json:"redis"`
}

type MongodbConf struct {
	Uri string `json:"uri"`
	Db  string `json:"db"`
}

type RedisConf struct {
	Addr     string `json:"addr"`
	Db       int    `json:"db"`
	Password string `json:"password"`
}
