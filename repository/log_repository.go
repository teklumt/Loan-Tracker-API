package repository

import (
	"context"
	"loan-tracker-api/domain"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogRepositoryImpl struct {
    LogCollection *mongo.Collection
}

func NewLogRepositoryImpl(logCollection *mongo.Collection) domain.LogRepository {
    return &LogRepositoryImpl{LogCollection: logCollection}
}

func (lr *LogRepositoryImpl) GetSystemLogs(logType, limit, page string) ([]domain.SystemLog, error) {
	var logs []domain.SystemLog

	// Build the query filter for log type or details search
	filter := bson.M{}
	if logType != "all" {
		filter = bson.M{"details": bson.M{"$regex": logType, "$options": "i"}} // Case-insensitive search in details field
	}

	// Convert limit to int64
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}

	// Convert page to int
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}

	// Set pagination options
	findOptions := options.Find()
	findOptions.SetLimit(int64(limitInt))                  // Set the limit
	findOptions.SetSkip(int64((pageInt - 1) * limitInt))    // Set the skip based on page and limit

	// Query the database with the filter and pagination options
	cursor, err := lr.LogCollection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Decode the results into the logs slice
	if err = cursor.All(context.Background(), &logs); err != nil {
		return nil, err
	}

	return logs, nil
}



func (lr *LogRepositoryImpl) CreateLog(log domain.SystemLog) error {
    log.ID = primitive.NewObjectID().Hex()
    _, err := lr.LogCollection.InsertOne(context.Background(), log)
    return err
}