package modelDashboard

type Menu struct {
	Items []*MenuItem
}

type MenuItem struct {
	Title string
	Link  string
}

// NewMenu - constructor of Menu struct
func NewMenu() *Menu {
	return &Menu{
		Items: []*MenuItem{
			{
				Title: "Notifications",
				Link:  "/notifications",
			},
			{
				Title: "Translation",
				Link:  "/translation",
			},
			{
				Title: "Docs",
				Link:  "/docs",
			},
			{
				Title: "About",
				Link:  "/about",
			},
		},
	}
}
