package bot

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	"github.com/yimincai/gopunch/internal/config"
	"github.com/yimincai/gopunch/internal/database"
	"github.com/yimincai/gopunch/internal/service"
	"github.com/yimincai/gopunch/pkg/logger"
	"github.com/yimincai/gopunch/repository"
)

type Bot struct {
	Session *discordgo.Session
	Cfg     *config.Config
	Repo    repository.Repository
	Svc     service.Service
	Cron    *cron.Cron
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

	l := logger.GetInstance()
	c := cron.New(cron.WithLogger(l), cron.WithChain(cron.Recover(l)))
	s := service.NewService(cfg, repo, session, c)

	return &Bot{
		Session: session,
		Svc:     s,
		Cfg:     cfg,
		Repo:    repo,
		Cron:    c,
	}
}

func (b *Bot) Run() {
	err := b.Svc.InitSchedules()
	if err != nil {
		logger.Panicf("Error scheduling: %v", err)
	}

	b.Session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMembers |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsDirectMessages)

	err = b.Session.Open()
	if err != nil {
		panic("Error opening connection to Discord: " + err.Error())
	}

	logger.Infof("Bot Login as %s, UserID: %s", b.Session.State.User.String(), b.Session.State.User.ID)
	logger.Info("Bot is now running. Press CTRL-C to exit.")
}

func (b *Bot) Close() {
	b.Session.Close()
}
