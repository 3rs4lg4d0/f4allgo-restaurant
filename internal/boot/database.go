package boot

import (
	"database/sql"
	"fmt"

	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"
	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/avito-tech/go-transaction-manager/trm/settings"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDatabaseConnection() (*gorm.DB, *sql.DB) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		GetConfig().DbHost,
		GetConfig().DbUser,
		GetConfig().DbPassword,
		GetConfig().DbName,
		GetConfig().DbPort)

	var ll logger.LogLevel
	if GetConfig().LogLevel < 3 {
		ll = logger.Warn
	} else {
		ll = logger.Silent
	}
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(ll),
	})
	if err != nil {
		panic("failed to connect to database")
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		panic("failed to get *sql.DB from *gorm.DB")
	}

	sqlDB.SetMaxOpenConns(GetConfig().DbMaxOpenConns)
	sqlDB.SetMaxIdleConns(GetConfig().DbMaxIdleConns)
	sqlDB.SetConnMaxIdleTime(GetConfig().DbConnMaxIdleTime)
	sqlDB.SetConnMaxLifetime(GetConfig().DbConnMaxLifetime)

	return gormDB, sqlDB
}

func GetTransactionManager(db *gorm.DB) *manager.Manager {
	return manager.Must(
		trmgorm.NewDefaultFactory(db),
		manager.WithSettings(trmgorm.MustSettings(
			settings.Must(
				settings.WithPropagation(trm.PropagationRequired))),
		),
	)
}
