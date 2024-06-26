package dal

import (
	"context"
	"errors"

	"github.com/willtowle1/parkn/internal/common/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	msgCreateOneSuccess = "successfully created one"
	msgGetOneSuccess    = "successfully got one"

	errCreateOne = "error while creating one"
	errGetOne    = "error while getting one"
)

type Dal[D any] struct {
	logger     logger.Logger
	collection mongo.Collection
}

func NewRepository[D any](logger logger.Logger, collection mongo.Collection) *Dal[D] {
	return &Dal[D]{
		logger:     logger,
		collection: collection,
	}
}

func (r *Dal[D]) CreateOne(ctx context.Context, input D) (string, error) {
	res, err := r.collection.InsertOne(ctx, input)
	if err != nil {
		return "", errors.New(errCreateOne)
	}

	id := res.InsertedID.(primitive.ObjectID)
	r.logger.Info(ctx, msgCreateOneSuccess, "id", id.Hex())
	return id.Hex(), nil
}

func (r *Dal[D]) GetOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (D, error) {
	var res D
	err := r.collection.FindOne(ctx, filter, opts...).Decode(&res)
	if err != nil {
		return res, errors.New(errGetOne)
	}

	r.logger.Info(ctx, msgGetOneSuccess, res)
	return res, nil
}

func (r *Dal[D]) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	res, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}
