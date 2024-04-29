package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/yimincai/gopunch/internal/enums"
	"github.com/yimincai/gopunch/pkg/logger"
	"github.com/yimincai/gopunch/pkg/snowflake"
	"gorm.io/gorm"
)

const DEFAULT_TIME_FORMAT = "15:04"
const DEFAULT_TOLERANCE = 30

var CronScheduledMap map[string]cron.EntryID

func init() {
	CronScheduledMap = make(map[string]cron.EntryID)
}

type Schedule struct {
	ID        string    `gorm:"primaryKey;type:varchar(100)" json:"-"`
	PunchIn   string    `gorm:"embedded" json:"punch_in"`
	PunchOut  string    `gorm:"embedded" json:"punch_out"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`

	// belongs to user
	UserID string `gorm:"type:varchar(100);not null" json:"user_id"`
	User   User   `json:"-"`
}

type Expression struct {
	PunchIn  string `json:"punch_in"`
	PunchOut string `json:"punch_out"`
}

// BeforeCreate will set snowflake id rather than numeric id.
func (s *Schedule) BeforeCreate(_ *gorm.DB) (err error) {
	s.ID = snowflake.GetID()

	return nil
}

func (s *Schedule) GetCronExpression() Expression {
	now := time.Now()
	year, month, day := now.Date()

	inStr := strings.Split(s.PunchIn, ":")
	inHour, err := strconv.Atoi(inStr[0])
	if err != nil {
		logger.Fatalf("Error converting string to int: %v", err)
	}
	inMinute, err := strconv.Atoi(inStr[1])
	if err != nil {
		logger.Fatalf("Error converting string to int: %v", err)
	}

	outStr := strings.Split(s.PunchOut, ":")
	outHour, err := strconv.Atoi(outStr[0])
	if err != nil {
		logger.Fatalf("Error converting string to int: %v", err)
	}
	outMinute, err := strconv.Atoi(outStr[1])
	if err != nil {
		logger.Fatalf("Error converting string to int: %v", err)
	}

	in := time.Date(year, month, day, inHour, inMinute, 0, 0, now.Location()).Add(-time.Duration(DEFAULT_TOLERANCE) * time.Minute)
	out := time.Date(year, month, day, outHour, outMinute, 0, 0, now.Location())

	return Expression{
		PunchIn:  fmt.Sprintf("%d %d * * *", in.Minute(), in.Hour()),
		PunchOut: fmt.Sprintf("%d %d * * *", out.Minute(), out.Hour()),
	}
}

func (s *Schedule) GetCronEntryKey(punchType enums.PunchType) string {
	return fmt.Sprintf("%s_%d", s.UserID, punchType)
}

func (s *Schedule) GetCronEntry(punchType enums.PunchType) cron.EntryID {
	key := fmt.Sprintf("%s_%d", s.UserID, punchType)
	if ok := CronScheduledMap[key]; ok != 0 {
		return CronScheduledMap[key]
	}

	return 0
}
