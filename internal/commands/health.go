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

func (c *CommandHealth) Invokes() []string {
	return []string{"Health", "health"}
}

func (c *CommandHealth) Description() string {
	return "Checking if DM user is registered, and login functionality"
}

func (c *CommandHealth) Exec(ctx *bot.Context) (err error) {
	user, err := c.Svc.Repo.GetUserByDiscordUserID(ctx.Message.Author.ID)
	if err != nil {
		logger.Errorf("Error getting user: %s", err)
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Error getting user, did you registered?")
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
			return
		}
		return
	}

	if !user.IsEnable {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, errs.ErrUserDisabled.Error())
		return
	}

	_, err = c.Svc.Login(ctx.Message.Author.ID)
	if err != nil {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Error logging in, please update your settings")
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
		}
		logger.Errorf("Error logging in: %s", err)
		return
	}

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf("Login as %s successfully, you are good to go!", user.Account))
	if err != nil {
		logger.Errorf("Error sending message: %s", err)
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)
	return
}
