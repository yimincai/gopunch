package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	"github.com/yimincai/gopunch/domain"
	"github.com/yimincai/gopunch/internal/config"
	"github.com/yimincai/gopunch/internal/enums"
	"github.com/yimincai/gopunch/internal/errs"
	"github.com/yimincai/gopunch/pkg/logger"
	"github.com/yimincai/gopunch/pkg/utils"
	"github.com/yimincai/gopunch/repository"
	"gorm.io/gorm"
)

var R *rand.Rand

type Service struct {
	Session *discordgo.Session
	Cfg     *config.Config
	Repo    repository.Repository
	Cron    *cron.Cron
}

func init() {
	R = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (s *Service) InitSchedules() error {
	err := s.initUsersSchedules()
	if err != nil {
		return err
	}

	err = s.initDefaultSchedules()
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) initUsersSchedules() error {
	// schedule all users punch from schedule table
	schedules, err := s.Repo.FindAllSchedules()
	if err != nil {
		return err
	}

	if schedules == nil {
		logger.Info("No users schedules found")
		return nil
	}

	logger.Infof("Found %d users schedules", len(schedules))

	for _, schedule := range schedules {
		// check if user is enabled
		if !schedule.User.IsEnable {
			logger.Infof("User %s is disabled, pass scheduling", schedule.User.Account)
			continue
		}

		err = s.AddSchedulePunch(schedule)
		if err != nil {
			logger.Error(err)
		}

		logger.Infof("User %s schedule initialized", schedule.User.Account)
	}

	logger.Info("Users schedules initialized")

	return nil
}

func (s *Service) initDefaultSchedules() error {
	// default schedule punch all users at 07:30 every workday
	_, err := s.Cron.AddFunc("30 7 * * *", func() {
		now := time.Now()
		weekday := now.Weekday().String()

		if weekday == "Saturday" || weekday == "Sunday" {
			logger.Infof("Today is %v, don't need to punch", weekday)
			return
		} else {
			err := s.DefaultSchedulePunchAllUsers()
			if err != nil {
				logger.Error(err)
			}
			logger.Info("All Users Punch Done")
		}
	})
	if err != nil {
		logger.Errorf("Error adding default schedule punch all users at 07:30 every workday: %s", err)
	}

	// default schedule punch for all users at 18:00 every workday
	_, err = s.Cron.AddFunc("0 18 * * *", func() {
		now := time.Now()
		weekday := now.Weekday().String()

		if weekday == "Saturday" || weekday == "Sunday" {
			logger.Infof("Today is %v, don't need to punch", weekday)
			return
		} else {
			err := s.DefaultSchedulePunchAllUsers()
			if err != nil {
				logger.Error(err)
			}
			logger.Info("All Users Punch Done")
		}
	})
	if err != nil {
		logger.Errorf("Error adding default schedule punch for all users at 18:00 every workday: %s", err)
	}

	logger.Info("Default schedules initialized")

	return nil
}

func (s *Service) AddSchedulePunch(schedule *domain.Schedule) error {
	// check if user has day off
	now := time.Now()
	dayoff, err := s.Repo.FindUserDayOffByDate(schedule.UserID, now.Year(), utils.MonthToInt(now.Month()), now.Day())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// no day off record found, continue
	} else if err != nil {
		return err
	}

	if dayoff != nil {
		logger.Infof("User %s has day off on %d/%d/%d, don't need to punch", schedule.User.Account, dayoff.Year, dayoff.Month, dayoff.Date)
		return nil
	}
	// schedule punch in
	expression := schedule.GetCronExpression()
	pInEntryID, err := s.Cron.AddFunc(expression.PunchIn, func() {
		randomDelay := s.RamdomDelayInThirtyMinutes()
		logger.Debugf("User: %s punch out, random delay: %s", schedule.User.Account, randomDelay)
		time.Sleep(randomDelay)
		s.Punch(schedule.UserID)
	})
	if err != nil {
		logger.Errorf("Error adding schedule for user %s: %s", schedule.User.Account, err)
		return err
	}
	domain.CronScheduledMap[schedule.GetCronEntryKey(enums.PunchType_In)] = pInEntryID

	// schedule punch out
	pOutEntryID, err := s.Cron.AddFunc(expression.PunchOut, func() {
		randomDelay := s.RamdomDelayInThirtyMinutes()
		logger.Debugf("User: %s punch out, random delay: %s", schedule.User.Account, randomDelay)
		time.Sleep(randomDelay)
		s.Punch(schedule.UserID)
	})
	if err != nil {
		logger.Errorf("Error adding schedule for user %s: %s", schedule.User.Account, err)
		return err
	}
	domain.CronScheduledMap[schedule.GetCronEntryKey(enums.PunchType_Out)] = pOutEntryID

	return nil
}

// remove old schedule
func (s *Service) RemoveSchedulePunch(schedule *domain.Schedule) {
	pInEntryID := schedule.GetCronEntry(enums.PunchType_In)
	pOutEntryID := schedule.GetCronEntry(enums.PunchType_Out)

	s.Cron.Remove(pInEntryID)
	s.Cron.Remove(pOutEntryID)
}

func (s *Service) Punch(userID string) {
	// check if user has day off
	now := time.Now()
	weekday := now.Weekday().String()

	if weekday == "Saturday" || weekday == "Sunday" {
		logger.Infof("Today is %v, don't need to punch", weekday)
		return
	}

	user, err := s.Repo.FindUserByID(userID)
	if err != nil {
		logger.Errorf("Error finding user %s: %s", userID, err)
		return
	}

	dayoff, err := s.Repo.FindUserDayOffByDate(userID, now.Year(), utils.MonthToInt(now.Month()), now.Day())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// no day off record found, continue
	} else if err != nil {
		logger.Errorf("Error finding day off: %s", err)
		return
	}

	if dayoff != nil {
		logger.Infof("User %s has day off on %d/%d/%d, don't need to punch", user.Account, dayoff.Year, dayoff.Month, dayoff.Date)
		return
	}

	channel, err := s.Session.UserChannelCreate(user.DiscordUserID)
	if err != nil {
		logger.Errorf("Error creating user DM channel: %s", err)
		return
	}

	accessToken, err := s.Login(user.DiscordUserID)
	if err != nil {
		// notify user that login is failed with bot
		_, err = s.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("❌ %s scheduled login failed at %s", user.Account, utils.TimeFormat(time.Now())))
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
		}
		logger.Errorf("Error while login user %s, schedule skipped: %s", user.Account, err)
		return
	}

	err = s.WebPunch(accessToken)
	if err != nil {
		// notify user that punch is failed with bot
		_, err = s.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("❌ %s scheduled punch failed at %s", user.Account, utils.TimeFormat(time.Now())))
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
		}
		logger.Errorf("Error punching user %s: %s", user.Account, err)
		return
	}

	// notify user that punch is done with bot
	_, err = s.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("✅ %s scheduled punched successfully at %s", user.Account, utils.TimeFormat(time.Now())))
	if err != nil {
		logger.Errorf("Error sending message: %s", err)
	}
}

func (s *Service) Login(discordUserID string) (string, error) {
	user, err := s.Repo.FindUserByDiscordUserID(discordUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errs.ErrUserNotFound
		}
		return "", err
	}

	payload := &domain.LoginRequest{
		Account:  user.Account,
		Password: user.Password,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", s.Cfg.Endpoint+s.Cfg.LoginApiPath, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errs.ErrLoginFailed
	}

	var response domain.LoginResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	return response.Result.AccessToken, nil
}

func (s *Service) TryToLogin(account, password string) (string, error) {
	payload := &domain.LoginRequest{
		Account:  account,
		Password: password,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", s.Cfg.Endpoint+s.Cfg.LoginApiPath, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errs.ErrLoginFailed
	}

	var response domain.LoginResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	return response.Result.AccessToken, nil
}

func (s *Service) WebPunch(accessToken string) error {
	req, err := http.NewRequest("POST", s.Cfg.Endpoint+s.Cfg.PunchApiPath, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errs.ErrPunchFailed
	}

	var response domain.PunchResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		logger.Error(err)
		return err
	}

	if !response.Result.IsSuccess {
		return errs.ErrPunchFailed
	}

	return nil
}

func (s *Service) DefaultSchedulePunchAllUsers() error {
	var users []*domain.User
	users, err := s.Repo.FindUsers()
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	for _, user := range users {
		user := user
		wg.Add(1)
		go func(u *domain.User) {
			defer wg.Done()

			// check if user is enabled
			if !u.IsEnable {
				logger.Infof("User %s is disabled, pass default scheduling", u.Account)
				return
			}

			// check if user has scheduled punch
			schedule, err := s.Repo.FindScheduleByUserID(u.ID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// no schedule found, continue
				} else {
					logger.Errorf("Error finding schedule: %s", err)
					return
				}
			}
			if schedule != nil {
				logger.Infof("User %s has scheduled punch, pass default scheduling", u.Account)
				return
			}

			// check if user has day off
			now := time.Now()
			dayoff, err := s.Repo.FindUserDayOffByDate(u.ID, now.Year(), utils.MonthToInt(now.Month()), now.Day())
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// no day off record found, continue
			} else if err != nil {
				logger.Error(err)
				return
			}

			if dayoff != nil {
				logger.Infof("User %s has day off on %d/%d/%d, don't need to punch", u.Account, dayoff.Year, dayoff.Month, dayoff.Date)
				return
			}

			ramdomDelay := time.Duration(R.Intn(29))*time.Minute + time.Duration(R.Intn(60))*time.Second
			logger.Debugf("%s punch will be delayed for %s, executed at %s", user.Account, ramdomDelay, utils.TimeFormat(time.Now().Add(ramdomDelay)))
			time.Sleep(ramdomDelay)
			accessToken, err := s.Login(u.DiscordUserID)
			if err != nil {
				logger.Error(err)
				return
			}

			channel, err := s.Session.UserChannelCreate(u.DiscordUserID)
			if err != nil {
				logger.Errorf("Error creating user DM channel: %s", err)
				return
			}

			err = s.WebPunch(accessToken)
			if err != nil {
				// notify user that punch is failed with bot
				_, err = s.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("❌ %s schedule punch failed at %s", user.Account, utils.TimeFormat(time.Now())))
				if err != nil {
					logger.Errorf("Error sending message: %s", err)
				}
				logger.Error(err)
				return
			}

			// notify user that punch is done with bot
			_, err = s.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("✅ %s schedule punched successfully at %s", user.Account, utils.TimeFormat(time.Now())))
			if err != nil {
				logger.Errorf("Error sending message: %s", err)
			}
		}(user)
	}

	wg.Wait()
	return nil
}

func (s *Service) Register(user *domain.User) error {
	_, err := s.Repo.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) SetDayOff(discordUserID string, year, month, day int) error {
	user, err := s.Repo.FindUserByDiscordUserID(discordUserID)
	if err != nil {
		return errs.ErrUserNotFound
	}

	if !user.IsEnable {
		return errs.ErrUserNotEnabled
	}

	dayoff := &domain.DayOff{
		UserID: user.ID,
		Year:   year,
		Month:  month,
		Date:   day,
	}

	return s.Repo.SetDayOff(dayoff)
}

func (s *Service) PunchOnTime(accessToken string, punchTime time.Time) error {
	type payload struct {
		Timestamp string `json:"timestamp"`
	}

	var p payload
	p.Timestamp = punchTime.Format("2006/01/02 15:04:05")

	body, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.Cfg.Endpoint+s.Cfg.PunchApiPath, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errs.ErrPunchFailed
	}

	var response domain.PunchResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		logger.Error(err)
		return err
	}

	if !response.Result.IsSuccess {
		return errs.ErrPunchFailed
	}

	return nil
}

func (s *Service) RamdomDelayInThirtyMinutes() time.Duration {
	return time.Duration(R.Intn(29))*time.Minute + time.Duration(R.Intn(60))*time.Second
}

func NewService(cfg *config.Config, repo repository.Repository, session *discordgo.Session, cron *cron.Cron) Service {
	return Service{
		Cfg:     cfg,
		Repo:    repo,
		Session: session,
		Cron:    cron,
	}
}
