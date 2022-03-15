package service

import (
	"log"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/util"
)

type UserService struct {
	manager *manager.Manager
	Cache   map[string]*modelDB.User
}

// NewUserService - constructor of UserService
func NewUserService(manager *manager.Manager) *UserService {
	return &UserService{
		manager: manager,
		Cache:   map[string]*modelDB.User{},
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
