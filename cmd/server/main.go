package main

import (
	"context"
	"fmt"
	"github.com/TheTeaParty/monitor/internal"
	"github.com/TheTeaParty/monitor/internal/domain/report"
	"github.com/TheTeaParty/monitor/internal/domain/service"
	"github.com/TheTeaParty/monitor/internal/handler"
	"github.com/TheTeaParty/monitor/internal/util"
	monitorAPI "github.com/TheTeaParty/monitor/pkg/api/openapi"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"time"
)

func main() {

	app := internal.New("mus", "monitor")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	connection, err := util.NewMongoDBSessionWithEnv(ctx)
	if err != nil {
		app.Logger.Fatal("Error connecting mongodb",
			zap.Error(err))
	}

	db := connection.Database(
		util.GetEnv("MONGODB_DATABASE", fmt.Sprintf("%v-%v", app.AppName, app.ServiceName)))

	dir, _ := os.Getwd()
	app.ReportRepository = report.NewMongoDB(db)
	app.ServiceRepository = service.NewFile(filepath.Join(dir, "data/services"))

	monitorCtx, _ := context.WithCancel(context.Background())
	defer monitorCtx.Done()

	go app.RunServiceWatcher(monitorCtx)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	h := handler.NewHandler(app)
	r.Mount("/", monitorAPI.HandlerFromMux(h, r))

	if err := app.RunHTTP(r); err != nil {
		app.Logger.Fatal("Can't start http handler",
			zap.Error(err))
	}
}
