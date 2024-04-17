package database

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/yimincai/gopunch/domain"
	"github.com/yimincai/gopunch/internal/config"
	"github.com/yimincai/gopunch/internal/enums"
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

func MigrateUserMissingData(db *gorm.DB, s *discordgo.Session) {
	var allUsers []*domain.User
	result := db.Find(&allUsers)
	if result.Error != nil {
		logger.Errorf("Failed to get all users: %s", result.Error)
		return
	}

	for _, user := range allUsers {
		if user.Nickname == "" {
			dUser, err := s.User(user.DiscordUserID)
			if err != nil {
				logger.Errorf("User %s not found in discord: %s", user.DiscordUserID, err)
				continue
			}

			if user.Account == "AD0017" {
				user.Role = enums.RoleType_Admin
			} else {
				user.Role = enums.RoleType_Normal
			}

			user.Nickname = dUser.Username
			db.Save(user)

			logger.Infof("User %s updated as %s member", user.Account, user.Role)
		}
	}
}
