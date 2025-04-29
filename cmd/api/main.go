package main

import (
	"barcode-checker/internal/config"
	"barcode-checker/internal/controller"
	"barcode-checker/internal/middleware"
	"barcode-checker/internal/migration"
	"barcode-checker/internal/repository"
	"barcode-checker/internal/service"
	"barcode-checker/internal/utils"
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	logger := utils.SetupLogger()
	defer logger.Sync()

	gormDB, err := initPostgreSQL(cfg, logger)
	if err != nil {
		logger.Fatal("PostgreSQL initialization failed", zap.Error(err))
	}

	if err := migration.InitializeDatabase(gormDB, cfg, logger); err != nil {
		logger.Fatal("Database initialization failed", zap.Error(err))
	}

	mongoDB, err := initMongoDB(cfg, logger)
	if err != nil {
		logger.Fatal("MongoDB initialization failed", zap.Error(err))
	}
	defer func() {
		if err := mongoDB.Client().Disconnect(context.Background()); err != nil {
			logger.Error("Failed to disconnect MongoDB", zap.Error(err))
		}
	}()

	authController, productController, historyController, adminController := initServicesAndControllers(
		cfg,
		gormDB,
		mongoDB,
		logger,
	)

	startServer(cfg, authController, productController, historyController, adminController, logger)
}

func initPostgreSQL(cfg *config.Config, logger *zap.Logger) (*gorm.DB, error) {
	dsn := cfg.GetPostgresDSN()
	logger.Info("Connecting to PostgreSQL", zap.String("dsn", "host="+cfg.Database.Host+" port="+cfg.Database.Port+" user="+cfg.Database.User+" dbname="+cfg.Database.Name))

	var gormDB *gorm.DB
	var err error

	maxAttempts := 3
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		gormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}

		logger.Warn("Connection attempt failed",
			zap.Int("attempt", attempt),
			zap.Error(err))

		if attempt < maxAttempts {
			time.Sleep(time.Duration(attempt) * time.Second)
		}
	}

	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	logger.Info("Successfully connected to PostgreSQL")
	return gormDB, nil
}

func initMongoDB(cfg *config.Config, logger *zap.Logger) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Database.MongoDBURI))
	if err != nil {
		return nil, err
	}

	if err = mongoClient.Ping(ctx, nil); err != nil {
		return nil, err
	}

	logger.Info("Successfully connected to MongoDB")

	testColl := mongoClient.Database(cfg.Database.MongoDBName).Collection("test_connection")
	if _, err = testColl.InsertOne(ctx, bson.M{"test": time.Now()}); err != nil {
		return nil, err
	}

	logger.Info("Test MongoDB insert successful")
	return mongoClient.Database(cfg.Database.MongoDBName), nil
}

func initServicesAndControllers(
	cfg *config.Config,
	gormDB *gorm.DB,
	mongoDB *mongo.Database,
	logger *zap.Logger,
) (*controller.AuthController, *controller.ProductController, *controller.HistoryController, *controller.AdminController) {
	userRepo := repository.NewUserRepository(gormDB)
	historyRepo := repository.NewHistoryRepository(mongoDB)

	authService := service.NewAuthService(userRepo, cfg.Auth.JWTSecret, cfg.Auth.JWTDuration)
	barcodeChecker := service.NewLocalBarcodeChecker()
	productService := service.NewProductService(barcodeChecker, historyRepo)
	historyService := service.NewHistoryService(historyRepo)

	return controller.NewAuthController(authService),
		controller.NewProductController(productService),
		controller.NewHistoryController(historyService),
		controller.NewAdminController(authService)
}

func startServer(
	cfg *config.Config,
	authController *controller.AuthController,
	productController *controller.ProductController,
	historyController *controller.HistoryController,
	adminController *controller.AdminController,
	logger *zap.Logger,
) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware(logger))

	authGroup := r.Group("/auth")
	authGroup.Use(middleware.RateLimiter(cfg.RateLimit.AuthLimit, cfg.RateLimit.AuthBurst, time.Minute))
	{
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", authController.Login)
	}

	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.RateLimiter(cfg.RateLimit.APILimit, cfg.RateLimit.APIBurst, time.Minute))
	apiGroup.Use(middleware.JWTAuth(cfg.Auth.JWTSecret))
	{
		apiGroup.POST("/check", productController.CheckProduct)
		apiGroup.GET("/history", historyController.GetHistory)
		apiGroup.DELETE("/history/:id", historyController.DeleteHistoryItem)
	}

	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.JWTAuth(cfg.Auth.JWTSecret))
	adminGroup.Use(middleware.AdminOnly())
	{
		adminGroup.GET("/users", adminController.ListUsers)
		adminGroup.DELETE("/users/:id", adminController.DeleteUser)
	}

	logger.Info("Starting server", zap.String("port", cfg.Server.Port))
	if err := r.Run(cfg.Server.Port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
