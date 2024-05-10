package commands

import (
	"errors"
	"fmt"
	"time"

	"github.com/yimincai/gopunch/internal/bot"
	"github.com/yimincai/gopunch/internal/errs"
	"github.com/yimincai/gopunch/internal/service"
	"github.com/yimincai/gopunch/pkg/logger"
	"github.com/yimincai/gopunch/pkg/utils"
	"gorm.io/gorm"
)

type CommandPunchOnTime struct {
	Svc service.Service
}

func (c *CommandPunchOnTime) IsAdminRequired() bool {
	return false
}

func (c *CommandPunchOnTime) Invokes() []string {
	return []string{"PunchOnTime", "pot", "punchontime", "punchOnTime"}
}

func (c *CommandPunchOnTime) Description() string {
	return "Punch on given time"
}

func (c *CommandPunchOnTime) Exec(ctx *bot.Context) (err error) {
	if len(ctx.Args) != 1 {
		usage := "Usage: %sPunchOnTime <time>, format: YYYY-MM-DD HH:MM:SS"
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, usage)
		if err != nil {
			return errs.ErrSendingMessage
		}
	}

	t, err := time.Parse("2006/01/02 15:04:05 ", ctx.Args[1])
	if err != nil {
		response := "❌ Invalid time format, please use `YYYY-MM-DD HH:MM:SS`"
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
		if err != nil {
			return errs.ErrSendingMessage
		}
		return
	}

	// check if the user is already registered
	user, err := c.Svc.Repo.GetUserByDiscordUserID(ctx.Message.Author.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := fmt.Sprintf("You are not registered yet, please register first using `%sRegister <account> <password>`", c.Svc.Cfg.Prefix)
			_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
			if err != nil {
				return errs.ErrSendingMessage
			}
			return
		} else {
			return errs.ErrInternalError
		}
	}

	accessToken, err := c.Svc.Login(ctx.Message.Author.ID)
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
