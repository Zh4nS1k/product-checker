package migration

import (
	"barcode-checker/internal/config"
	"barcode-checker/internal/model"
	"barcode-checker/internal/repository"
	_ "barcode-checker/internal/utils"
	"errors"
	"fmt"
	"go.uber.org/zap"
	_ "log"
	"os"
	_ "time"

	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	defaultAdminUsername = "admin"
	defaultAdminEmail    = "admin@example.com"
	defaultAdminPassword = "SecureAdminPassword123!"
)

func RunPostgresMigrations(db *gorm.DB, cfg *config.Config, logger *zap.Logger) error {
	if _, err := os.Stat(cfg.Migration.Dir); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", cfg.Migration.Dir)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %v", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	if err := goose.Up(sqlDB, cfg.Migration.Dir); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	logger.Info("PostgreSQL migrations applied successfully")
	return nil
}

func CreateCustomMigrations(db *gorm.DB, logger *zap.Logger) error {
	userRepo := repository.NewUserRepository(db)

	admin, err := userRepo.GetByEmail(defaultAdminEmail)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check for existing admin: %v", err)
	}

	if admin != nil {
		logger.Info("Default admin already exists, skipping creation")
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(getAdminPassword()),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	adminUser := &model.User{
		Username: defaultAdminUsername,
		Email:    defaultAdminEmail,
		Password: string(hashedPassword),
		Role:     "admin",
	}

	if err := userRepo.Create(adminUser); err != nil {
		return fmt.Errorf("failed to create admin user: %v", err)
	}

	logger.Info("Default admin user created successfully")
	return nil
}

func getAdminPassword() string {
	if pass := os.Getenv("DEFAULT_ADMIN_PASSWORD"); pass != "" {
		return pass
	}
	return defaultAdminPassword
}

func InitializeDatabase(db *gorm.DB, cfg *config.Config, logger *zap.Logger) error {
	if err := RunPostgresMigrations(db, cfg, logger); err != nil {
		return fmt.Errorf("database initialization failed: %v", err)
	}

	if err := CreateCustomMigrations(db, logger); err != nil {
		return fmt.Errorf("custom migrations failed: %v", err)
	}

	return nil
}
