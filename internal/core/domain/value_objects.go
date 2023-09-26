package domain

import "math/big"

// --------------------------------------------------------------------------------
// VO :: Address
// --------------------------------------------------------------------------------

// Address is a value object to represent addresses (e.g. a restaurant Address).
type Address struct {
	street string
	city   string
	state  string
	zip    string
}

func NewAddress(street string, city string, state string, zip string) *Address {
	return &Address{street: street, city: city, state: state, zip: zip}
}

func (a *Address) Street() string {
	return a.street
}

func (a *Address) City() string {
	return a.city
}

func (a *Address) State() string {
	return a.state
}

func (a *Address) Zip() string {
	return a.zip
}

// --------------------------------------------------------------------------------
// VO :: Menu
// --------------------------------------------------------------------------------

// Menu is a value object to represent restaurant menus. Is composed by menu items.
type Menu struct {
	items []*MenuItem
}

func NewMenu(items []*MenuItem) *Menu {
	return &Menu{items: items}
}

func (m *Menu) GetItems() []*MenuItem {
	return m.items
}

// --------------------------------------------------------------------------------
// VO :: MenuItem
// --------------------------------------------------------------------------------

// MenuItem is a value object to represent a menu item.
type MenuItem struct {
	id    int16
	name  string
	price *big.Float
}

func NewMenuItem(id int16, name string, price *big.Float) *MenuItem {
	return &MenuItem{id: id, name: name, price: price}
}

func (mi *MenuItem) GetId() int16 {
	return mi.id
}

func (mi *MenuItem) GetName() string {
	return mi.name
}

func (mi *MenuItem) GetPrice() *big.Float {
	return mi.price
}
