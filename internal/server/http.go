package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"spotsync/internal/apperror"
	"spotsync/internal/auth"
	"spotsync/internal/config"
	"spotsync/internal/domain/reservation"
	"spotsync/internal/domain/user"
	"spotsync/internal/domain/zone"
	"spotsync/internal/httpresponse"
	"spotsync/internal/middlewares"

	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}

func Start(db *gorm.DB, cfg *config.Config) {
	if cfg.Environment == "development" {
		db.AutoMigrate(&user.User{}, &zone.ParkingZone{}, &reservation.Reservation{})
	}

	e := echo.New()

	e.HTTPErrorHandler = func(c *echo.Context, err error) {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			var errDetails any
			if appErr.Err != nil {
				errDetails = appErr.Err.Error()
			}
			_ = httpresponse.Error(c, appErr.Code, appErr.Message, errDetails)
			return
		}

		var he *echo.HTTPError
		if errors.As(err, &he) {
			_ = httpresponse.Error(c, he.Code, fmt.Sprintf("%v", he.Message), nil)
			return
		}

		_ = httpresponse.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
	}

	// global validator which is validate http request body
	e.Validator = &CustomValidator{validator: validator.New()}

	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	e.Use(middleware.RequestLoggerWithConfig(
		middleware.RequestLoggerConfig{
			LogURI:       true,
			LogMethod:    true,
			LogStatus:    true,
			LogLatency:   true,
			LogRemoteIP:  true,
			LogUserAgent: true,

			LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
				logger.Info("request",
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"latency", v.Latency,
					"remote_ip", v.RemoteIP,
					"user_agent", v.UserAgent,
				)
				return nil
			},
		},
	))
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
	}))

	e.GET("/health", func(c *echo.Context) error {
		return c.String(http.StatusOK, "The app is healthy!")
	})

	// Setup dependencies
	jwtService, err := auth.NewJWTService(cfg.JwtSecret)
	if err != nil {
		slog.Error("failed to initialize jwt service", "error", err)
		os.Exit(1)
	}
	authMw := middlewares.AuthMiddleware(jwtService)
	adminMw := middlewares.RequireRole("admin")

	// User domain
	userRepo := user.NewRepository(db)
	userSvc := user.NewService(userRepo, jwtService)
	userHandler := user.NewHandler(userSvc)
	user.RegisterRoutes(e, userHandler)

	// Zone domain
	zoneRepo := zone.NewRepository(db)
	zoneSvc := zone.NewService(zoneRepo)
	zoneHandler := zone.NewHandler(zoneSvc)
	zone.RegisterRoutes(e, zoneHandler, authMw, adminMw)

	// Reservation domain
	reservationRepo := reservation.NewRepository(db)
	reservationSvc := reservation.NewService(reservationRepo, zoneRepo)
	reservationHandler := reservation.NewHandler(reservationSvc)
	reservation.RegisterRoutes(e, reservationHandler, authMw, adminMw)

	port := fmt.Sprintf(":%s", cfg.Port)
	if err := e.Start(port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
