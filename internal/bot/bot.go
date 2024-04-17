package bot

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/yimincai/gopunch/internal/config"
	"github.com/yimincai/gopunch/internal/cronjob"
	"github.com/yimincai/gopunch/internal/database"
	"github.com/yimincai/gopunch/internal/service"
	"github.com/yimincai/gopunch/pkg/logger"
	"github.com/yimincai/gopunch/repository"
	"gorm.io/gorm"
)

type Bot struct {
	Session *discordgo.Session
	Cfg     *config.Config
	Repo    repository.Repository
	Svc     service.Service
	Cron    *cronjob.Cron
	Db      *gorm.DB
}

func New() *Bot {
	// create folder if not exists
	if _, err := os.Stat("./bot_files"); os.IsNotExist(err) {
		err := os.Mkdir("./bot_files", os.ModePerm)
		if err != nil {
			logger.Panicf("Error creating bot_files folder: " + err.Error())
		}
	}

	cfg := config.New()
	db, err := database.New(cfg)
	if err != nil {
		logger.Panicf("Error creating database: " + err.Error())
	}
	repo := repository.New(db)

	session, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		logger.Panicf("Error creating Discord session: " + err.Error())
	}

	s := service.NewService(cfg, repo, session)

	return &Bot{
		Session: session,
		Svc:     s,
		Cfg:     cfg,
		Repo:    repo,
		Cron:    cronjob.New(logger.GetInstance(), s),
		Db:      db,
	}
}

func (b *Bot) Run() {
	b.Session.Identify.Intents = discordgo.IntentDirectMessages

	err := b.Session.Open()
	if err != nil {
		panic("Error opening connection to Discord: " + err.Error())
	}

	database.MigrateUserMissingData(b.Db, b.Session)

	logger.Infof("Bot Login as %s, UserID: %s", b.Session.State.User.String(), b.Session.State.User.ID)
	logger.Info("Bot is now running. Press CTRL-C to exit.")
}

func (b *Bot) Close() {
	b.Session.Close()
}
