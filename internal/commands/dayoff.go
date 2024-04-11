package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/yimincai/gopunch/internal/bot"
	"github.com/yimincai/gopunch/internal/errs"
	"github.com/yimincai/gopunch/internal/service"
	"github.com/yimincai/gopunch/pkg/logger"
	"github.com/yimincai/gopunch/pkg/utils"
	"gorm.io/gorm"
)

type CommandDayOff struct {
	Svc service.Service
}

func (c *CommandDayOff) Invokes() []string {
	return []string{"DayOff", "dayoff", "dayOff"}
}

func (c *CommandDayOff) Description() string {
	return "Set the day not to punch"
}

func (c *CommandDayOff) Exec(ctx *bot.Context) (err error) {
	if len(ctx.Args) != 1 {
		usage := fmt.Sprintf("Usage: %sDayOff <date>, format example: %sDayOff 2024/12/31", c.Svc.Cfg.Prefix, c.Svc.Cfg.Prefix)
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, usage)
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
			return err
		}
		return
	}

	sp := strings.Split(ctx.Args[0], "/")
	if len(sp) != 3 {
		usage := fmt.Sprintf("Usage: %sDayOff <date>, format example: %sDayOff 2024/12/31", c.Svc.Cfg.Prefix, c.Svc.Cfg.Prefix)
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, usage)
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
			return err
		}
		return
	}

	user, err := c.Svc.Repo.GetUserByDiscordUserID(ctx.Message.Author.ID)
	if err != nil {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, errs.ErrUserNotFound.Error())
		return
	}

	if !user.IsEnable {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, errs.ErrUserDisabled.Error())
		return
	}

	year := sp[0]
	month := sp[1]
	day := sp[2]

	//case to int
	y, err := strconv.Atoi(year)
	if err != nil {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Invalid date")
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
			return err
		}
		return
	}

	m, err := strconv.Atoi(month)
	if err != nil {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Invalid date")
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
			return err
		}
	}

	d, err := strconv.Atoi(day)
	if err != nil {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Invalid date")
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
			return err
		}
	}

	// validate date
	if !utils.IsVaildateDate(y, m, d) {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Invalid date")
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
			return err
		}
		return
	}

	existsDayOff, err := c.Svc.Repo.FindUserDayOffByDate(user.ID, y, m, d)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, errs.ErrInternalError.Error())
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
			return err
		}
		return
	}

	if existsDayOff != nil {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Day off already set")
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
			return err
		}
		return
	}

	err = c.Svc.SetDayOff(ctx.Message.Author.ID, y, m, d)
	if err != nil {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, err.Error())
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
			return err
		}
		return
	}

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf("Day off set successfully, %s/%s/%s", year, month, day))
	if err != nil {
		logger.Errorf("Error sending message: %s", err)
		return
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)
	return
}
