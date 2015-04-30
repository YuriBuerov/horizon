package horizon

import (
	"github.com/rcrowley/go-metrics"
	"github.com/rs/cors"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

type Web struct {
	router *web.Mux

	requestTimer metrics.Timer
	failureMeter metrics.Meter
	successMeter metrics.Meter
}

func NewWeb(app *App) (*Web, error) {

	api := web.New()
	result := Web{
		router:       api,
		requestTimer: metrics.NewTimer(),
		failureMeter: metrics.NewMeter(),
		successMeter: metrics.NewMeter(),
	}

	app.metrics.Register("requests.total", result.requestTimer)
	app.metrics.Register("requests.succeeded", result.successMeter)
	app.metrics.Register("requests.failed", result.failureMeter)

	installMiddleware(api, app)
	installActions(api)

	return &result, nil
}

func installMiddleware(api *web.Mux, app *App) {
	api.Use(middleware.RequestID)
	api.Use(middleware.Logger)
	api.Use(middleware.Recoverer)
	api.Use(middleware.AutomaticOptions)
	api.Use(appMiddleware(app))
	api.Use(requestMetricsMiddleware)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	api.Use(c.Handler)
}

func installActions(api *web.Mux) {
	api.Get("/", rootAction)
	api.Get("/metrics", metricsAction)

	// ledger actions
	api.Get("/ledgers", ledgerIndexAction)
	api.Get("/ledgers/:id", ledgerShowAction)
	api.Get("/ledgers/:ledger_id/transactions", notFoundAction)
	api.Get("/ledgers/:ledger_id/operations", notFoundAction)
	api.Get("/ledgers/:ledger_id/effects", notFoundAction)

	// account actions
	api.Get("/accounts", notFoundAction)
	api.Get("/accounts/:id", notFoundAction)
	api.Get("/accounts/:account_id/transactions", notFoundAction)
	api.Get("/accounts/:account_id/operations", notFoundAction)
	api.Get("/accounts/:account_id/effects", notFoundAction)

	// transaction actions
	api.Get("/transactions", notFoundAction)
	api.Get("/transactions/:id", notFoundAction)
	api.Get("/transactions/:tx_id/operations", notFoundAction)
	api.Get("/transactions/:tx_id/effects", notFoundAction)

	// transaction actions
	api.Get("/operations", notFoundAction)
	api.Get("/operations/:id", notFoundAction)
	api.Get("/operations/:tx_id/effects", notFoundAction)

	api.Get("/stream", streamAction)

	api.NotFound(notFoundAction)
}
