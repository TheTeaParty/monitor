package handler

import (
	"github.com/TheTeaParty/monitor/internal"
	"github.com/TheTeaParty/monitor/internal/domain"
	monitorAPI "github.com/TheTeaParty/monitor/pkg/api/openapi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type handler struct {
	app *internal.Application
}

func (h *handler) GetReports(w http.ResponseWriter, r *http.Request) {

	serviceURLs := strings.Split(
		strings.ReplaceAll(r.URL.Query().Get("serviceURLs"), " ", ""), ",")

	criteria := domain.ReportCriteria{
		ServiceURLs: serviceURLs,
		Status:      -1,
	}

	criteria.ReportedAtFrom = time.Now().Add(-1 * time.Minute).Unix()

	if r.URL.Query().Get("status") != "" {
		criteria.Status = domain.ServiceStatusUnavailable
		if r.URL.Query().Get("status") == "available" {
			criteria.Status = domain.ServiceStatusAvailable
		}
	}

	if r.URL.Query().Get("reportedAtFrom") != "" {
		criteria.ReportedAtFrom, _ = strconv.ParseInt(r.URL.Query().Get("reportedAtFrom"), 10, 64)
	}
	if r.URL.Query().Get("reportedAtTo") != "" {
		criteria.ReportedAtTo, _ = strconv.ParseInt(r.URL.Query().Get("reportedAtTo"), 10, 64)
	}
	if r.URL.Query().Get("responseTimeMoreThen") != "" {
		criteria.ResponseTimeMoreThen, _ = strconv.ParseInt(r.URL.Query().Get("responseTimeMoreThen"), 10, 64)
	}
	if r.URL.Query().Get("responseTimeLessThen") != "" {
		criteria.ResponseTimeLessThen, _ = strconv.ParseInt(r.URL.Query().Get("responseTimeLessThen"), 10, 64)
	}

	reports, err := h.app.ReportRepository.GetMatching(r.Context(), criteria)
	if err != nil {
		_ = render.Render(w, r, h.errRender(err))
	}

	reportsResponse := make([]*monitorAPI.Report, len(reports))
	for i, r := range reports {

		status := "available"
		if r.Status == domain.ServiceStatusUnavailable {
			status = "unavailable"
		}

		reportsResponse[i] = &monitorAPI.Report{
			CreatedAt:    int(r.CreatedAt),
			Id:           r.ID,
			ResponseTime: int(r.ResponseTime),
			ServiceURL:   r.ServiceURL,
			Status:       status,
			Details:      r.Details,
		}
	}

	render.JSON(w, r, reportsResponse)
}

func NewHandler(app *internal.Application) monitorAPI.ServerInterface {
	return &handler{app: app}
}

func (h *handler) errRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}
