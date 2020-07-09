package internal

import (
	"context"
	"github.com/TheTeaParty/monitor/internal/domain"
	"github.com/TheTeaParty/monitor/internal/util"
	"github.com/TheTeaParty/monitor/internal/worker"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"net"
	"net/http"
)

type Application struct {
	AppName     string
	ServiceName string

	ServiceRepository domain.ServiceRepository
	ReportRepository  domain.ReportRepository

	Logger *zap.Logger

	listener net.Listener
}

func (a *Application) RunServiceWatcher(ctx context.Context) error {

	services, err := a.ServiceRepository.GetAll()
	if err != nil {
		a.Logger.Error("Error getting services",
			zap.Error(err))
		return err
	}

	out, _, err := worker.RunServiceMonitor(ctx, services)
	if err != nil {
		a.Logger.Error("Error running monitor",
			zap.Error(err))
		return err
	}

	defer close(out)

	go func() {
		for reports := range out {
			for _, r := range reports {
				if err := a.ReportRepository.Create(ctx, r); err != nil {
					a.Logger.Error("Error saving report",
						zap.Error(err))
				}
			}
		}
	}()

	<-ctx.Done()
	return nil
}

// RunHTTP run service HTTP handler
func (a *Application) RunHTTP(r http.Handler) error {
	listener, err := net.Listen("tcp", util.GetEnv("PORT", ":0"))
	if err != nil {
		return err
	}

	a.listener = listener

	a.Logger.Info("Starting http handler",
		zap.String("port", a.listener.Addr().String()))
	if err := http.Serve(a.listener, r); err != nil {
		return err
	}

	return nil
}

func (a *Application) Stop() error {

	return a.listener.Close()
}

func New(appName, serviceName string) *Application {

	_ = godotenv.Load(".env")

	environment := util.GetEnv("ENVIRONMENT", "development")
	logger, _ := zap.NewProduction()

	if environment == "development" {
		logger, _ = zap.NewDevelopment()
	}

	return &Application{
		AppName:     appName,
		ServiceName: serviceName,
		Logger:      logger,
	}
}
