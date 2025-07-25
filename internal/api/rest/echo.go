package rest

import (
	"DataLinks/internal/api/rest/handler"
	contexth "DataLinks/pkg/context_helper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func Router() {
	e := echo.New()

	e.Use(middleware.Logger(),
		middleware.Recover(),
	)
	h := handler.Handler{}
	e.POST("/reg", h.RegHandler)
	e.POST("/login", h.LogHandler, contexth.WithTimeout(4))

	server := &http.Server{
		Addr: ,
		Handler: e,
		ReadTimeout: ,
		WriteTimeout: ,
		IdleTimeout: ,
	}

	server.ListenAndServe()
}
