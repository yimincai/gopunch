package repository

import (
	"strconv"
	"time"

	"github.com/yimincai/gopunch/domain"
	"github.com/yimincai/gopunch/internal/errs"
)

// CreateUser implements Repository.
func (r *Repo) CreateUser(user *domain.User) error {
	if err := r.Db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

// GetUserByAccount implements Repository.
func (r *Repo) GetUserByAccount(account string) (*domain.User, error) {
	var user *domain.User
	if err := r.Db.Where("account = ?", account).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByDiscordUserID implements Repository.
func (r *Repo) GetUserByDiscordUserID(discordUserID string) (*domain.User, error) {
	var user *domain.User
	if err := r.Db.Where("discord_user_id = ?", discordUserID).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// GetUsers implements Repository.
func (r *Repo) GetUsers() ([]*domain.User, error) {
	var users []*domain.User

	if err := r.Db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// DeleteUserByAccount implements Repository.
func (r *Repo) DeleteUserByAccount(account string) error {
	var user domain.User
	if err := r.Db.Where("account = ?", account).Delete(&user).Error; err != nil {
		return errs.ErrUserNotFound
	}

	// update user's account to let it be unique
	now := time.Now().Unix()
	dummy := account + "_deleted_" + strconv.FormatInt(now, 10)
	if err := r.Db.Unscoped().Where("account = ?", account).Update("account", dummy).Error; err != nil {
		return err
	}

	return nil
}
