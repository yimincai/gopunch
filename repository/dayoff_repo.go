package repository

import "github.com/yimincai/gopunch/domain"

// FindUserDayOffByDate implements Repository.
func (r *Repo) FindUserDayOffByDate(userID string, year int, month int, date int) (*domain.DayOff, error) {
	var dayOff *domain.DayOff
	if err := r.Db.Where("user_id = ? AND year = ? AND month = ? AND date = ?", userID, year, month, date).First(&dayOff).Error; err != nil {
		return nil, err
	}

	return dayOff, nil
}

// SetDayOff implements Repository.
func (r *Repo) SetDayOff(dayOff *domain.DayOff) error {
	if err := r.Db.Create(dayOff).Error; err != nil {
		return err
	}

	return nil
}
