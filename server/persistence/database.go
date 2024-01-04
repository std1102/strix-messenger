package persistence

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strix-server/system"
)

var DatabaseContext *gorm.DB

func InitDb() {
	var err error
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		system.SystemConfig.Db.Host,
		system.SystemConfig.Db.Username,
		system.SystemConfig.Db.Password,
		system.SystemConfig.Db.Database,
		system.SystemConfig.Db.Port,
		system.SystemConfig.Db.TimeZone,
	)
	DatabaseContext, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		system.Logger.Fatal("Cannot connect to database", err)
	}
	system.Logger.Infof("Connected to database %s@%s:%s", system.SystemConfig.Db.Database, system.SystemConfig.Db.Host, system.SystemConfig.Db.Port)
	if system.SystemConfig.Db.Mirgate {
		migrate()
	}
}

func migrate() {
	system.Logger.Info("Creating database")
	DatabaseContext.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	_migrate(User{})
	_migrate(PreKeys{})
	_migrate(Device{})
	_migrate(ChatSession{})
	_migrate(PendingMessage{})
	_migrate(UploadedFile{})
}

func _migrate(model interface{}) {
	err := DatabaseContext.AutoMigrate(&model)
	if err != nil {
		system.Logger.Fatal("Cannot migrate database")
	}
}
