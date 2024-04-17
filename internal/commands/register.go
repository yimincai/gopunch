package commands

import (
	"fmt"

	"github.com/yimincai/gopunch/domain"
	"github.com/yimincai/gopunch/internal/bot"
	"github.com/yimincai/gopunch/internal/errs"
	"github.com/yimincai/gopunch/internal/service"
	"github.com/yimincai/gopunch/pkg/logger"
)

type CommandRegister struct {
	Svc service.Service
}

func (c *CommandRegister) Invokes() []string {
	return []string{"Register", "register", "rg", "r"}
}

func (c *CommandRegister) Description() string {
	return "Register user to the system"
}

func (c *CommandRegister) Exec(ctx *bot.Context) (err error) {
	// check if the user is already registered
	existUser, err := c.Svc.Repo.GetUserByDiscordUserID(ctx.Message.Author.ID)
	if err == nil {
		response := fmt.Sprintf("You are already registered as %s \nif you want to overwrite it, please using `%sForceRegister <account> <password>`", existUser.Account, c.Svc.Cfg.Prefix)
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
		if err != nil {
			return errs.ErrSendingMessage
		}
		return
	}

	if len(ctx.Args) != 2 {
		usage := fmt.Sprintf("Usage: %sRegister <account> <password>", c.Svc.Cfg.Prefix)
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, usage)
		if err != nil {
			return errs.ErrSendingMessage
		}
		return
	}

	user := &domain.User{
		Nickname:      ctx.Message.Author.Username,
		Account:       ctx.Args[0],
		Password:      ctx.Args[1],
		DiscordUserID: ctx.Message.Author.ID,
		IsEnable:      true,
	}

	_, err = c.Svc.TryToLogin(user.Account, user.Password)
	if err != nil {
		return errs.ErrLoginVerifyFailed
	}

	err = c.Svc.Register(user)
	if err != nil {
		return errs.ErrInternalError
	}

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Registered successfully, You are now able to use the system")
	if err != nil {
		return errs.ErrSendingMessage
	}

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf("You are registered as %s", user.Account))
	if err != nil {
		return errs.ErrSendingMessage
	}

	response := "```"
	response += "系統會自動在需打卡的第一個時間前30分鐘內隨機打卡，並在需打卡的最後一個時間後30分鐘內隨機打卡。\n"
	response += fmt.Sprintf("國定假日尚未實作，請自行使用 %sdayoff 指令調整。\n\n", c.Svc.Cfg.Prefix)
	response += "Monday: 8:00 - 18:00\n"
	response += "Tuesday: 8:00 - 18:00\n"
	response += "Wednesday: 8:00 - 18:00\n"
	response += "Thursday: 8:00 - 18:00\n"
	response += "Friday: 8:00 - 18:00\n"
	response += "Saturday: None\n"
	response += "Sunday: None\n"
	response += "```"

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
	if err != nil {
		return errs.ErrSendingMessage
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)
	return
}
