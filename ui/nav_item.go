package ui

type NavItem interface {
	View
	Page() Page
	URL() PageURL
	NavTitle() string
	Children() []NavItem
}
