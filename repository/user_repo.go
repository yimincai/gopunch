package repository

import (
	"github.com/yimincai/gopunch/domain"
	"github.com/yimincai/gopunch/internal/errs"
)

// CreateUser implements Repository.
func (r *Repo) CreateUser(user *domain.User) (*domain.User, error) {
	if err := r.Db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// FindUserByAccount implements Repository.
func (r *Repo) FindUserByAccount(account string) (*domain.User, error) {
	var user *domain.User
	if err := r.Db.Where("account = ?", account).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// FindUserByDiscordUserID implements Repository.
func (r *Repo) FindUserByDiscordUserID(discordUserID string) (*domain.User, error) {
	var user *domain.User
	if err := r.Db.Where("discord_user_id = ?", discordUserID).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// FindUsers implements Repository.
func (r *Repo) FindUsers() ([]*domain.User, error) {
	var users []*domain.User

	if err := r.Db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// FindUserByID implements Repository.
func (r *Repo) FindUserByID(userID string) (*domain.User, error) {
	var user *domain.User
	if err := r.Db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUserByAccount implements Repository.
func (r *Repo) DeleteUserByAccount(account string) error {
	if err := r.Db.Where("account = ?", account).Delete(&domain.User{}).Error; err != nil {
		return errs.ErrUserNotFound
	}

	return nil
}

// UpdateAccount implements Repository.
func (r *Repo) UpdateAccount(user *domain.User) error {
	if err := r.Db.Updates(user).Error; err != nil {
		return err
	}

	return nil
}
