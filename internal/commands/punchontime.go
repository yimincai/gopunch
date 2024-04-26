package commands

import (
	"github.com/yimincai/gopunch/internal/bot"
	"github.com/yimincai/gopunch/internal/config"
	"github.com/yimincai/gopunch/internal/errs"
	"github.com/yimincai/gopunch/pkg/logger"
)

type CommandPunchOnTime struct {
	Cfg *config.Config
}

func (c *CommandPunchOnTime) Invokes() []string {
	return []string{"PunchOnTime", "pot", "punchontime", "punchOnTime"}
}

func (c *CommandPunchOnTime) Description() string {
	return "Punch on given time"
}

func (c *CommandPunchOnTime) Exec(ctx *bot.Context) (err error) {
	response := "```"
	response += "PunchOnTime\n"
	response += "Not implemented yet\n"
	response += "```"

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
	if err != nil {
		return errs.ErrSendingMessage
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)
	return
}
