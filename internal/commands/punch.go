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

func (c *CommandPunch) Invokes() []string {
	return []string{"Punch", "punch", "p"}
}

func (c *CommandPunch) Description() string {
	return "Punch NOW"
}

func (c *CommandPunch) Exec(ctx *bot.Context) (err error) {
	user, err := c.Svc.Repo.GetUserByDiscordUserID(ctx.Message.Author.ID)
	if err != nil {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, errs.ErrUserNotFound.Error())
		return
	}

	if !user.IsEnable {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, errs.ErrUserDisabled.Error())
		return
	}

	assessToken, err := c.Svc.Login(ctx.Message.Author.ID)
	if err != nil {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, errs.ErrLoginFailed.Error())
		return
	}

	err = c.Svc.Punch(assessToken)
	if err != nil {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, errs.ErrPunchFailed.Error())
		return err
	}

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf("Punched successfully at %s", utils.TimeFormat(time.Now())))
	if err != nil {
		logger.Errorf("Error sending message: %s", err)
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)
	return
}
