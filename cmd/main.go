package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/yimincai/gopunch/internal/bot"
	"github.com/yimincai/gopunch/internal/commands"
	"github.com/yimincai/gopunch/internal/events"
	"github.com/yimincai/gopunch/pkg/logger"
)

var (
	goVersion = runtime.Version()
	oSArch    = fmt.Sprintf("%v/%v", runtime.GOOS, runtime.GOARCH)
)

func main() {
	logger.Infof("Go Version: %s", goVersion)
	logger.Infof("OS/Arch: %s", oSArch)
	server := bot.New()

	// Register events
	registerEvents(server)

	// Register commands
	registerCommands(server)

	server.Run()
	defer func() {
		server.Close()
		logger.Info("Bot closed")
	}()

	go func() {
		server.Cron.Start()
		// defer func() {
		// 	server.Cron.Stop()
		// }()
	}()

	// ============================== Graceful shutdown ==============================
	// Wait for interrupt signal to gracefully shut down the server with a timeout.
	quit := make(chan os.Signal, 1)
	// kill (no param) default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catching, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// ========================= Server block and serve here =========================
	q := <-quit
	// ============================== Graceful shutdown ==============================

	logger.Infof("Got signal: %s", q.String())
	logger.Info("Shutdown server ...")
}

func registerEvents(s *bot.Bot) {
	s.Session.AddHandler(events.NewMessageHandler(s.Svc).Handler)
}

func registerCommands(b *bot.Bot) {
	// Register commands here
	cmdHandler := bot.NewCommandHandler(b.Cfg.Prefix)
	cmdHandler.OnError = func(ctx *bot.Context, err error) {
		logger.Errorf("Error executing command: %v", err)
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "⚠️  An error occurred while executing the command. \n❌ error: "+err.Error())
		if err != nil {
			logger.Errorf("Error sending message: %v", err)
		}
	}

	// cmdHandler.RegisterCommand(&commands.CommandPing{})
	cmdHandler.RegisterCommand(&commands.CommandHelp{Cfg: b.Cfg})
	cmdHandler.RegisterCommand(&commands.CommandDefaultSchedule{Svc: b.Svc})
	cmdHandler.RegisterCommand(&commands.CommandHealth{Svc: b.Svc})
	cmdHandler.RegisterCommand(&commands.CommandGetUsers{Svc: b.Svc})
	cmdHandler.RegisterCommand(&commands.CommandPunch{Svc: b.Svc})
	cmdHandler.RegisterCommand(&commands.CommandRegister{Svc: b.Svc})
	cmdHandler.RegisterCommand(&commands.CommandUpdateAccount{Svc: b.Svc})
	cmdHandler.RegisterCommand(&commands.CommandWhoAmI{Svc: b.Svc})
	cmdHandler.RegisterCommand(&commands.CommandDayOff{Svc: b.Svc})
	cmdHandler.RegisterCommand(&commands.CommandSetSchedule{Svc: b.Svc})
	cmdHandler.RegisterCommand(&commands.CommandPrintSchedule{Svc: b.Svc})
	cmdHandler.RegisterCommand(&commands.CommandPunchOnTime{Svc: b.Svc})

	b.Session.AddHandler(cmdHandler.HandleMessage)
}
