package service

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/util"
)

var dateTimeLayout = "2006-01-02 15:04:05"
var dateTimeShortLayout = "2006-01-02 15"
var mx *sync.Mutex

type UserService struct {
	// deps.
	manager   *manager.Manager
	tzFetcher *TimeZoneFetcherService
	// vals.
	Cache map[string]*modelDB.User
}

// NewUserService - constructor of UserService
func NewUserService(manager *manager.Manager, tzFetcher *TimeZoneFetcherService) *UserService {
	return &UserService{
		manager:   manager,
		tzFetcher: tzFetcher,
		Cache:     map[string]*modelDB.User{},
	}
}

// GetUser - getting a user by the username field. If the user is not found, it will be created.
func (userService *UserService) GetUser(username string, chatId int64) (*modelDB.User, error) {
	// trying to find user into the cache
	if _, issetInCache := userService.Cache[username]; !issetInCache {
		// trying to find user into database
		dbUser, err := userService.manager.Repository.User().FindByUsername(username)
		if err != nil {
			log.Println(util.Trace(err))
			return nil, err
		} else {
			// user was not fonud, then create and store it
			if dbUser == nil {
				newUser := modelDB.NewUser()
				newUser.ChatId = chatId
				newUser.Username = username

				// store user
				dbUser, err = userService.manager.Repository.User().Create(newUser)
				if err != nil {
					log.Println(util.Trace(err))
					return nil, err
				}
			}
		}

		// store user to cache
		userService.Cache[username] = dbUser
	}

	return userService.Cache[username], nil
}

// SetToken - setting token to user
func (userService *UserService) SetToken(user *modelDB.User, token string) (*modelDB.User, error) {
	user, err := userService.manager.Repository.User().SetTokenById(user, token)
	if err != nil {
		log.Println(util.Trace(err))
		return nil, err
	}

	// update user instance in cache
	userService.Cache[user.Username] = user

	return user, nil
}

// GetUserTimeZone - goes over all timezones and trying determine the target one.
func (userService *UserService) GetUserTimeZone(dateTime string) (string, error) {
	var loc *time.Location

	userDateTime, err := time.Parse(dateTimeLayout, dateTime)
	if err != nil {
		return "", err
	}

	now := time.Now()
	for _, tz := range userService.tzFetcher.GetTimeZones() {
		loc, _ = time.LoadLocation(tz)

		if now.In(loc).Format(dateTimeShortLayout) == userDateTime.Format(dateTimeShortLayout) {
			return tz, nil
		}
	}

	return "", errors.New("Unable determine timezone of user provided string.")
}
