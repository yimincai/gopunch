package database

import (
	"log"
	"os"

	"github.com/yimincai/gopunch/domain"
	"github.com/yimincai/gopunch/internal/config"
	"github.com/yimincai/gopunch/pkg/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gLog "gorm.io/gorm/logger"
)

func New(env *config.Config) (*gorm.DB, error) {
	l := gLog.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gLog.Config{
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  gLog.Info,
		},
	)

	// create folder if not exists
	if _, err := os.Stat("./bot_files"); os.IsNotExist(err) {
		err := os.Mkdir("./bot_files", os.ModePerm)
		if err != nil {
			logger.Panicf("Error creating bot_files folder: " + err.Error())
		}
	}

	db, err := gorm.Open(sqlite.Open("./bot_files/bot.db"), &gorm.Config{
		Logger: l,
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&domain.User{}, &domain.DayOff{})
	if err != nil {
		return nil, err
	}

	logger.Info("Database connected")

	// InitUsers(db, env)

	return db, nil
}
