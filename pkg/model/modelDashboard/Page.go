package modelDashboard

type Page struct {
	Menu    *Menu
	Content interface{}
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
