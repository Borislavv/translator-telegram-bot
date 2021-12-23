package modelDashboard

type Page struct {
	Menu *Menu
}

// NewPage - constructor of Page struct
func NewPage() *Page {
	return &Page{}
}

// AddMenu - setter of Menu
func (page *Page) AddMenu(menu *Menu) {
	page.Menu = menu
}
