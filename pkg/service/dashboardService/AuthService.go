package dashboardService

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelInterface"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service"
)

const (
	// Cookies
	AuthCookieToken    = "token"
	AuthCookieUsername = "username"

	// Definitions
	LoginActionTitle  = "Login"
	LoginActionLink   = "/login"
	LogoutActionTitle = "Logout"
	LogoutActionLink  = "/logout"
)

type AuthService struct {
	manager     *manager.Manager
	userService *service.UserService
}

// NewAuthService - constructor of AuthService struct
func NewAuthService(
	manager *manager.Manager,
	userService *service.UserService,
) *AuthService {
	return &AuthService{
		manager:     manager,
		userService: userService,
	}
}

// Login - check the user isset and check the token is matches, set it to cookie
func (auth *AuthService) Login(authData modelInterface.UserAuthInterface) (*modelDB.User, error) {
	user, err := auth.manager.Repository.User().FindByToken(authData.GetToken())
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.New("User not found by token")
	} else if user.Token == "" {
		return nil, errors.New("Token is not defined, please write `/token` to bot chat")
	}

	http.SetCookie(authData.GetWriter(), &http.Cookie{
		Name:  AuthCookieToken,
		Value: user.Token,
	})

	http.SetCookie(authData.GetWriter(), &http.Cookie{
		Name:  AuthCookieUsername,
		Value: user.Username,
	})

	return auth.GetUserFromCache(authData)
}

// IsGrantedAccessByCookies - check the access be cookies
func (auth *AuthService) IsGrantedAccessByCookies(w http.ResponseWriter, r *http.Request) error {
	authData, err := model.NewAuthCookieData(w, r)
	if err != nil {
		return err
	}

	if !auth.IsAuthorized(w, authData) {
		return errors.New("You are not logged in, please visit the login page.")
	}

	return nil
}

// IsAuthorized - check the user is authorized and isset in cache if not do logout action
func (auth *AuthService) IsAuthorized(w http.ResponseWriter, authData modelInterface.UserAuthInterface) bool {
	if !auth.isAuthorized(authData) {
		auth.Logout(authData)

		return false
	}

	return true
}

// isAuthorized - check the user is authorized and isset in cache
func (auth *AuthService) isAuthorized(authData modelInterface.UserAuthInterface) bool {
	// trying to find user into cache
	user, err := auth.GetUserFromCache(authData)
	if err == nil {
		return user.Token == authData.GetToken()
	}

	// trying to find user into database
	user, err = auth.manager.Repository.User().FindByToken(authData.GetToken())
	if err == nil {
		return strings.ToLower(user.Username) == strings.ToLower(authData.GetUsername())
	}

	return false
}

// GetUserFromCache - loading user from or into cache by authData
func (auth *AuthService) GetUserFromCache(authData modelInterface.UserAuthInterface) (*modelDB.User, error) {
	if _, issetInCache := auth.userService.Cache[authData.GetUsername()]; issetInCache {
		if auth.userService.Cache[authData.GetUsername()].Token != "" {
			return auth.userService.Cache[authData.GetUsername()], nil
		}
	} else {
		dbUser, err := auth.manager.Repository.User().FindByUsername(authData.GetUsername())
		if err != nil {
			return nil, err
		}

		if dbUser != nil {
			auth.userService.Cache[authData.GetUsername()] = dbUser

			return dbUser, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("User `%v` not found", authData.GetUsername()))
}

// Logout - set cookies to empty value
func (auth *AuthService) Logout(authData modelInterface.UserAuthInterface) {
	http.SetCookie(authData.GetWriter(), &http.Cookie{
		Name:    AuthCookieToken,
		Value:   "",
		Expires: time.Unix(0, 0),
	})

	http.SetCookie(authData.GetWriter(), &http.Cookie{
		Name:    AuthCookieUsername,
		Value:   "",
		Expires: time.Unix(0, 0),
	})
}
