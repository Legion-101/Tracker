package main

import (
	"net/http"
	"os"

	op "tracker/internal"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
)


func main() {
	file, err := os.OpenFile(
		"/var/log/tracker.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	e := echo.New()

	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method"},
	)
	requestDuration := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests.",
		},
		[]string{"method"},
	)

	server := op.NewServer(file, requestsTotal, requestDuration)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	})) // CORS (Cross-Origin Resource Sharing())

	e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{}))
	e.GET("/metrics", echoprometheus.NewHandler())

	op.RegisterHandlers(e, server)
	addres := ":8080"

	e.Logger.Fatal(e.Start(addres))
}
