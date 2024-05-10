package commands

import (
	"github.com/yimincai/gopunch/internal/bot"
	"github.com/yimincai/gopunch/internal/config"
	"github.com/yimincai/gopunch/internal/errs"
	"github.com/yimincai/gopunch/pkg/logger"
)

type CommandHelp struct {
	Cfg *config.Config
}

func (c *CommandHelp) IsAdminRequired() bool {
	return false
}

func (c *CommandHelp) Invokes() []string {
	return []string{"Help", "h", "help"}
}

func (c *CommandHelp) Description() string {
	return "Show help message"
}

func (c *CommandHelp) Exec(ctx *bot.Context) (err error) {
	response := "```"
	for _, command := range ctx.Commands {
		response += "Command: " + c.Cfg.Prefix + command.Invokes()[0] + ": "
		response += command.Description() + "\n"
	}
	response += "```"

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
	if err != nil {
		return errs.ErrSendingMessage
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)
	return
}
