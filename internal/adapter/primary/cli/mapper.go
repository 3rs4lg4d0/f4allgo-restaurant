package cli

import (
	"math/big"

	"f4allgo-restaurant/internal/core/domain"
)

// Mapper maps structs and slices from the rest handler layer to the
// core layer and vice versa.
type Mapper interface {

	// toDomainRestaurant maps a Restaurant struct into a domain.Restaurant.
	toDomainRestaurant(*Restaurant) *domain.Restaurant

	// toDomainRestaurants maps a slice of Restaurant into a slice of domain.Restaurant.
	toDomainRestaurants([]*Restaurant) []*domain.Restaurant

	// toDomainAddress maps a Address struct into a domain.Address.
	toDomainAddress(*Address) *domain.Address

	// toDomainMenu maps a Menu struct into a domain.Menu.
	toDomainMenu(*Menu) *domain.Menu

	// fromDomainRestaurant maps a domain.Restaurant struct into a Restaurant.
	fromDomainRestaurant(*domain.Restaurant) *Restaurant

	// fromDomainRestaurants maps a slice of domain.Restaurant into a slice of Restaurant.
	fromDomainRestaurants([]*domain.Restaurant) []*Restaurant

	// fromDomainAddress maps a domain.Address struct into an Address.
	fromDomainAddress(*domain.Address) *Address

	// fromDomainMenu maps a domain.Menu struct into a Menu.
	fromDomainMenu(*domain.Menu) *Menu
}

// DefaultMapper is the default implementation of Mapper.
// It provides a plain manual mapping without using any third party library.
type DefaultMapper struct{}

// Interface compliance verification.
var _ Mapper = (*DefaultMapper)(nil)

// toDomainRestaurant maps a Restaurant struct into a domain.Restaurant.
func (dm DefaultMapper) toDomainRestaurant(r *Restaurant) *domain.Restaurant {
	domainRestaurant := domain.Restaurant{}
	domainRestaurant.Name = r.Name
	domainRestaurant.Address = dm.toDomainAddress(r.Address)
	domainRestaurant.Menu = dm.toDomainMenu(r.Menu)
	return &domainRestaurant
}

// ToDomainRestaurants maps a slice of Restaurant into a slice of domain.Restaurant.
func (dm DefaultMapper) toDomainRestaurants(restaurants []*Restaurant) []*domain.Restaurant {
	items := []*domain.Restaurant{}
	for _, item := range restaurants {
		items = append(items, dm.toDomainRestaurant(item))
	}

	return items
}

// ToDomainAddress maps a Address struct into a domain.Address.
func (DefaultMapper) toDomainAddress(a *Address) *domain.Address {
	return domain.NewAddress(a.Street, a.City, a.State, a.Zip)
}

// ToDomainMenu maps a Menu struct into a domain.Menu.
func (DefaultMapper) toDomainMenu(m *Menu) *domain.Menu {
	domainItems := []*domain.MenuItem{}
	for _, item := range m.Items {
		f := new(big.Float)
		f.SetString(item.Price)
		domainItems = append(domainItems, domain.NewMenuItem(int16(item.Id), item.Name, f))
	}
	return domain.NewMenu(domainItems)
}

// FromDomainRestaurant maps a domain.Restaurant struct into a Restaurant.
func (dm DefaultMapper) fromDomainRestaurant(r *domain.Restaurant) *Restaurant {
	restRestaurant := Restaurant{}
	restRestaurant.Id = r.Id
	restRestaurant.Name = r.Name
	restRestaurant.Address = dm.fromDomainAddress(r.Address)
	restRestaurant.Menu = dm.fromDomainMenu(r.Menu)
	return &restRestaurant
}

// FromDomainRestaurants maps a slice of domain.Restaurant into a slice of Restaurant.
func (dm DefaultMapper) fromDomainRestaurants(restaurants []*domain.Restaurant) []*Restaurant {
	items := []*Restaurant{}
	for _, item := range restaurants {
		items = append(items, dm.fromDomainRestaurant(item))
	}

	return items
}

// FromDomainAddress maps a domain.Address struct into a Address.
func (DefaultMapper) fromDomainAddress(a *domain.Address) *Address {
	return &Address{Street: a.Street(), City: a.City(), State: a.State(), Zip: a.Zip()}
}

// FromDomainMenu maps a domain.Menu struct into a Menu.
func (DefaultMapper) fromDomainMenu(menu *domain.Menu) *Menu {
	items := []MenuItem{}
	for _, item := range menu.GetItems() {
		items = append(items, MenuItem{Id: int32(item.GetId()), Name: item.GetName(), Price: item.GetPrice().String()})
	}
	return &Menu{Items: items}
}
