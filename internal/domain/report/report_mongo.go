package report

import (
	"context"
	"fmt"
	"github.com/TheTeaParty/monitor/internal/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const ReportsCollection = "reports"

type reportMongodb struct {
	db *mongo.Database
}

func (r *reportMongodb) Create(ctx context.Context, report *domain.Report) error {

	report.ID = uuid.New().String()
	report.CreatedAt = time.Now().Unix()

	c := r.db.Collection(ReportsCollection)

	if _, err := c.InsertOne(ctx, report); err != nil {
		return err
	}

	return nil
}

func (r *reportMongodb) GetMatching(ctx context.Context, criteria domain.ReportCriteria) ([]*domain.Report, error) {

	c := r.db.Collection(ReportsCollection)

	filter := bson.D{}

	if len(criteria.ServiceURLs) > 0 && criteria.ServiceURLs[0] != "" {
		filter = append(filter, bson.E{
			Key: "serviceUrl",
			Value: bson.D{{
				Key: "$in", Value: criteria.ServiceURLs,
			}},
		})
	}

	filter = append(filter, bson.E{Key: "responseTime",
		Value: bson.D{{"$gt", criteria.ResponseTimeMoreThen}}})
	filter = append(filter, bson.E{Key: "reportedAt",
		Value: bson.D{{"$gt", criteria.ReportedAtFrom}}})

	if criteria.ResponseTimeLessThen != 0 {
		filter = append(filter, bson.E{Key: "responseTime",
			Value: bson.D{{"$lt", criteria.ResponseTimeLessThen}}})
	}

	if criteria.ReportedAtTo != 0 {
		filter = append(filter, bson.E{Key: "reportedAt",
			Value: bson.D{{"$lt", criteria.ReportedAtTo}}})
	}

	if criteria.Status != -1 {
		filter = append(filter, bson.E{Key: "status", Value: criteria.Status})
	}

	fmt.Println(filter)

	cur, err := c.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var reports []*domain.Report
	if err := cur.All(ctx, &reports); err != nil {
		return nil, err
	}

	return reports, nil
}

func NewMongoDB(db *mongo.Database) domain.ReportRepository {
	return &reportMongodb{db: db}
}
