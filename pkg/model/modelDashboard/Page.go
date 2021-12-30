package modelDashboard

import "github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"

type Page struct {
	Menu    *Menu
	Content interface{}
	Funcs   interface{}
	User    *modelDB.User
}

// NewPage - constructor of Page struct
func NewPage() *Page {
	return &Page{}
}

// AddMenu - setter of the Menu
func (page *Page) AddMenu(menu *Menu) {
	page.Menu = menu
}

// AddConent - setter of the Conent
func (page *Page) AddConent(content interface{}) {
	page.Content = content
}

// AddTemplateFuncs - setter of the Funcs (functions which will exec. into templates)
func (page *Page) AddTemplateFuncs(funcs interface{}) {
	page.Funcs = funcs
}

// AddUser - user data which will used into tempaltes
func (page *Page) AddUser(user *modelDB.User) {
	page.User = user
}
