package repository

import (
	"github.com/yimincai/gopunch/domain"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(user *domain.User) error
	GetUsers() ([]*domain.User, error)
	GetUserByAccount(account string) (*domain.User, error)
	GetUserByDiscordUserID(discordUserID string) (*domain.User, error)
	DeleteUserByAccount(account string) error
	SetDayOff(dayOff *domain.DayOff) error
	FindUserDayOffByDate(userID string, year, month, date int) (*domain.DayOff, error)
	FindAllSchedules() ([]*domain.Schedule, error)
	FindScheduleByUserID(userID string) (*domain.Schedule, error)
	CreateSchedule(schedule *domain.Schedule) error
	UpdateSchedule(schedule *domain.Schedule) error
	UpdateAccount(user *domain.User) error
}

type Repo struct {
	Db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &Repo{Db: db}
}
