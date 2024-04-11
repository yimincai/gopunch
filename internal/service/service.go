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
	"github.com/yimincai/gopunch/domain"
	"github.com/yimincai/gopunch/internal/config"
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
}

func init() {
	R = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (s *Service) Login(discordUserID string) (string, error) {
	user := &domain.User{}
	user, err := s.Repo.GetUserByDiscordUserID(discordUserID)
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

func (s *Service) Punch(accessToken string) error {
	now := time.Now()
	weekday := now.Weekday().String()

	if weekday == "Saturday" || weekday == "Sunday" {
		logger.Infof("Today is %v, no need to punch", weekday)
		return nil
	}

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

	return nil
}

func (s *Service) DefaultSchedulePunchAllUsers() error {
	var users []*domain.User
	users, err := s.Repo.GetUsers()
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
				logger.Infof("User %s is disabled, no need to punch", u.Account)
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
				logger.Infof("User %s has day off on %d/%d/%d, no need to punch", u.Account, dayoff.Year, dayoff.Month, dayoff.Date)
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
			err = s.Punch(accessToken)
			if err != nil {
				logger.Error(err)
				return
			}

			channel, err := s.Session.UserChannelCreate(u.DiscordUserID)
			if err != nil {
				logger.Errorf("Error creating user DM channel: %s", err)
				return
			}

			// notify user that punch is done with bot
			_, err = s.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("✅ %s punched successfully at %s", user.Account, utils.TimeFormat(time.Now())))
			if err != nil {
				logger.Errorf("Error sending message: %s", err)
			}
		}(user)
	}

	wg.Wait()
	return nil
}

func (s *Service) Register(user *domain.User) error {
	return s.Repo.CreateUser(user)
}

func (s *Service) SetDayOff(discordUserID string, year, month, day int) error {
	user, err := s.Repo.GetUserByDiscordUserID(discordUserID)
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

func NewService(cfg *config.Config, repo repository.Repository, session *discordgo.Session) Service {
	return Service{
		Cfg:     cfg,
		Repo:    repo,
		Session: session,
	}
}
