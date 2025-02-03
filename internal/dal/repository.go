package dal

import (
	"context"
	"errors"

	"github.com/willtowle1/parkn/internal/common/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	msgCreateOneSuccess = "successfully created one"
	msgGetSuccess       = "successfully got from collection"

	errCreateOne = "error while creating one"
	errGet       = "error while getting from collection"
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

	id := res.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

func (r *Dal[D]) Get(ctx context.Context, filter interface{}) ([]D, error) {
	var res []D

	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return res, errors.New(errGet)
	}
	err = cur.All(ctx, &res)
	if err != nil {
		return res, errors.New(errGet)
	}

	return res, nil
}

func (r *Dal[D]) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	res, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}
