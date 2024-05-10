package commands

import (
	"errors"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/yimincai/gopunch/internal/bot"
	"github.com/yimincai/gopunch/internal/errs"
	"github.com/yimincai/gopunch/internal/service"
	"github.com/yimincai/gopunch/pkg/logger"
	"gorm.io/gorm"
)

type CommandPrintSchedule struct {
	Svc service.Service
}

func (c *CommandPrintSchedule) IsAdminRequired() bool {
	return false
}

func (c *CommandPrintSchedule) Invokes() []string {
	return []string{"PrintSchedule", "ps", "printschedule", "printSchedule"}
}

func (c *CommandPrintSchedule) Description() string {
	return "Print the schedule"
}

func (c *CommandPrintSchedule) Exec(ctx *bot.Context) (err error) {
	// check if the user is already registered
	user, err := c.Svc.Repo.GetUserByDiscordUserID(ctx.Message.Author.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := fmt.Sprintf("You are not registered yet, please register first using `%sRegister <account> <password>`", c.Svc.Cfg.Prefix)
			_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
			if err != nil {
				return errs.ErrSendingMessage
			}
			return
		} else {
			return errs.ErrInternalError
		}
	}

	// check if the user has already set up the schedule
	schedule, err := c.Svc.Repo.FindScheduleByUserID(user.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg := fmt.Sprintf("You have not set up the schedule yet, now is using default schedule,\n if you want to set up the schedule, please using `%sSetSchedule <start time> <end time>`\nDefault Schedule:\n", c.Svc.Cfg.Prefix)
			t := table.NewWriter()
			t.AppendHeader(table.Row{"User", "Punch In", "Punch Out", "Tolerance"})
			t.AppendRow([]interface{}{user.Account, "08:00", "18:00", "30 minutes"})
			t.SetStyle(table.StyleLight)
			response := "```" + msg + t.Render() + "```"
			_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
			if err != nil {
				return errs.ErrSendingMessage
			}
			return

		} else {
			return errs.ErrInternalError
		}
	}

	// a table of schedule
	t := table.NewWriter()
	t.AppendHeader(table.Row{"User", "Punch In", "Punch Out", "Tolerance"})
	t.AppendRow([]interface{}{user.Account, schedule.PunchIn, schedule.PunchOut, "30 minutes"})
	t.SetStyle(table.StyleLight)
	response := "```" + "User schedule:\n" + t.Render() + "```"

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
	if err != nil {
		return errs.ErrSendingMessage
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)
	return
}
