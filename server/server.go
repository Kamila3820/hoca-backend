package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Kamila3820/hoca-backend/config"
	"github.com/Kamila3820/hoca-backend/modules/account/misc"
	_oauth2Controller "github.com/Kamila3820/hoca-backend/modules/oauth2/controller"
	_oauth2Service "github.com/Kamila3820/hoca-backend/modules/oauth2/service"
	_userRepository "github.com/Kamila3820/hoca-backend/modules/user/repository"
	"github.com/Kamila3820/hoca-backend/pkg/databases"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type echoServer struct {
	app  *echo.Echo
	db   databases.Database
	conf *config.Config
}

var (
	once   sync.Once
	server *echoServer
)

func NewEchoServer(conf *config.Config, db databases.Database) *echoServer {
	echoApp := echo.New()
	echoApp.Logger.SetLevel(log.DEBUG)

	once.Do(func() {
		server = &echoServer{
			app:  echoApp,
			db:   db,
			conf: conf,
		}
	})

	return server
}

func (s *echoServer) Start() {
	corsMiddleware := getCORSMiddleware(s.conf.Server.AllowOrigins)
	bodyLimitMiddleware := getBodyLimitMiddleware(s.conf.Server.BodyLimit)
	timeoutMiddleware := getTimeOutMiddleware(s.conf.Server.TimeOut)

	// s.app.Use(middleware.Recover())
	s.app.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  log.ERROR,
	}))
	s.app.Use(middleware.Logger())
	s.app.Use(corsMiddleware)
	s.app.Use(bodyLimitMiddleware)
	s.app.Use(timeoutMiddleware)

	// authorizingMiddleware := s.getAuthorizingMiddleware()

	s.app.GET("/v1/health", s.healthCheck)
	s.app.GET("/v1/panic", func(c echo.Context) error {
		panic("Panic!")
	})

	// s.initOAuth2Router()
	s.initAccountRouter()
	s.initPostRouter()
	s.initUserRatingRouter()
	s.initOrderRouter()
	s.initHistoryRouter()
	s.initNotificationRouter()

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)
	go s.gracefullyShutdown(quitCh)

	s.httpListening()
}

func (s *echoServer) httpListening() {
	serverURL := fmt.Sprintf(":%d", s.conf.Server.Port)

	if err := s.app.Start(serverURL); err != nil && err != http.ErrServerClosed {
		s.app.Logger.Fatalf("Error: %s", err.Error())
	}

}

func (s *echoServer) gracefullyShutdown(quitCh chan os.Signal) {
	ctx := context.Background()

	<-quitCh
	s.app.Logger.Fatalf("Shutting down server...")

	if err := s.app.Shutdown(ctx); err != nil {
		s.app.Logger.Fatalf("Error: %s", err.Error())
	}
}

func (s *echoServer) healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

// Middleware
func getTimeOutMiddleware(timeout time.Duration) echo.MiddlewareFunc {
	return middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Request Timeout",
		Timeout:      timeout * time.Second,
	})
}

func getCORSMiddleware(allowOrigins []string) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: allowOrigins,
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	})
}

func getBodyLimitMiddleware(bodyLimit string) echo.MiddlewareFunc {
	return middleware.BodyLimit(bodyLimit)
}

func Jwt() echo.MiddlewareFunc {
	config := echojwt.Config{
		SigningKey:  []byte("babycomeandtakemylovenadruinit"),
		TokenLookup: "header:x-auth-token",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(misc.UserClaim)
		},
		ErrorHandler: func(c echo.Context, err error) error {
			fmt.Println("JWT Validation Error:", err)

			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"message": "JWT validation failure",
				"error":   err.Error(),
			})
		},
	}

	return echojwt.WithConfig(config)
}

func (s *echoServer) getAuthorizingMiddleware() *authorizingMiddleware {
	userRepository := _userRepository.NewUserRepositoryImpl(s.db, s.app.Logger)

	oauth2Service := _oauth2Service.NewGoogleOAuth2Service(userRepository)
	oauth2Controller := _oauth2Controller.NewGoogleOAuth2Controller(oauth2Service, s.conf.OAuth2, s.app.Logger)

	return &authorizingMiddleware{
		oauth2Controller: oauth2Controller,
		oauth2Conf:       s.conf.OAuth2,
		logger:           s.app.Logger,
	}

}
