package service

import (
	"bufio"
	"fmt"
	"github.com/TheTeaParty/monitor/internal/domain"
	"os"
)

type serviceFile struct {
	fileName string
}

func (r *serviceFile) GetAll() ([]*domain.Service, error) {
	file, err := os.OpenFile(r.fileName, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var services []*domain.Service
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		services = append(services, &domain.Service{
			ID:  fmt.Sprintf("%v", lineNum),
			URL: scanner.Text(),
		})
	}

	return services, nil
}

func NewFile(fileName string) domain.ServiceRepository {
	return &serviceFile{fileName: fileName}
}
