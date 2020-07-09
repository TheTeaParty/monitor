package domain

type Service struct {
	ID  string `json:"id" bson:"_id"`
	URL string `json:"url" bson:"url"`
}

type ServiceRepository interface {
	GetAll() ([]*Service, error)
}
