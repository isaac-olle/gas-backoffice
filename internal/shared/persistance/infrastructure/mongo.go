package infrastructure

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrInsertingMongoObject app_errors.ErrorCode = "INSERTING_MONGO_OBJECT"

type MongoRepository[T any] struct {
	client   *mongo.Client
	database string
}

func NewMongoRepository[T any](client *mongo.Client, database string) *MongoRepository[T] {
	return &MongoRepository[T]{client: client, database: database}
}

func (m *MongoRepository[T]) Insert(ctx context.Context, collection string, object any) error {
	_, err := m.client.Database(m.database).Collection(collection).InsertOne(ctx, object)
	if err != nil {
		return app_errors.InternalError(ErrInsertingMongoObject).
			WithDescription(fmt.Sprintf("cannot insert object into collection: %s", err.Error()))
	}
	return nil
}

func (m *MongoRepository[T]) SearchById(ctx context.Context, collection string, key, id string, toDomain func(result *mongo.SingleResult) (T, error)) (T, error) {
	filter := bson.D{{Key: key, Value: id}}
	result := m.client.Database(m.database).Collection(collection).FindOne(ctx, filter)
	return toDomain(result)
}
