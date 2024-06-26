package domain

import (
	"time"

	"github.com/yimincai/gopunch/internal/enums"
	"github.com/yimincai/gopunch/pkg/snowflake"
	"gorm.io/gorm"
)

type User struct {
	ID            string         `gorm:"primaryKey;type:varchar(100)" json:"-"`
	Name          string         `gorm:"type:varchar(100)" json:"name"`
	Account       string         `gorm:"type:varchat(100);not null" json:"account"`
	Password      string         `gorm:"type:varchar(255);not null" json:"password"`
	DiscordUserID string         `gorm:"type:varchar(255);not null" json:"discord_user_id"`
	IsEnable      bool           `gorm:"index" json:"-"`
	Role          enums.RoleType `gorm:"type:int" json:"role"`
	CreatedAt     time.Time      `json:"-"`
	UpdatedAt     time.Time      `json:"-"`
	DeletedAt     gorm.DeletedAt `json:"-"`
}

// BeforeCreate will set snowflake id rather than numeric id.
func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	u.ID = snowflake.GetID()
	return nil
}
