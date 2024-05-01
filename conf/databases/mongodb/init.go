package mongodb

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func(c *Conf) InitMongoDB() *mongo.Database {
	defer c.ctxCancel()
	if client != nil {
		return client.Database(c.DBName)
	}

	opt := options.Client().ApplyURI(c.Dsn)
	client, err := mongo.NewClient(opt)
	if err != nil {
		return nil
	}

	err = client.Connect(c.ctx)
	if err != nil {
		return nil
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		panic(err)
	} else {
		logrus.Println("Connected!")
	}

	return client.Database(c.DBName)
}

func(c *Conf) CloseDB() {
	if c.client == nil {
		return
	}

	if err := client.Disconnect(c.ctx); err != nil {
		panic(err)
	}

	log.Println("Disconnected!")
}