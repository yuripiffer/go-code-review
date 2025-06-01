package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"coupon_service/internal/api/auth"
	"coupon_service/internal/config"
	"coupon_service/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type API struct {
	srv *http.Server
	mux *gin.Engine
	svc service.CouponService
	cfg config.Config
}

// New creates a new API instance with the provided configuration and service.
//
// @title Coupon Service API
// @version 1.0
// @description This is an API for managing coupons.
// @host localhost:8080
// @BasePath /api/
// @schemes http
func New(cfg config.Config, svc service.CouponService) *API {
	router := SetupRouter(cfg)

	api := &API{
		mux: router,
		cfg: cfg,
		svc: svc,
	}
	return api.withServer().withRoutes()
}

// SetupRouter sets gin router according to the environment.
func SetupRouter(cfg config.Config) *gin.Engine {
	env := cfg.Env.Environment

	switch env {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)

	}

	var router *gin.Engine
	if env == config.ProductionEnv {
		router = setupProductionRouter()
	} else {
		router = setupDevRouter()
	}
	return router
}

// production-specific auth
func setupProductionRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://schwarz.es"}, // just an example, adapt to real needs
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return router
}

// development/staging auth
func setupDevRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.Default())

	return router
}

func (a *API) withServer() *API {
	a.srv = &http.Server{
		Addr:              fmt.Sprintf(":%d", a.cfg.Env.Port),
		Handler:           a.mux,
		ReadTimeout:       10 * time.Second, // Max time to read the full request
		WriteTimeout:      10 * time.Second, // Max time to write the response
		IdleTimeout:       20 * time.Second, // Max time for idle Keep-Alive connections
		ReadHeaderTimeout: 3 * time.Second,  // Max time to read request headers
	}
	return a
}

func (a *API) withRoutes() *API {
	authMiddleware := auth.TokenMiddleware(a.cfg.Env.AuthConfig.JWTSecret)

	apiGroup := a.mux.Group("/api")
	apiGroup.Use(authMiddleware)

	// Admin-only endpoints
	adminGroup := apiGroup.Group("")
	adminGroup.Use(auth.RequireRoles(auth.RoleAdmin))
	{
		adminGroup.POST("/coupon", a.CreateCoupon)
	}

	// Endpoints that both users and admins can access
	userGroup := apiGroup.Group("")
	userGroup.Use(auth.RequireRoles(auth.RoleUser, auth.RoleAdmin))
	{
		userGroup.POST("/coupon/validation", a.ApplyCoupon)
		userGroup.GET("/coupons", a.GetCoupons)
	}
	return a
}

func (a *API) Start() error {
	log.Printf("Starting the service on port: %v", a.cfg.Env.Port)
	log.Printf("Running in %s mode", a.cfg.Env.Environment)
	return a.srv.ListenAndServe()
}

func (a *API) Shutdown() {
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer shutdownCancel()

	err := a.srv.Shutdown(shutdownCtx)
	if err != nil {
		log.Println(err)
		if errors.Is(err, context.DeadlineExceeded) {
			if err := a.srv.Close(); err != nil {
				log.Fatalf("Forced service close failed: %v\n", err)
			}
		}
	}
}
