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
	if _, err := os.Stat("./bot_files/db"); os.IsNotExist(err) {
		os.Mkdir("./bot_files/db", os.ModePerm)
	}

	db, err := gorm.Open(sqlite.Open("./bot_files/db/bot.db"), &gorm.Config{
		Logger: l,
	})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&domain.User{}, &domain.DayOff{})

	logger.Info("Database connected")

	// InitUsers(db, env)

	return db, nil
}

// func InitUsers(db *gorm.DB, env *config.Config) error {
// 	for _, user := range env.Users {
// 		var u domain.User
// 		result := db.Where("account = ?", user.Account).First(&u)
// 		if result.Error != nil {
// 			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 				if user.Account == "" || user.Password == "" || user.DiscordUserID == "" {
// 					logger.Warnf("User information is not complete, skipped: %s", user.Account)
// 					continue
// 				}
// 				result := db.Create(&domain.User{
// 					Account:       user.Account,
// 					Password:      user.Password,
// 					DiscordUserID: user.DiscordUserID,
// 					IsEnable:      true,
// 				})
// 				if result.Error != nil {
// 					logger.Errorf("Init database error: %s", result.Error)
// 					return result.Error
// 				}
//
// 				logger.Infof("User created: %s", user.Account)
// 				continue
// 			} else {
// 				logger.Errorf("Init database error: %s", result.Error)
// 			}
// 		}
//
// 		logger.Infof("User already exists: %s", user.Account)
// 	}
//
// 	return nil
// }
