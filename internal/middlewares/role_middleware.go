package middlewares

import (
	"github.com/yimincai/gopunch/internal/bot"
	"github.com/yimincai/gopunch/internal/enums"
	"github.com/yimincai/gopunch/internal/errs"
	"github.com/yimincai/gopunch/repository"
)

type RequiredAdminPermission struct {
	Repo repository.Repository
}

func (a *RequiredAdminPermission) Exec(ctx *bot.Context, cmd bot.Command) (next bool, err error) {
	if !cmd.IsAdminRequired() {
		return true, nil
	}

	user, err := a.Repo.GetUserByDiscordUserID(ctx.Message.Author.ID)
	if err != nil {
		return false, errs.ErrUserNotFound
	}

	if !user.IsEnable {
		return false, errs.ErrUserNotEnabled
	}

	if user.Role != enums.RoleType_Admin {
		return false, errs.ErrForbidden
	}

	return true, nil
}

func NewRequiredAdminPermission(repo repository.Repository) *RequiredAdminPermission {
	return &RequiredAdminPermission{
		Repo: repo,
	}
}
