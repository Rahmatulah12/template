package mysql

import (
	"context"
	"time"

	_ "github.com/newrelic/go-agent/v3/integrations/nrmysql"
	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func(m *MysqlDbConf) InitGorm() *gorm.DB {
	if m.MaxIdleConn == 0 { m.MaxIdleConn = 25 }

	if m.MaxOpenConn == 0 { m.MaxOpenConn = 250 }

	if m.MaxConnLifeTime == 0 { m.MaxConnLifeTime = 24 }

	if m.MaxIdleConnLifeTime == 0 { m.MaxIdleConnLifeTime = 1 }

	if m.Timeout == 0 { m.Timeout = 60 }

	dsn, err := m.FormatDSN()

	if err != nil {
		panic("Failed connect to database. Error :" + err.Error())
	}

	dialector := mysql.Dialector{
		Config: &mysql.Config{
			DriverName:                    "nrmysql",
			DSN:                           dsn,
		},
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: NewLogger,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		QueryFields: true,
	})

	if err != nil {
		panic("Failed connect to database. Error :" + err.Error())
	}

	// Enable Database Connection Pool, for reusable connection. if connection is available. use existing connection
	sqlDB, err := db.DB()

	if err != nil {
		panic("Failed connect to database. Error :" + err.Error())
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(m.MaxIdleConn)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(m.MaxOpenConn)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Duration(m.MaxConnLifeTime) * time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Duration(m.MaxIdleConnLifeTime) * time.Hour)

	if m.AutoMigrate {
		/*
			* run auto migrate gorm, only for development. initiate entity/model here
			* call model here
			* example
			* db.AutoMigrate( 
				&model.User{},
			  )
		*/
		err = db.AutoMigrate()
		if err != nil {
			panic("Failed to run auto migration table. Error :" + err.Error())
		}
	}
	gormTransactionTrace := m.newrelicTransaction.StartTransaction("GORM Operation")
	gormTransactionContext := newrelic.NewContext(context.Background(), gormTransactionTrace)
	tracedDB := db.WithContext(gormTransactionContext)
	return tracedDB
}