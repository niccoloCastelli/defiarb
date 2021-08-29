package storage

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/niccoloCastelli/defiarb/config"
	"github.com/niccoloCastelli/defiarb/storage/models"
)

func NewDb(conf *config.DbConfig) (*gorm.DB, error) {
	return gorm.Open("postgres", conf.ConnectionString())
}

func Migrate(db *gorm.DB) error {
	tables := []interface{}{
		models.Token{},
		models.LiquidityPool{},
		models.Transaction{},
	}
	return db.AutoMigrate(tables...).Error
}
