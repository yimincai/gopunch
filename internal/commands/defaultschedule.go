package commands

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/yimincai/gopunch/internal/bot"
	"github.com/yimincai/gopunch/internal/errs"
	"github.com/yimincai/gopunch/internal/service"
	"github.com/yimincai/gopunch/pkg/logger"
)

type CommandDefaultSchedule struct {
	Svc service.Service
}

func (c *CommandDefaultSchedule) IsAdminRequired() bool {
	return false
}

func (c *CommandDefaultSchedule) Invokes() []string {
	return []string{"DefaultSchedule", "ds", "defaultschedule", "defaultSchedule"}
}

func (c *CommandDefaultSchedule) Description() string {
	return "Show default schedule"
}

func (c *CommandDefaultSchedule) Exec(ctx *bot.Context) (err error) {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Punch In", "Punch Out"})
	t.AppendRow([]interface{}{"08:00", "18:00"})
	t.SetStyle(table.StyleLight)

	response := "```"
	response += "系統會自動在需打卡的第一個時間前30分鐘內隨機打卡，並在需打卡的最後一個時間後30分鐘內隨機打卡。\n"
	response += "六日不打卡！若有特殊需求，請使用 `SetSchedule` 指令設定。\n"
	response += fmt.Sprintf("國定假日尚未實作，請自行使用 %sdayoff 指令調整。\n\n", c.Svc.Cfg.Prefix)
	response += "Default Schedule:\n"
	response += t.Render()
	response += "```"

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
	if err != nil {
		return errs.ErrSendingMessage
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)
	return
}
