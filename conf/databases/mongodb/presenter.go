package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	client      *mongo.Client
)

type Conf struct {
	Dsn		string
	DBName	string
	Timeout	int
	client *mongo.Client
	ctx context.Context
	ctxCancel context.CancelFunc
}

func(c *Conf) NewConn() *Conf {
	if c.Dsn == "" { panic("DSN could not be empty.") }

	if c.DBName == "" { panic("DB Name could not be empty.") }

	if c.Timeout == 0 { c.Timeout = 10 }

	var ctx, cancel = context.WithTimeout(context.Background(), time.Duration(c.Timeout) * time.Second)

	c.client = client
	c.ctx = ctx
	c.ctxCancel = cancel

	return c
}