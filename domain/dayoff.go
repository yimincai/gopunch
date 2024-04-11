package domain

import (
	"time"

	"github.com/yimincai/gopunch/pkg/snowflake"
	"gorm.io/gorm"
)

type DayOff struct {
	ID        string         `gorm:"primaryKey;type:varchar(100)" json:"-"`
	UserID    string         `gorm:"type:varchar(100);not null" json:"user_id"`
	Year      int            `json:"year"`
	Month     int            `json:"month"`
	Date      int            `json:"date"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

// BeforeCreate will set snowflake id rather than numeric id.
func (d *DayOff) BeforeCreate(_ *gorm.DB) (err error) {
	d.ID = snowflake.GetID()
	return nil
}
