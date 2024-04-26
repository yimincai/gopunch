package commands

import (
	"errors"
	"fmt"

	"github.com/yimincai/gopunch/internal/bot"
	"github.com/yimincai/gopunch/internal/errs"
	"github.com/yimincai/gopunch/internal/service"
	"github.com/yimincai/gopunch/pkg/logger"
	"gorm.io/gorm"
)

type CommandUpdateAccount struct {
	Svc service.Service
}

func (c *CommandUpdateAccount) Invokes() []string {
	return []string{"UpdateAccount", "ua", "updateaccount", "updateAccount"}
}

func (c *CommandUpdateAccount) Description() string {
	return "Update punch service's account and password"
}

func (c *CommandUpdateAccount) Exec(ctx *bot.Context) (err error) {
	if len(ctx.Args) != 2 {
		usage := fmt.Sprintf("Usage: %sForceRegister <account> <password>", c.Svc.Cfg.Prefix)
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, usage)
		if err != nil {
			return errs.ErrSendingMessage
		}
		return
	}

	user, err := c.Svc.Repo.GetUserByDiscordUserID(ctx.Message.Author.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrUserNotFound
		}
		return errs.ErrInternalError
	}

	_, err = c.Svc.TryToLogin(ctx.Args[0], ctx.Args[1])
	if err != nil {
		return errs.ErrLoginVerifyFailed
	}

	user.Account = ctx.Args[0]
	user.Password = ctx.Args[1]

	err = c.Svc.Repo.UpdateAccount(user)
	if err != nil {
		return errs.ErrInternalError
	}

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "✅ Update account and password successfully")
	if err != nil {
		return errs.ErrSendingMessage
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)
	return
}
