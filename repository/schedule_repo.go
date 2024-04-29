package repository

import (
	"github.com/yimincai/gopunch/domain"
)

// FindAllSchedules implements Repository.
func (r *Repo) FindAllSchedules() ([]*domain.Schedule, error) {
	var schedules []*domain.Schedule

	if err := r.Db.Preload("User").Find(&schedules).Error; err != nil {
		return nil, err
	}

	if len(schedules) == 0 {
		return nil, nil
	}

	return schedules, nil
}

// CreateSchedule implements Repository.
func (r *Repo) CreateSchedule(schedule *domain.Schedule) (*domain.Schedule, error) {
	if err := r.Db.Create(schedule).Error; err != nil {
		return nil, err
	}

	return schedule, nil
}

// FindScheduleByUserID implements Repository.
func (r *Repo) FindScheduleByUserID(userID string) (*domain.Schedule, error) {
	var schedule *domain.Schedule

	if err := r.Db.Preload("User").Where("user_id = ?", userID).First(&schedule).Error; err != nil {
		return nil, err
	}

	return schedule, nil
}

// UpdateSchedule implements Repository.
func (r *Repo) UpdateSchedule(schedule *domain.Schedule) (*domain.Schedule, error) {
	if err := r.Db.Updates(schedule).Error; err != nil {
		return nil, err
	}

	return schedule, nil
}
