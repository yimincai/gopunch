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
			return err
		}
		return
	}

	sp := strings.Split(ctx.Args[0], "/")
	if len(sp) != 3 {
		usage := fmt.Sprintf("Usage: %sDayOff <date>, format example: %sDayOff 2024/12/31", c.Svc.Cfg.Prefix, c.Svc.Cfg.Prefix)
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, usage)
		if err != nil {
			return errs.ErrSendingMessage
		}
		return
	}

	user, err := c.Svc.Repo.GetUserByDiscordUserID(ctx.Message.Author.ID)
	if err != nil {
		return errs.ErrUserNotFound
	}

	if !user.IsEnable {
		return errs.ErrUserNotEnabled
	}

	year := sp[0]
	month := sp[1]
	day := sp[2]

	//case to int
	y, err := strconv.Atoi(year)
	if err != nil {
		return errs.ErrInvalidDate
	}

	m, err := strconv.Atoi(month)
	if err != nil {
		return errs.ErrInvalidDate
	}

	d, err := strconv.Atoi(day)
	if err != nil {
		return errs.ErrInvalidDate
	}

	// validate date
	if !utils.IsVaildateDate(y, m, d) {
		return errs.ErrInvalidDate
	}

	existsDayOff, err := c.Svc.Repo.FindUserDayOffByDate(user.ID, y, m, d)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errs.ErrInternalError
	}

	if existsDayOff != nil {
		return errs.ErrDayOffAlreadySet
	}

	err = c.Svc.SetDayOff(ctx.Message.Author.ID, y, m, d)
	if err != nil {
		return errs.ErrInternalError
	}

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf("Day off set successfully, %s/%s/%s", year, month, day))
	if err != nil {
		return errs.ErrSendingMessage
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)
	return
}
