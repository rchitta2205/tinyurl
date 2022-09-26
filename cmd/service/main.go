package main

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"sync"
	"tinyurl/pkg"
	"tinyurl/pkg/config"
)

func main() {
	var ctx context.Context

	ctx = context.Background()
	logEntry := logrus.NewEntry(logrus.New())

	// Initialize the configuration
	cfg := config.NewConfig()

	dbOptions := options.Client().ApplyURI(cfg.DbAddr)
	db, err := mongo.Connect(ctx, dbOptions)
	if err != nil {
		logEntry.Warnf("db not initialized, using mock db")
	} else {
		err = db.Ping(ctx, readpref.Primary())
		if err != nil {
			logEntry.Warnf("db not connected, using mock db")
			db = nil
		}
	}

	cacheOptions := &redis.Options{Addr: cfg.CacheAddr}
	cache := redis.NewClient(cacheOptions)
	err = cache.Ping().Err()
	if err != nil {
		logEntry.Warnf("cache not connected...")
		cache = nil
	}

	grpcService := pkg.NewGrpcService(ctx, cfg, logEntry, db, cache)
	restService := pkg.NewRestService(ctx, cfg, logEntry)

	wg := &sync.WaitGroup{}
	err = grpcService.Register()
	if err != nil {
		logEntry.Fatalf("Failed to create grpc server object and register apps: %+v", err.Error())
	}

	err = grpcService.Serve(wg)
	if err != nil {
		logEntry.Fatalf("Failed to start the gRPC server: %+v", err.Error())
	}

	err = restService.Register()
	if err != nil {
		logEntry.Fatalf("Failed to create rest server and register apps: %+v", err.Error())
	}

	err = restService.Serve(wg)
	if err != nil {
		logEntry.Fatalf("Failed to start the rest server: %+v", err.Error())
	}

	// Wait for all goroutines to be completed
	wg.Wait()

	// Close the database
	err = db.Disconnect(ctx)
	if err != nil {
		logEntry.Fatalf("Failed to disconnect from the database: %+v", err.Error())
	}

	// Close the cache
	err = cache.Close()
	if err != nil {
		logEntry.Fatalf("Failed to disconnect from the cache: %+v", err.Error())
	}
}
