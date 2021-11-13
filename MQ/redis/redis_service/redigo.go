package redis_service

import (
	"fmt"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

type RedisConfig struct {
	Alias    string `mapstructure:"alias"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	MaxIdles int    `mapstructure:"max_idles"`
	MaxConns int    `mapstructure:"max_conns"`
}

type _rdConf struct {
	db           int
	host         string
	port         int
	password     string
	poolMaxIdle  int
	poolMaxConns int
}

type WRedisPool struct {
	pool *redis.Pool
}

func (p *WRedisPool) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	st := time.Now().UnixNano() / 1e6
	r := p.pool.Get()
	defer r.Close()
	if et := time.Now().UnixNano() / 1e6; et-st >= 30 {
		logrus.Warnf("[redigo] Do: get conn from pool over time: %dms", et-st)
	}
	return r.Do(commandName, args...)
}

func (p *WRedisPool) GetOut() redis.Conn {
	st := time.Now().UnixNano() / 1e6
	r := p.pool.Get()
	if et := time.Now().UnixNano() / 1e6; et-st >= 30 {
		logrus.Warnf("[redigo] Do: get conn from pool over time: %dms", et-st)
	}
	return r
}

var (
	snl   sync.Mutex
	confs map[string]_rdConf
	pools map[string]*WRedisPool
)

func init() {
	confs = make(map[string]_rdConf)
	pools = make(map[string]*WRedisPool)
}

func getConf(sectionName string, db int) (_rdConf, bool) {
	cf, ok := confs[sectionName]
	if ok {
		cf.db = db
	}
	return cf, ok
}

func createPool(cf _rdConf) *WRedisPool {
	// redis新增连接池最大连接数配置
	if cf.poolMaxIdle <= 0 {
		cf.poolMaxIdle = 5 // 默认5个连接
	}
	return &WRedisPool{&redis.Pool{
		MaxIdle:     cf.poolMaxIdle,
		MaxActive:   cf.poolMaxConns,
		Wait:        true,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(`tcp`, fmt.Sprintf(`%s:%d`, cf.host, cf.port), redis.DialPassword(cf.password), redis.DialConnectTimeout(time.Second*10), redis.DialReadTimeout(time.Second*5), redis.DialWriteTimeout(time.Second*3), redis.DialDatabase(cf.db))
		},
	}}
}

// SetConnection 设置redis连接信息
func SetConnection(c RedisConfig) {
	if c.Alias == `` {
		c.Alias = `default`
	}
	if c.MaxIdles == 0 {
		c.MaxIdles = 5
	}
	if c.MaxConns == 0 {
		c.MaxConns = 5
	}
	var sectionName = c.Alias

	confs[sectionName] = _rdConf{
		host:         c.Host,
		port:         c.Port,
		password:     c.Password,
		poolMaxIdle:  c.MaxIdles,
		poolMaxConns: c.MaxConns,
	}
	logrus.Debug("register redis conn pool: " + sectionName)
}

// GetPool 获取一个连接池
func GetPool(sectionName string, db int) (*WRedisPool, error) {
	cacheName := fmt.Sprintf("%s:%d", sectionName, db)
	p, ok := pools[cacheName]
	if ok {
		return p, nil
	}
	snl.Lock()
	defer snl.Unlock()
	p, ok = pools[cacheName]
	if ok {
		return p, nil
	}
	conf, ok := getConf(sectionName, db)
	if !ok {
		return nil, fmt.Errorf("failed to find redis conf %s %d", sectionName, db)
	}
	inst := createPool(conf)
	pools[cacheName] = inst
	return inst, nil
}

// 获取redis连接池
func GetInstance(db int, section ...string) *WRedisPool {
	var sectionName = "default"
	if len(section) > 0 {
		sectionName = section[0]
	}
	p, err := GetPool(sectionName, db)
	if err != nil {
		panic(fmt.Errorf("section not register: %s", sectionName))
	}
	return p
}

func GetPubSubConn(section ...string) (*redis.PubSubConn, error) {
	var sectionName = "default"
	if len(section) > 0 {
		sectionName = section[0]
	}
	cf, ok := getConf(sectionName, 0)
	if !ok {
		return nil, fmt.Errorf("failed to find redis conf %s", sectionName)
	}
	const healthCheckPeriod = time.Minute * 1
	c, err := redis.Dial(`tcp`, fmt.Sprintf(`%s:%d`, cf.host, cf.port), redis.DialPassword(cf.password), redis.DialConnectTimeout(time.Second*10),
		redis.DialReadTimeout(healthCheckPeriod+time.Second*5), redis.DialWriteTimeout(time.Second*3))
	if err != nil {
		return nil, err
	}
	return &redis.PubSubConn{Conn: c}, err
}
