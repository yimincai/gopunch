package commands

import (
	"errors"
	"fmt"

	"github.com/yimincai/gopunch/domain"
	"github.com/yimincai/gopunch/internal/bot"
	"github.com/yimincai/gopunch/internal/errs"
	"github.com/yimincai/gopunch/internal/service"
	"github.com/yimincai/gopunch/pkg/logger"
	"gorm.io/gorm"
)

type CommandForceRegister struct {
	Svc service.Service
}

func (c *CommandForceRegister) Invokes() []string {
	return []string{"ForceRegister", "fr", "forceregister", "forceRegister"}
}

func (c *CommandForceRegister) Description() string {
	return "Delete exists user and register again"
}

func (c *CommandForceRegister) Exec(ctx *bot.Context) (err error) {
	if len(ctx.Args) != 2 {
		usage := fmt.Sprintf("Usage: %sForceRegister <account> <password>", c.Svc.Cfg.Prefix)
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, usage)
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
			return err
		}
		return
	}

	existsUser, err := c.Svc.Repo.GetUserByDiscordUserID(ctx.Message.Author.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Error getting user, did you registered?")
			if err != nil {
				logger.Errorf("Error sending message: %s", err)
				return
			}
			return
		}
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, errs.ErrInternalError.Error())
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
			return
		}
		return
	}

	err = c.Svc.Repo.DeleteUserByAccount(existsUser.Account)
	if err != nil {
		logger.Errorf("Error deleting user: %s", err)
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, errs.ErrInternalError.Error())
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
		}
		return
	}

	user := &domain.User{
		Account:       ctx.Args[0],
		Password:      ctx.Args[1],
		DiscordUserID: ctx.Message.Author.ID,
		IsEnable:      true,
	}

	_, err = c.Svc.TryToLogin(user.Account, user.Password)
	if err != nil {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Verify failed, please check your account and password")
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
			return err
		}
		return
	}

	err = c.Svc.Register(user)
	if err != nil {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, err.Error())
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
			return err
		}
		return
	}

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Registered successfully, You are now able to use the system")
	if err != nil {
		logger.Errorf("Error sending message: %s", err)
		return
	}

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf("You are registered as %s", user.Account))
	if err != nil {
		logger.Errorf("Error sending message: %s", err)
		return
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
		logger.Errorf("Error sending message: %s", err)
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)
	return
}
