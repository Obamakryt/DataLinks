package rest

import (
	"DataLinks/internal/api/rest/handler"
	"DataLinks/internal/service/jwt_hash"
	"DataLinks/pkg/config/parser"
	contextmiddleware "DataLinks/pkg/context_helper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"
)

func Router(dataServer parser.Server, sigh jwt_hash.JWTSigh, h1 handler.AuthRegHandler, h2 handler.LinkHandler) error {
	e := echo.New()

	e.Use(middleware.Logger(),
		middleware.Recover(),
		middleware.RequestID(),
		contextmiddleware.MWUserId(sigh),
		contextmiddleware.MWRequestId,
	)

	e.POST("/reg", h1.RegHandler, contextmiddleware.WithTimeout(3))
	e.POST("/login", h1.LoginHandler, contextmiddleware.WithTimeout(4))

	e.GET("/mylinks", h2.GetUserLinksHandler, contextmiddleware.WithTimeout(3))
	e.POST("/add", h2.AddNewLinkHandler, contextmiddleware.WithTimeout(3))
	e.PATCH("/change", h2.SwapLinks, contextmiddleware.WithTimeout(3))
	e.DELETE("/delete", h2.DeleteLinkHandler, contextmiddleware.WithTimeout(3))

	server := &http.Server{
		Addr:         dataServer.Settings.Addr,
		Handler:      e,
		ReadTimeout:  dataServer.Settings.ReadTimeout * time.Second,
		WriteTimeout: dataServer.Settings.ReadTimeout * time.Second,
		IdleTimeout:  dataServer.Settings.ReadTimeout * time.Second,
	}
	log.Info("starting server")
	err := server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
