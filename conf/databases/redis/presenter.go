package redis

import "github.com/sirupsen/logrus"

type Conf struct {
	ClientName string
	Addr string
	Username string
	Pass string
	DB	int
	MaxRetries int
	MaxActiveConns int
	ConnMaxIdleTime int
	ConnMaxLifetime int
	MinIdleConns int
	MaxIdleConns int
	DialTimeout int
	ReadTimeout int
	WriteTimeout int
	PoolTimeout int
	PoolSize int
	IsUseTls bool
	IsUseHooks bool
}

func NewConn(conf *Conf) *Conf {
	if conf == nil { logrus.Panicln("Failed to connect redis") }

	if conf.MaxRetries == 0 { conf.MaxRetries = 3 }

	if conf.DialTimeout == 0 { conf.DialTimeout = 10 } // seconds

	if conf.ReadTimeout == 0 { conf.ReadTimeout = 10 } // seconds

	if conf.WriteTimeout == 0 { conf.WriteTimeout = 10 } // seconds

	if conf.MaxActiveConns == 0 { conf.MaxActiveConns = 250 }

	if conf.PoolSize == 0 { conf.PoolSize = conf.MaxActiveConns }

	if conf.PoolTimeout == 0 { conf.PoolTimeout = 10 }

	if conf.MinIdleConns == 0 { conf.MinIdleConns = 25 }

	if conf.MaxIdleConns == 0 { conf.MaxIdleConns = 50 }

	if conf.ConnMaxIdleTime == 0 { conf.ConnMaxIdleTime = 60 } // seconds

	if conf.ConnMaxLifetime == 0 { conf.ConnMaxLifetime = 60 } // minutes

	return conf
}