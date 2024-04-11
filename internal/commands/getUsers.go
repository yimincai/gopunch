package commands

import (
	"github.com/yimincai/gopunch/domain"
	"github.com/yimincai/gopunch/internal/bot"
	"github.com/yimincai/gopunch/internal/errs"
	"github.com/yimincai/gopunch/internal/service"
	"github.com/yimincai/gopunch/pkg/logger"
)

type CommandGetUsers struct {
	Svc service.Service
}

func (c *CommandGetUsers) Invokes() []string {
	return []string{"GetUsers", "gus", "getusers", "getUsers"}
}

func (c *CommandGetUsers) Description() string {
	return "Get all users in db"
}

func (c *CommandGetUsers) Exec(ctx *bot.Context) (err error) {
	user, err := c.Svc.Repo.GetUserByDiscordUserID(ctx.Message.Author.ID)
	if err != nil {
		return errs.ErrUserNotFound
	}

	if !user.IsEnable {
		return errs.ErrUserNotEnabled
	}

	var users []*domain.User
	users, err = c.Svc.Repo.GetUsers()
	if err != nil {
		return errs.ErrInternalError
	}

	var response string
	for _, user := range users {
		response += user.Account + "\n"
	}

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
	if err != nil {
		return errs.ErrSendingMessage
	}

	logger.Info("Command Executed: ", c.Invokes(), " UserID: ", ctx.Message.Author.ID)
	return
}
