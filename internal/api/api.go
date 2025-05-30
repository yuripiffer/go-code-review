package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"coupon_service/internal/service"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Host string
	Port int
}

type API struct {
	srv *http.Server
	mux *gin.Engine
	svc service.CouponService
	cfg Config
}

// New creates a new API instance with the provided configuration and service.
//
// @title My API
// @version 1.0
// @description This is an API for managing coupons.
// @host localhost:8080
// @BasePath /api/
// @schemes http
func New(cfg Config, svc service.CouponService) *API {
	gin.SetMode(gin.ReleaseMode)
	r := new(gin.Engine)
	r = gin.New()
	r.Use(gin.Recovery())

	api := &API{
		mux: r,
		cfg: cfg,
		svc: svc,
	}
	return api.withServer().withRoutes()
}

func (a *API) withServer() *API {
	a.srv = &http.Server{
		// TODO
		//Addr:              fmt.Sprintf(":%d", a.cfg.Port),
		Addr:              fmt.Sprintf(":%s", "8080"),
		Handler:           a.mux,
		ReadTimeout:       10 * time.Second, // Max time to read the full request
		WriteTimeout:      10 * time.Second, // Max time to write the response
		IdleTimeout:       20 * time.Second, // Max time for idle Keep-Alive connections
		ReadHeaderTimeout: 3 * time.Second,  // Max time to read request headers
	}
	return a
}

func (a *API) withRoutes() *API {
	apiGroup := a.mux.Group("/api")

	apiGroup.POST("/coupon/validation", a.ApplyCoupon)
	apiGroup.POST("/coupon", a.CreateCoupon)
	apiGroup.GET("/coupons", a.GetCoupons)

	return a
}

func (a *API) Start() {
	err := a.srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (a *API) Close() {
	<-time.After(5 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}
