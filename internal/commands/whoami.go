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

type CommandWhoAmI struct {
	Svc service.Service
}

func (c *CommandWhoAmI) IsAdminRequired() bool {
	return false
}

func (c *CommandWhoAmI) Invokes() []string {
	return []string{"WhoAmI", "whoami", "whoAmI", "w"}
}

func (c *CommandWhoAmI) Description() string {
	return "Returns the user's account name and discord user id"
}

func (c *CommandWhoAmI) Exec(ctx *bot.Context) (err error) {
	user, err := c.Svc.Repo.FindUserByDiscordUserID(ctx.Message.Author.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrUserNotFound
		}
		return errs.ErrInternalError
	}

	if !user.IsEnable {
		return errs.ErrUserNotEnabled
	}

	response := fmt.Sprintf("Your account name is: `%s` and your discord user id is: `%s`", user.Account, ctx.Message.Author.ID)
	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
	if err != nil {
		return errs.ErrSendingMessage
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)
	return
}
