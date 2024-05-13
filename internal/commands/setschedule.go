package commands

import (
	"errors"
	"fmt"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/yimincai/gopunch/domain"
	"github.com/yimincai/gopunch/internal/bot"
	"github.com/yimincai/gopunch/internal/errs"
	"github.com/yimincai/gopunch/internal/service"
	"github.com/yimincai/gopunch/pkg/logger"
	"gorm.io/gorm"
)

type CommandSetSchedule struct {
	Svc service.Service
}

func (c *CommandSetSchedule) IsAdminRequired() bool {
	return false
}

func (c *CommandSetSchedule) Invokes() []string {
	return []string{"SetSchedule", "ss", "setschedule", "setSchedule"}
}

func (c *CommandSetSchedule) Description() string {
	return "Set up the schedule for the user"
}

func (c *CommandSetSchedule) Exec(ctx *bot.Context) (err error) {
	if len(ctx.Args) != 2 {
		usage := fmt.Sprintf("Usage: %sSetSchedule <start time> <end time>", c.Svc.Cfg.Prefix)
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, usage)
		if err != nil {
			return errs.ErrSendingMessage
		}
		return
	}

	// validate the time format
	stime, err := time.Parse("15:04", ctx.Args[0])
	if err != nil {
		response := "❌ Invalid start time format, please use `HH:MM` in 24-hour format"
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
		if err != nil {
			return errs.ErrSendingMessage
		}
		return
	}

	etime, err := time.Parse("15:04", ctx.Args[1])
	if err != nil {
		response := "❌ Invalid start time format, please use `HH:MM` in 24-hour format"
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
		if err != nil {
			return errs.ErrSendingMessage
		}
		return
	}

	// check if the user is already registered
	user, err := c.Svc.Repo.FindUserByDiscordUserID(ctx.Message.Author.ID)
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
			// create a new schedule
			schedule, err = c.Svc.Repo.CreateSchedule(&domain.Schedule{
				PunchIn:  stime.Format("15:04"),
				PunchOut: etime.Format("15:04"),
				UserID:   user.ID,
			})

			if err != nil {
				return errs.ErrInternalError
			}

			err = c.Svc.AddSchedulePunch(schedule)
			if err != nil {
				return errs.ErrInternalError
			}

			response := "```"
			response += "✅ Schedule set up successfully\n"
			t := table.NewWriter()
			t.AppendHeader(table.Row{"User", "Punch In", "Punch Out", "Tolerance"})
			t.AppendRow([]interface{}{user.Account, schedule.PunchIn, schedule.PunchOut, "30 minutes"})
			t.SetStyle(table.StyleLight)
			response += t.Render()
			response += "```"
			_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
			if err != nil {
				return errs.ErrSendingMessage
			}
			return

		} else {
			return errs.ErrInternalError
		}
	}

	// remove the schedule from the cron job
	c.Svc.RemoveSchedulePunch(schedule)

	// update the schedule
	schedule.PunchIn = stime.Format("15:04")
	schedule.PunchOut = etime.Format("15:04")

	schedule, err = c.Svc.Repo.UpdateSchedule(&domain.Schedule{
		ID:       schedule.ID,
		PunchIn:  stime.Format("15:04"),
		PunchOut: etime.Format("15:04"),
		UserID:   schedule.UserID,
		User:     domain.User{},
	})
	if err != nil {
		return errs.ErrInternalError
	}

	err = c.Svc.AddSchedulePunch(schedule)
	if err != nil {
		return errs.ErrInternalError
	}

	response := "✅ Schedule updated successfully"
	t := table.NewWriter()
	t.AppendHeader(table.Row{"User", "Punch In", "Punch Out", "Tolerance"})
	t.AppendRow([]interface{}{user.Account, schedule.PunchIn, schedule.PunchOut, "30 minutes"})
	t.SetStyle(table.StyleLight)
	response += "```" + t.Render() + "```"
	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
	if err != nil {
		return errs.ErrSendingMessage
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)
	return
}
