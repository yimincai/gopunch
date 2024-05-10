package commands

import (
	"fmt"
	"time"

	"github.com/yimincai/gopunch/internal/bot"
	"github.com/yimincai/gopunch/internal/errs"
	"github.com/yimincai/gopunch/internal/service"
	"github.com/yimincai/gopunch/pkg/logger"
	"github.com/yimincai/gopunch/pkg/utils"
)

type CommandPunchUserOnTime struct {
	Svc service.Service
}

func (c *CommandPunchUserOnTime) IsAdminRequired() bool {
	return true
}

func (c *CommandPunchUserOnTime) Invokes() []string {
	return []string{"PunchUserOnTime", "punchuserontime", "puot"}
}

func (c *CommandPunchUserOnTime) Description() string {
	return "Punch user on given time"
}

func (c *CommandPunchUserOnTime) Exec(ctx *bot.Context) (err error) {
	if len(ctx.Args) != 3 {
		usage := "Usage: %sPunchUserOnTime <Account> <time>, format: YYYY/MM/DD HH:MM:SS"
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, usage)
		if err != nil {
			return errs.ErrSendingMessage
		}
	}

	account := ctx.Args[0]
	inputDate := ctx.Args[1]
	inputTime := ctx.Args[2]

	t, err := time.Parse("2006/01/02 15:04:05 ", inputDate+" "+inputTime)
	if err != nil {
		response := "❌ Invalid time format, please use `YYYY/MM/DD HH:MM:SS`"
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
		if err != nil {
			return errs.ErrSendingMessage
		}
		return
	}

	// find account user
	user, err := c.Svc.Repo.GetUserByAccount(account)
	if err != nil {
		return errs.ErrUserNotFound
	}

	accessToken, err := c.Svc.Login(user.DiscordUserID)
	if err != nil {
		logger.Error(err)
		return
	}

	err = c.Svc.PunchOnTime(accessToken, t)
	if err != nil {
		return errs.ErrPunchOnTimeFailed
	}

	response := fmt.Sprintf("✅ %s punch on time successfully at %s", user.Account, utils.TimeFormat(t))

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
	if err != nil {
		return errs.ErrSendingMessage
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)
	return
}
