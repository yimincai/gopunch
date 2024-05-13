package repository

import (
	"github.com/yimincai/gopunch/domain"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(user *domain.User) (*domain.User, error)
	FindUsers() ([]*domain.User, error)
	FindUserByAccount(account string) (*domain.User, error)
	FindUserByDiscordUserID(discordUserID string) (*domain.User, error)
	FindUserByID(userID string) (*domain.User, error)
	DeleteUserByAccount(account string) error
	SetDayOff(dayOff *domain.DayOff) error
	FindUserDayOffByDate(userID string, year, month, date int) (*domain.DayOff, error)
	FindAllSchedules() ([]*domain.Schedule, error)
	FindScheduleByUserID(userID string) (*domain.Schedule, error)
	CreateSchedule(schedule *domain.Schedule) (*domain.Schedule, error)
	UpdateSchedule(schedule *domain.Schedule) (*domain.Schedule, error)
	UpdateAccount(user *domain.User) error
}

type Repo struct {
	Db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &Repo{Db: db}
}
