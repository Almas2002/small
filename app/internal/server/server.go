package server

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/signal"
	"small/internal/config"
	productHttp "small/internal/modules/product/delivery/http"
	productRepo "small/internal/modules/product/repository"
	productUc "small/internal/modules/product/usecase"
	userHttp "small/internal/modules/user/delivery/http"
	userRepo "small/internal/modules/user/repository"
	userUc "small/internal/modules/user/usecase"
	"small/pkg/store/postgres"
	"small/pkg/tracing"
	"small/pkg/type/logger"
	"syscall"
)

type server struct {
	log  logger.Logger
	cfg  *config.Config
	pg   *pgxpool.Pool
	v    *validator.Validate
	echo *echo.Echo
}

func New(log logger.Logger, cfg *config.Config) *server {
	return &server{log: log, cfg: cfg, echo: echo.New(), v: validator.New()}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	pgxConn, err := postgres.NewPgxConn(s.cfg.Postgres)
	if err != nil {
		return errors.Wrap(err, "postgresql.NewPgxConn")
	}

	s.log.InfoF("postgres connected: %v", pgxConn.Stat().TotalConns())

	err = s.migrations(pgxConn)
	if err != nil {
		return errors.Wrap(err, "migrations")
	}
	// init repositories
	userRepository := userRepo.New(pgxConn, userRepo.Options{}, s.log)
	productRepository := productRepo.New(pgxConn, productRepo.Options{}, s.log)

	// init useCases
	userUseCase := userUc.New(s.log, userRepository)
	productUseCae := productUc.New(s.log, productRepository, userUseCase, s.cfg)

	userHandlers := userHttp.NewHandlers(s.echo.Group(s.cfg.Http.BaseUserPath), s.log, s.cfg, userUseCase, s.v)
	userHandlers.MapRoutes()

	productHandlers := productHttp.NewHandlers(s.echo.Group(s.cfg.Http.BaseProductPath), s.log, s.cfg, productUseCae, s.v)
	productHandlers.MapRoutes()

	if s.cfg.Jaeger.Enable {
		tracer, closer, err := tracing.NewJaegerTracer(s.cfg.Jaeger)
		if err != nil {
			return err
		}
		defer func(closer io.Closer) {
			err := closer.Close()
			if err != nil {
				s.log.WarnMsg("tracing.Close", err)
				return
			}
		}(closer)
		opentracing.SetGlobalTracer(tracer)
	}

	go func() {
		if err := s.runHttpServer(); err != nil {
			s.log.Errorf(" s.runHttpServer: %v", err)
			cancel()
		}
	}()
	s.log.InfoF("small is listening on PORT: %s", s.cfg.Http.Port)

	<-ctx.Done()
	if err := s.echo.Server.Shutdown(ctx); err != nil {
		s.log.WarnMsg("echo.Server.Shutdown", err)
	}
	return nil
}
