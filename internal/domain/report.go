package domain

import "context"

type ServiceStatus int

const (
	ServiceStatusAvailable ServiceStatus = iota
	ServiceStatusUnavailable
)

type Report struct {
	ID           string        `json:"id" bson:"_id"`
	CreatedAt    int64         `json:"createdAt" bson:"createdAt"`
	ReportedAt   int64         `json:"reportedAt" bson:"reportedAt"`
	ServiceURL   string        `json:"serviceUrl" bson:"serviceUrl"`
	ResponseTime int64         `json:"responseTime" bson:"responseTime"`
	Status       ServiceStatus `json:"status" bson:"status"`
	Details      string        `json:"details" bson:"details"`
}

type ReportCriteria struct {
	ServiceURLs          []string
	ResponseTimeMoreThen int64
	ResponseTimeLessThen int64
	ReportedAtFrom       int64
	ReportedAtTo         int64
	Status               ServiceStatus
}

type ReportRepository interface {
	Create(ctx context.Context, report *Report) error
	GetMatching(ctx context.Context, criteria ReportCriteria) ([]*Report, error)
}
