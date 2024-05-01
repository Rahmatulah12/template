package redis

import (
	"crypto/tls"
	"time"

	nrRedis "github.com/newrelic/go-agent/v3/integrations/nrredis-v9"
	"github.com/redis/go-redis/v9"
)

var (
	redisOpts    *redis.Options
	redisClusterOpts *redis.ClusterOptions
)

func(r *Conf) InitClient() *redis.Client {
	redisOpts = &redis.Options{
		Addr:       r.Addr,
		ClientName: r.ClientName,
		Username: r.Username,
		Password: r.Pass,
		DB:                    r.DB,
		MaxRetries:            r.MaxRetries,
		DialTimeout:           time.Duration(r.DialTimeout) * time.Second,
		ReadTimeout:           time.Duration(r.ReadTimeout) * time.Second,
		WriteTimeout:          time.Duration(r.WriteTimeout) * time.Second,
		PoolSize:              r.PoolSize,
		PoolTimeout:           time.Duration(r.PoolTimeout) * time.Second,
		MinIdleConns:          r.MinIdleConns,
		MaxIdleConns:          r.MaxIdleConns,
		MaxActiveConns:        r.MaxActiveConns,
		ConnMaxIdleTime:       time.Duration(r.ConnMaxIdleTime) * time.Second,
		ConnMaxLifetime:       time.Duration(r.ConnMaxLifetime) * time.Minute,
		TLSConfig: &tls.Config{InsecureSkipVerify: false},
		ContextTimeoutEnabled: true,
	}

	if r.IsUseTls {
		redisOpts.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	
	redisClient := redis.NewClient(redisOpts)

	if r.IsUseHooks {
		redisClient.AddHook(nrRedis.NewHook(redisOpts))
	}

	return redisClient
}

func(r *Conf) InitClusterClient() *redis.ClusterClient {
	redisClusterOpts = &redis.ClusterOptions{
		Addrs:      []string{r.Addr},
		ClientName: r.ClientName,
		Username: r.Username,
		Password: r.Pass,
		MaxRetries:            r.MaxRetries,
		DialTimeout:           time.Duration(r.DialTimeout) * time.Second,
		ReadTimeout:           time.Duration(r.ReadTimeout) * time.Second,
		WriteTimeout:          time.Duration(r.WriteTimeout) * time.Second,
		PoolSize:              r.PoolSize,
		PoolTimeout:           time.Duration(r.PoolTimeout) * time.Second,
		MinIdleConns:          r.MinIdleConns,
		MaxIdleConns:          r.MaxIdleConns,
		MaxActiveConns:        r.MaxActiveConns,
		ConnMaxIdleTime:       time.Duration(r.ConnMaxIdleTime) * time.Second,
		ConnMaxLifetime:       time.Duration(r.ConnMaxLifetime) * time.Minute,
		TLSConfig: &tls.Config{InsecureSkipVerify: false},
		ContextTimeoutEnabled: true,
	}

	redisClusterClient := redis.NewClusterClient(redisClusterOpts)

	if r.IsUseHooks {
		redisClusterClient.AddHook(nrRedis.NewHook(nil))
	}

	return redisClusterClient
}