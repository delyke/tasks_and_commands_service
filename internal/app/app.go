package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/delyke/tasks_and_commands_service/internal/handler/middleware"
	analytics2 "github.com/delyke/tasks_and_commands_service/internal/repository/mysql/analytics"
	task2 "github.com/delyke/tasks_and_commands_service/internal/repository/mysql/task"
	"github.com/delyke/tasks_and_commands_service/internal/repository/mysql/task_comment"
	"github.com/delyke/tasks_and_commands_service/internal/repository/mysql/task_history"
	team2 "github.com/delyke/tasks_and_commands_service/internal/repository/mysql/team"
	"github.com/delyke/tasks_and_commands_service/internal/repository/mysql/team_member"
	mysqlRepo "github.com/delyke/tasks_and_commands_service/internal/repository/mysql/user"
	"github.com/delyke/tasks_and_commands_service/internal/service/analytics"
	"github.com/delyke/tasks_and_commands_service/internal/service/auth"
	"github.com/delyke/tasks_and_commands_service/internal/service/jwt"
	"github.com/delyke/tasks_and_commands_service/internal/service/task"
	"github.com/delyke/tasks_and_commands_service/internal/service/team"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/delyke/tasks_and_commands_service/closer"
	"github.com/delyke/tasks_and_commands_service/internal/config"
	"github.com/delyke/tasks_and_commands_service/internal/handler"
	"github.com/delyke/tasks_and_commands_service/internal/infrastructure/mysql"
	"github.com/delyke/tasks_and_commands_service/internal/infrastructure/redis"
	"github.com/delyke/tasks_and_commands_service/logger"
	api "github.com/delyke/tasks_and_commands_service/pkg/openapi/tasks/v1"
)

type App struct {
	httpServer *http.Server

	// Infrastructure
	mysqlClient *mysql.Client
	redisClient *redis.Client

	// Services
	authService      *auth.Service
	teamService      *team.Service
	taskService      *task.Service
	analyticsService *analytics.Service

	// JWT
	jwtGenerator *jwt.Generator
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		return a.runHTTPServer(ctx)
	})

	return g.Wait()
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("HTTP server listening on %s", config.AppConfig().HTTP.Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initLogger,
		a.initCloser,
		a.initMySQL,
		a.initRedis,
		a.initServices,
		a.initHTTPServer,
	}
	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initMySQL(ctx context.Context) error {
	client, err := mysql.NewClient(config.AppConfig().MySQL)
	if err != nil {
		return fmt.Errorf("failed to create mysql client: %w", err)
	}

	if err := client.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping mysql: %w", err)
	}

	a.mysqlClient = client

	closer.AddNamed("mysql", func(ctx context.Context) error {
		return a.mysqlClient.Close()
	})

	logger.Info(ctx, "MySQL connected")
	return nil
}

func (a *App) initRedis(ctx context.Context) error {
	client, err := redis.NewClient(config.AppConfig().Redis)
	if err != nil {
		return fmt.Errorf("failed to create redis client: %w", err)
	}

	if err := client.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}

	a.redisClient = client

	closer.AddNamed("redis", func(ctx context.Context) error {
		return a.redisClient.Close()
	})

	logger.Info(ctx, "Redis connected")
	return nil
}

func (a *App) initServices(_ context.Context) error {
	db := a.mysqlClient.DB()

	// Repositories
	userRepo := mysqlRepo.NewUserRepository(db)
	teamRepo := team2.NewTeamRepository(db)
	memberRepo := team_member.NewTeamMemberRepository(db)
	taskRepo := task2.NewTaskRepository(db)
	historyRepo := task_history.NewTaskHistoryRepository(db)
	commentRepo := task_comment.NewTaskCommentRepository(db)
	analyticsRepo := analytics2.NewAnalyticsRepository(db)

	// JWT Generator
	a.jwtGenerator = jwt.NewJWTGenerator(config.AppConfig().JWT.SecretKey())

	// Services
	a.authService = auth.NewAuthService(userRepo, a.jwtGenerator, config.AppConfig().JWT)
	a.teamService = team.NewTeamService(teamRepo, memberRepo, userRepo)
	a.taskService = task.NewTaskService(taskRepo, historyRepo, commentRepo, memberRepo, a.redisClient)
	a.analyticsService = analytics.NewAnalyticsService(analyticsRepo)

	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initHTTPServer(_ context.Context) error {
	db := a.mysqlClient.DB()
	userRepo := mysqlRepo.NewUserRepository(db)

	h := handler.NewHandler(a.authService, a.teamService, a.taskService, userRepo)
	secHandler := middleware.NewSecurityHandler(a.jwtGenerator)

	srv, err := api.NewServer(h, secHandler)
	if err != nil {
		return fmt.Errorf("failed to create ogen server: %w", err)
	}

	a.httpServer = &http.Server{
		Addr:        config.AppConfig().HTTP.Address(),
		Handler:     srv,
		ReadTimeout: config.AppConfig().HTTP.ReadTimeout(),
	}

	closer.AddNamed("http-server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})

	return nil
}
