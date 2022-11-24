package repository

import (
	"context"
	"fmt"
	"moonspace/model"
	"moonspace/repository/types"
	"time"

	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MakeRepositoryClient(cfg types.Config) any {
	switch cfg.Type {
	case types.Mongo:
		return createMongoCLI(cfg)
	case types.Postgres:
		return createPostgresCLI(cfg)
	default:
		return nil
	}
}

func createPostgresCLI(cfg types.Config, config ...any) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.Url))
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&model.Category{}, &model.Order{}, &model.Cart{}, &model.Product{})
	if err != nil {
		panic(fmt.Errorf("Error creating database: %w", err))
	}

	return db
}

func createUniqueIndexMongo(cli *qmgo.Client, dbName string) {
	tables := []string{"cart", "order"}

	for _, t := range tables {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		col := cli.Database(dbName).Collection(t)
		key := t

		if key == "cart" {
			key = "user_id"
		} else if key == "order" {
			key = key + "_id"
		}

		unique := true
		indexModel := options.IndexModel{
			Key: []string{key},
			IndexOptions: &mongo_options.IndexOptions{
				Unique: &unique,
			},
		}

		if err := col.CreateOneIndex(ctx, indexModel); err != nil {
			panic(err)
		}
	}
}

func createMongoCLI(cfg types.Config) *qmgo.QmgoClient {
	connectionTimeout := types.Timeout.Milliseconds()
	ctx, cancel := context.WithTimeout(context.Background(), types.Timeout)
	defer cancel()
	qmgoCfg := &qmgo.Config{
		Uri:              cfg.Url,
		Database:         cfg.Database,
		ConnectTimeoutMS: &connectionTimeout,
	}

	cli, err := qmgo.Open(ctx, qmgoCfg)
	if err != nil {
		panic("error creating a mongo repository: " + err.Error())
	}

	// createUniqueIndexMongo(cli.Client, cfg.Database)

	return cli
}
