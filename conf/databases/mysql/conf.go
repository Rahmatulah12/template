package mysql

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/gorm/logger"
)

var (
	NewLogger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)
)

type MysqlDbConf struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Timeout int
	MaxIdleConn int
	MaxOpenConn int
	MaxConnLifeTime int
	MaxIdleConnLifeTime int
	AutoMigrate bool
	newrelicTransaction *newrelic.Application
}

func (m *MysqlDbConf) FormatDSN() (string, error) {
	location, err := time.LoadLocation("Asia/Jakarta")

	if err != nil {
		return "", err
	}

	if m.Timeout == 0 { m.Timeout = 30 }

	conn := mysql.Config{
		User:                 m.User,
		Passwd:               m.Password,
		DBName:               m.Database,
		Addr:                 fmt.Sprintf("%s:%s", m.Host, m.Port),
		Net:                  "tcp",
		ParseTime:            true,
		Loc:                  location,
		AllowNativePasswords: true,
		Timeout:              time.Duration(m.Timeout) * time.Second,
	}

	return conn.FormatDSN(), nil
}