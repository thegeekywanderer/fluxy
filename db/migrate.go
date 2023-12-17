package database

import "github.com/thegeekywanderer/fluxy/models"

// Migrate handles migrations
func Migrate(db *DB) {
	var migrationModels = []interface{}{&models.Client{}}
	err := db.Database.AutoMigrate(migrationModels...)
	if err != nil {
		return
	}
}
