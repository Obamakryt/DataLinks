package rest

import (
	"DataLinks/internal/api/rest/handler"
	"DataLinks/pkg/config/parser"
	contextmiddleware "DataLinks/pkg/context_helper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
)

func Router(dataServer parser.ServerSettings) {
	e := echo.New()

	e.Use(middleware.Logger(),
		middleware.Recover(),
		middleware.RequestID(),
	)
	h := handler.Handler{}
	e.POST("/reg", h.RegHandler)
	e.POST("/login", h.LogHandler, contextmiddleware.WithTimeout(4))

	server := &http.Server{
		Addr:         dataServer.Addr,
		Handler:      e,
		ReadTimeout:  dataServer.ReadTimeout * time.Second,
		WriteTimeout: dataServer.ReadTimeout * time.Second,
		IdleTimeout:  dataServer.ReadTimeout * time.Second,
	}

	server.ListenAndServe()
}
