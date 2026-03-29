package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "backend/docs"
)

// @title			Backend API
// @version		1.0
// @description	Backend Service
// @BasePath		/api
func main() {
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"X-Requested-With", "Content-Type", "Authorization", "Sync-Session-Id"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowCredentials: true,
	}))

	router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		LogLatency:  true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				log.Printf("method=%s uri=%s status=%d latency=%s",
					v.Method, v.URI, v.Status, v.Latency)
			} else {
				log.Printf("method=%s uri=%s status=%d latency=%s err=%s",
					v.Method, v.URI, v.Status, v.Latency, v.Error)
			}
			return nil
		},
	}))

	router.GET("/swagger/*", echoSwagger.WrapHandler)

	api := router.Group("/api")
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	authRouteController.RegisterRoute(api)
	userRouteController.RegisterRoute(api)

	router.HTTPErrorHandler = func(err error, c echo.Context) {
		if he, ok := err.(*echo.HTTPError); ok {
			_ = c.JSON(he.Code, map[string]interface{}{
				"status":  "fail",
				"message": he.Message,
			})
			return
		}
		_ = c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "error",
			"message": "Internal server error.",
		})
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.ServerPort),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server listen on port %s: %v", config.ServerPort, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}

	log.Println("shutdown gracefully")
}
