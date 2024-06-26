package commands

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/yimincai/gopunch/domain"
	"github.com/yimincai/gopunch/internal/bot"
	"github.com/yimincai/gopunch/internal/errs"
	"github.com/yimincai/gopunch/internal/service"
	"github.com/yimincai/gopunch/pkg/logger"
)

type CommandGetUsers struct {
	Svc service.Service
}

func (c *CommandGetUsers) IsAdminRequired() bool {
	return true
}

func (c *CommandGetUsers) Invokes() []string {
	return []string{"GetUsers", "gus", "getusers", "getUsers"}
}

func (c *CommandGetUsers) Description() string {
	return "Get all users in db"
}

func (c *CommandGetUsers) Exec(ctx *bot.Context) (err error) {
	var users []*domain.User
	users, err = c.Svc.Repo.FindUsers()
	if err != nil {
		return errs.ErrInternalError
	}

	// a table of users
	t := table.NewWriter()
	t.AppendHeader(table.Row{"ID", "Name", "Account", "Role", "Is_Enable"})

	var response string
	for _, user := range users {
		t.AppendRow([]interface{}{user.ID, user.Name, user.Account, user.Role, user.IsEnable})
	}

	t.SetStyle(table.StyleLight)
	response = "```" + t.Render() + "```"

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
	if err != nil {
		return errs.ErrSendingMessage
	}

	logger.Info("Command Executed: ", c.Invokes(), " UserID: ", ctx.Message.Author.ID)
	return
}
