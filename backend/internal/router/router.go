package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/config"
	"github.com/gbstudyapply/gbstudyapply/internal/dto"
	"github.com/gbstudyapply/gbstudyapply/internal/handler"
	"github.com/gbstudyapply/gbstudyapply/internal/middleware"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
	"github.com/gbstudyapply/gbstudyapply/internal/service"
	"github.com/gbstudyapply/gbstudyapply/internal/util"
)

// Setup builds the gin engine.
func Setup(cfg *config.Config, db *gorm.DB, minio *util.MinIOClient, logger *slog.Logger) *gin.Engine {
	userRepo := repository.NewUserRepository(db)
	univRepo := repository.NewUniversityRepository(db)
	appRepo := repository.NewApplicationProjectRepository(db)
	docRepo := repository.NewDocumentRepository(db)
	verRepo := repository.NewDocumentVersionRepository(db)
	annRepo := repository.NewAnnotationRepository(db)
	matRepo := repository.NewMaterialItemRepository(db)
	tlRepo := repository.NewTimelineNodeRepository(db)
	recRepo := repository.NewRecommendationRepository(db)
	msgRepo := repository.NewMessageRepository(db)

	userService := service.NewUserService(userRepo, logger, cfg)
	univService := service.NewUniversityService(univRepo, logger)
	appService := service.NewApplicationService(appRepo, univRepo, logger)
	docService := service.NewDocumentService(db, docRepo, verRepo, annRepo, appRepo, logger)
	matService := service.NewMaterialService(matRepo, logger)
	tlService := service.NewTimelineService(tlRepo, logger)
	recService := service.NewRecommendationService(recRepo, univRepo, logger)
	msgService := service.NewMessageService(msgRepo, logger)
	_ = service.NewNotificationService(tlRepo, msgRepo, logger)

	userHandler := handler.NewUserHandler(userService, logger)
	univHandler := handler.NewUniversityHandler(univService, logger)
	appHandler := handler.NewApplicationHandler(appService, logger)
	docHandler := handler.NewDocumentHandler(docService, logger)
	matHandler := handler.NewMaterialHandler(matService, logger)
	tlHandler := handler.NewTimelineHandler(tlService, logger)
	recHandler := handler.NewRecommendationHandler(recService, logger)
	msgHandler := handler.NewMessageHandler(msgService, logger)
	uploadHandler := handler.NewUploadHandler(minio, logger)
	dashboardHandler := handler.NewDashboardHandler(appService, univService, matService, logger)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.CORS(cfg))
	r.Use(middleware.ErrorHandler(logger))

	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, dto.OK(gin.H{"status": "ok"})) })

	limiter := middleware.NewRateLimiter(cfg.RateLimitReq, cfg.RateLimitWin)
	v1 := r.Group("/api/v1")
	{
		v1.GET("/files/:key", uploadHandler.Get)
		registerUserRoutes(v1, cfg, userHandler, limiter)
		registerUniversityRoutes(v1, cfg, univHandler, limiter)
		registerApplicationRoutes(v1, cfg, appHandler, docHandler, matHandler, tlHandler, limiter)
		registerDocumentRoutes(v1, cfg, docHandler, limiter)
		registerTimelineRoutes(v1, cfg, tlHandler, limiter)
		registerRecommendationRoutes(v1, cfg, recHandler, limiter)
		registerMessageRoutes(v1, cfg, msgHandler, limiter)
		v1.GET("/dashboard/stats", middleware.AuthRequired(cfg), dashboardHandler.Stats)
		v1.POST("/uploads", middleware.AuthRequired(cfg), limiter.Limit(), uploadHandler.Upload)
	}
	return r
}
