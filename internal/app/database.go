package app

import (
	"context"

	"github.com/willtowle1/parkn/internal/common/logger"
	"github.com/willtowle1/parkn/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDatabase(ctx context.Context, logger logger.Logger, errs chan error, config config.Config) (*mongo.Client, error) {

	connectOptions := options.Client().
		ApplyURI(config.MongoConnectionString).
		SetAppName(config.MongoAppName)

	client, err := mongo.Connect(ctx, connectOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}
