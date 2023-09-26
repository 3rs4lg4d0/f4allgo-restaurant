package boot

import (
	"fmt"

	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"
	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/avito-tech/go-transaction-manager/trm/settings"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDatabaseConnection() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		GetConfig().DbHost,
		GetConfig().DbUser,
		GetConfig().DbPassword,
		GetConfig().DbName,
		GetConfig().DbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		panic("failed to connect to database")
	}

	return db
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
