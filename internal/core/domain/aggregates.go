package domain

import "errors"

// --------------------------------------------------------------------------------
// Aggregate Root :: Restaurant
// --------------------------------------------------------------------------------

// Restaurant is an entity representing a restaurant in the system that is placed
// in an address and offers to customers a menu. It's also the aggregate root to
// manage addresses and menus.
type Restaurant struct {
	Id      uint64
	Name    string
	Address *Address
	Menu    *Menu
}

// UpdateMenu updates the menu of a restaurant, enforcing some invariants on
// the menu (e.g. it cannot be empty).
func (r *Restaurant) UpdateMenu(menu *Menu) error {
	// Enforce invariants for the menu, even if you have previous validations in
	// external layers.
	if !isValidMenu(menu) {
		return errors.New("invalid menu. Check the requirements of a menu in isValidMenu(menu *Menu) func")
	}
	r.Menu = menu
	return nil
}

func isValidMenu(menu *Menu) bool {
	return menu != nil && menu.items != nil && len(menu.items) > 0 && len(menu.items) <= 1000
}
