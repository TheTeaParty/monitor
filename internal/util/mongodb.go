package util

import (
	"context"
	"crypto/tls"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"strings"
	"time"
)

// NewSessionWithEnv new session using environment variables
func NewMongoDBSessionWithEnv(ctx context.Context) (*mongo.Client, error) {
	opts := options.Client()
	opts.SetHosts(strings.Split(os.Getenv("MONGODB_HOSTS"), ","))

	if os.Getenv("MONGODB_USERNAME") != "" && os.Getenv("MONGODB_PASSWORD") != "" {
		opts.SetAuth(options.Credential{
			Username: os.Getenv("MONGODB_USERNAME"),
			Password: os.Getenv("MONGODB_PASSWORD"),
		})
	}

	opts.SetConnectTimeout(10 * time.Second)

	if os.Getenv("MONGODB_TLS") == "yes" {
		tlsConfig := &tls.Config{}
		opts.SetTLSConfig(tlsConfig)
	}

	db, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return db, nil
}
