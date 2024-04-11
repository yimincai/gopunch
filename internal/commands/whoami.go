package commands

import (
	"fmt"

	"github.com/yimincai/gopunch/internal/bot"
	"github.com/yimincai/gopunch/internal/service"
	"github.com/yimincai/gopunch/pkg/logger"
)

type CommandWhoAmI struct {
	Svc service.Service
}

func (c *CommandWhoAmI) Invokes() []string {
	return []string{"WhoAmI", "whoami", "whoAmI", "w"}
}

func (c *CommandWhoAmI) Description() string {
	return "Returns the user's account name and discord user id"
}

func (c *CommandWhoAmI) Exec(ctx *bot.Context) (err error) {
	user, err := c.Svc.Repo.GetUserByDiscordUserID(ctx.Message.Author.ID)
	if err != nil {
		logger.Errorf("Error getting user: %s", err)
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Error getting user, are you registered?")
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
			return
		}
		return
	}

	if !user.IsEnable {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "User is disabled!!")
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
			return err
		}
		return
	}

	response := fmt.Sprintf("Your account name is: `%s` and your discord user id is: `%s`", user.Account, ctx.Message.Author.ID)
	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
	if err != nil {
		logger.Errorf("Error sending message: %s", err)
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)
	return
}
