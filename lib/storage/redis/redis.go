package redis

import (
	"crypto/tls"
	redis "github.com/go-redis/redis/v7"
	"github.com/godofcc/go-common/lib/log"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type RedisOptions struct {
	Host                  string   `json:"host" description:"Redis service host address"`
	Port                  int      `json:"port"`
	Addrs                 []string `json:"addrs"`
	Username              string   `json:"username"`
	Password              string   `json:"password"`
	Database              int      `json:"database"`
	MasterName            string   `json:"master-name"`
	MaxIdle               int      `json:"optimisation-max-idle"`
	MaxActive             int      `json:"optimisation-max-active"`
	Timeout               int      `json:"timeout"`
	EnableCluster         bool     `json:"enable-cluster"`
	UseSSL                bool     `json:"use-ssl"`
	SSLInsecureSkipVerify bool     `json:"ssl-insecure-skip-verify"`
}

const (
	RedisKeyPrefix      = "analytics-"
	defaultRedisAddress = "127.0.0.1:6379"
)

var redisClusterSingleton redis.UniversalClient

type RedisOpts redis.UniversalOptions

func (o *RedisOpts) failover() *redis.FailoverOptions {
	if len(o.Addrs) == 0 {
		o.Addrs = []string{"127.0.0.1:6379"}
	}
	return &redis.FailoverOptions{
		SentinelAddrs:      o.Addrs,
		MasterName:         o.MasterName,
		OnConnect:          o.OnConnect,
		DB:                 o.DB,
		Password:           o.Password,
		MaxRetries:         o.MaxRetries,
		MinRetryBackoff:    o.MinRetryBackoff,
		MaxRetryBackoff:    o.MaxRetryBackoff,
		DialTimeout:        o.DialTimeout,
		ReadTimeout:        o.ReadTimeout,
		WriteTimeout:       o.WriteTimeout,
		PoolSize:           o.PoolSize,
		MinIdleConns:       o.MinIdleConns,
		MaxConnAge:         o.MaxConnAge,
		PoolTimeout:        o.PoolTimeout,
		IdleTimeout:        o.IdleTimeout,
		IdleCheckFrequency: o.IdleCheckFrequency,
		TLSConfig:          o.TLSConfig,
	}
}

func (o *RedisOpts) cluster() *redis.ClusterOptions {
	if len(o.Addrs) == 0 {
		o.Addrs = []string{defaultRedisAddress}
	}

	return &redis.ClusterOptions{
		Addrs:              o.Addrs,
		OnConnect:          o.OnConnect,
		Password:           o.Password,
		MaxRedirects:       o.MaxRedirects,
		ReadOnly:           o.ReadOnly,
		RouteByLatency:     o.RouteByLatency,
		RouteRandomly:      o.RouteRandomly,
		MaxRetries:         o.MaxRetries,
		MinRetryBackoff:    o.MinRetryBackoff,
		MaxRetryBackoff:    o.MaxRetryBackoff,
		DialTimeout:        o.DialTimeout,
		ReadTimeout:        o.ReadTimeout,
		WriteTimeout:       o.WriteTimeout,
		PoolSize:           o.PoolSize,
		MinIdleConns:       o.MinIdleConns,
		MaxConnAge:         o.MaxConnAge,
		PoolTimeout:        o.PoolTimeout,
		IdleTimeout:        o.IdleTimeout,
		IdleCheckFrequency: o.IdleCheckFrequency,
		TLSConfig:          o.TLSConfig,
	}
}

func (o *RedisOpts) simple() *redis.Options {
	addr := defaultRedisAddress
	if len(o.Addrs) > 0 {
		addr = o.Addrs[0]
	}

	return &redis.Options{
		Addr:      addr,
		OnConnect: o.OnConnect,

		DB:       o.DB,
		Password: o.Password,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:  o.DialTimeout,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,

		PoolSize:           o.PoolSize,
		MinIdleConns:       o.MinIdleConns,
		MaxConnAge:         o.MaxConnAge,
		PoolTimeout:        o.PoolTimeout,
		IdleTimeout:        o.IdleTimeout,
		IdleCheckFrequency: o.IdleCheckFrequency,

		TLSConfig: o.TLSConfig,
	}
}

type RedisClusterStorageManager struct {
	db        redis.UniversalClient
	KeyPrefix string
	HashKeys  bool
	Config    RedisOptions
}

// redis 集群
func NewRedisClusterPool(forceReconnect bool, config RedisOptions) redis.UniversalClient {
	if !forceReconnect {
		if redisClusterSingleton != nil {
			return redisClusterSingleton
		}
	} else {
		if redisClusterSingleton != nil {
			redisClusterSingleton.Close()
		}
	}

	maxActive := 500
	if config.MaxActive > 0 {
		maxActive = config.MaxActive
	}
	timeout := 5 * time.Second
	if config.Timeout > 0 {
		timeout = time.Duration(config.Timeout) * time.Second
	}
	var tlsConfig *tls.Config
	if config.UseSSL {
		tlsConfig = &tls.Config{
			// nolint: gosec
			InsecureSkipVerify: config.SSLInsecureSkipVerify,
		}
	}
	var client redis.UniversalClient
	opts := &RedisOpts{
		MasterName:   config.MasterName,
		Addrs:        getRedisAddrs(config),
		DB:           config.Database,
		Password:     config.Password,
		PoolSize:     maxActive,
		IdleTimeout:  240 * time.Second,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		DialTimeout:  timeout,
		TLSConfig:    tlsConfig,
	}

	if opts.MasterName != "" {
		client = redis.NewFailoverClient(opts.failover())
	} else if config.EnableCluster {
		client = redis.NewClusterClient(opts.cluster())
	} else {
		client = redis.NewClient(opts.simple())
	}
	redisClusterSingleton = client

	return client
}

func getRedisAddrs(config RedisOptions) (addrs []string) {
	if len(config.Addrs) != 0 {
		addrs = config.Addrs
	}

	if len(addrs) == 0 && config.Port != 0 {
		addr := config.Host + ":" + strconv.Itoa(config.Port)
		addrs = append(addrs, addr)
	}

	return addrs
}

func (r *RedisClusterStorageManager) GetName() string {
	return "redis"
}

func (r *RedisClusterStorageManager) Init(config interface{}) error {
	r.Config = RedisOptions{}
	r.KeyPrefix = RedisKeyPrefix
	return nil
}

func (r *RedisClusterStorageManager) Connect() bool {
	if r.db == nil {
		r.db = NewRedisClusterPool(false, r.Config)
		return true
	}
	r.db = redisClusterSingleton
	return true
}

func (r *RedisClusterStorageManager) hashKey(in string) string {
	return in
}

func (r *RedisClusterStorageManager) fixKey(keyName string) string {
	setKeyName := r.KeyPrefix + r.hashKey(keyName)
	return setKeyName
}

func (r *RedisClusterStorageManager) GetAndDeleteSet(keyName string) []interface{} {
	if r.db == nil {
		r.Connect()
		return r.GetAndDeleteSet(keyName)
	}
	fixedKey := r.fixKey(keyName)
	var lrange *redis.StringSliceCmd
	r.db.TxPipelined(func(pipe redis.Pipeliner) error {
		lrange = pipe.LRange(fixedKey, 0, -1)
		pipe.Del(fixedKey)
		return nil
	})

	vals := lrange.Val()
	result := make([]interface{}, len(vals))
	for i, v := range vals {
		result[i] = v
	}
	return result

}

func (r *RedisClusterStorageManager) SetKey(keyName, session string, timeout int64) error {

	r.ensureConnection()
	err := r.db.Set(r.fixKey(keyName), session, 0).Err()
	if timeout > 0 {
		if expErr := r.SetExp(keyName, timeout); expErr != nil {
			return expErr
		}
	}
	if err != nil {
		log.Errorf("Error trying to set value: %s", err.Error())

		return errors.Wrap(err, "failed to set key")
	}
	return nil
}

func (r *RedisClusterStorageManager) SetExp(keyName string, timeout int64) error {
	err := r.db.Expire(r.fixKey(keyName), time.Duration(timeout)*time.Second).Err()
	if err != nil {
		log.Errorf("Could not EXPIRE key: %s", err.Error())
	}
	return errors.Wrap(err, "failed to set expire time for key")
}

func (r *RedisClusterStorageManager) ensureConnection() {
	if r.db != nil {
		// already connected
		return
	}
	for {
		r.Connect()
		if r.db != nil {
			// reconnection worked
			return
		}
	}
}
