package commands

import (
	"fmt"

	"github.com/yimincai/gopunch/internal/bot"
	"github.com/yimincai/gopunch/internal/errs"
	"github.com/yimincai/gopunch/internal/service"
	"github.com/yimincai/gopunch/pkg/logger"
)

type CommandHealth struct {
	Svc service.Service
}

func (c *CommandHealth) IsAdminRequired() bool {
	return false
}

func (c *CommandHealth) Invokes() []string {
	return []string{"Health", "health"}
}

func (c *CommandHealth) Description() string {
	return "Checking if DM user is registered, and login functionality"
}

func (c *CommandHealth) Exec(ctx *bot.Context) (err error) {
	user, err := c.Svc.Repo.FindUserByDiscordUserID(ctx.Message.Author.ID)
	if err != nil {
		return errs.ErrUserNotFound
	}

	if !user.IsEnable {
		return errs.ErrUserNotEnabled
	}

	_, err = c.Svc.Login(ctx.Message.Author.ID)
	if err != nil {
		return errs.ErrLoginFailed
	}

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf("Login as %s successfully, you are good to go!", user.Account))
	if err != nil {
		return errs.ErrSendingMessage
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)
	return
}
