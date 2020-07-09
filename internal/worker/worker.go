package worker

import (
	"context"
	"fmt"
	"github.com/TheTeaParty/monitor/internal/domain"
	"net/http"
	"runtime"
	"time"
)

func doRequests(urls []string, concurrencyLimit int) []*domain.Report {

	semaphoreChan := make(chan struct{}, concurrencyLimit)

	resultsChan := make(chan *domain.Report)

	defer func() {
		close(semaphoreChan)
		close(resultsChan)
	}()

	reportedAt := time.Now().Unix()

	for i, u := range urls {
		go func(i int, u string) {

			semaphoreChan <- struct{}{}

			start := time.Now()

			r := &domain.Report{
				ReportedAt:   reportedAt,
				ResponseTime: 0,
				ServiceURL:   u,
				Status:       domain.ServiceStatusUnavailable,
			}

			res, err := http.Get(fmt.Sprintf("https://%v", u))
			if err == nil && res.StatusCode != 503 {
				r.Details = "OK"
				r.Status = domain.ServiceStatusAvailable
			}

			if err != nil {
				r.Details = err.Error()
			}

			r.ResponseTime = time.Since(start).Nanoseconds()

			resultsChan <- r

			<-semaphoreChan

		}(i, u)
	}

	var results []*domain.Report

	for {
		result := <-resultsChan
		results = append(results, result)

		if len(results) == len(urls) {
			break
		}
	}

	return results
}

func RunServiceMonitor(ctx context.Context, serviceURLs []*domain.Service) (chan []*domain.Report, <-chan error, error) {
	out := make(chan []*domain.Report)
	errC := make(chan error, 1)

	ticker := time.NewTicker(5 * time.Second)

	go func() {
		defer close(out)
		defer close(errC)

		for range ticker.C {

			urls := make([]string, len(serviceURLs))
			for i, s := range serviceURLs {
				urls[i] = s.URL
			}

			results := doRequests(urls, runtime.GOMAXPROCS(runtime.NumCPU()))

			reports := make([]*domain.Report, len(serviceURLs))
			for i, result := range results {
				reports[i] = result
			}

			select {
			case out <- reports:
			case <-ctx.Done():
				return
			}
		}

	}()

	return out, errC, nil
}
