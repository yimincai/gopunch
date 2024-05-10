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

type CommandPunch struct {
	Svc service.Service
}

func (c *CommandPunch) IsAdminRequired() bool {
	return false
}

func (c *CommandPunch) Invokes() []string {
	return []string{"Punch", "punch", "p"}
}

func (c *CommandPunch) Description() string {
	return "Punch NOW"
}

func (c *CommandPunch) Exec(ctx *bot.Context) (err error) {
	user, err := c.Svc.Repo.GetUserByDiscordUserID(ctx.Message.Author.ID)
	if err != nil {
		return errs.ErrUserNotFound
	}

	if !user.IsEnable {
		return errs.ErrUserNotEnabled
	}

	assessToken, err := c.Svc.Login(ctx.Message.Author.ID)
	if err != nil {
		return errs.ErrLoginFailed
	}

	err = c.Svc.WebPunch(assessToken)
	if err != nil {
		return errs.ErrPunchFailed
	}

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf("✅ %s punched successfully at %s", user.Account, utils.TimeFormat(time.Now())))
	if err != nil {
		return errs.ErrSendingMessage
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)
	return
}
