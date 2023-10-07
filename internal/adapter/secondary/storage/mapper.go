package storage

import (
	"math/big"

	"f4allgo-restaurant/internal/core/domain"
)

type Mapper interface {

	// fromDomainRestaurant maps a domain.Restaurant struct into a Restaurant.
	fromDomainRestaurant(restaurant *domain.Restaurant) *Restaurant

	// fromDomainRestaurants maps a slice of domain.Restaurant into a slice of Restaurant.
	fromDomainRestaurants(restaurant []*domain.Restaurant) []*Restaurant

	// fromDomainAddress maps a domain.Address struct into an Address
	fromDomainAddress(address *domain.Address) *Address

	// fromDomainMenu maps a domain.Menu struct into a Menu.
	fromDomainMenu(menu *domain.Menu) []*MenuItem

	// toDomainRestaurant maps a Restaurant struct into a domain.Restaurant.
	toDomainRestaurant(restaurantDto *Restaurant) *domain.Restaurant

	// toDomainRestaurants maps a slice of Restaurant into a slice of domain.Restaurant.
	toDomainRestaurants(restaurants []*Restaurant) []*domain.Restaurant

	// toDomainAddress maps a Address struct into a domain.Address.
	toDomainAddress(address *Address) *domain.Address

	// toDomainMenu maps a Menu struct into a domain.Menu.
	toDomainMenu(menuItems []*MenuItem) *domain.Menu
}

// DefaultMapper is the default implementation of Mapper.
// It provides a plain manual mapping without using any third party library.
type DefaultMapper struct{}

// Interface compliance verification.
var _ Mapper = (*DefaultMapper)(nil)

func (dm DefaultMapper) fromDomainRestaurant(restaurant *domain.Restaurant) *Restaurant {
	if restaurant == nil {
		return nil
	}
	return &Restaurant{ID: restaurant.Id, Name: restaurant.Name, Address: dm.fromDomainAddress(restaurant.Address), Menu: dm.fromDomainMenu(restaurant.Menu)}
}

func (dm DefaultMapper) fromDomainRestaurants(restaurants []*domain.Restaurant) []*Restaurant {
	if restaurants == nil {
		return nil
	}
	dtoRestaurants := []*Restaurant{}
	for _, restaurant := range restaurants {
		dtoRestaurants = append(dtoRestaurants, dm.fromDomainRestaurant(restaurant))
	}

	return dtoRestaurants
}

func (DefaultMapper) fromDomainAddress(address *domain.Address) *Address {
	if address == nil {
		return nil
	}
	return &Address{Street: address.Street(), City: address.City(), State: address.State(), Zip: address.Zip()}
}

func (DefaultMapper) fromDomainMenu(menu *domain.Menu) []*MenuItem {
	if menu == nil {
		return nil
	}
	dtoItems := []*MenuItem{}
	for _, item := range menu.GetItems() {
		dtoItems = append(dtoItems, &MenuItem{Id: int32(item.GetId()), Name: item.GetName(), Price: item.GetPrice().Text('f', 2)})
	}

	return dtoItems
}

func (dm DefaultMapper) toDomainRestaurant(restaurantDto *Restaurant) *domain.Restaurant {
	if restaurantDto == nil {
		return nil
	}
	domainRestaurant := domain.Restaurant{}
	domainRestaurant.Id = restaurantDto.ID
	domainRestaurant.Name = restaurantDto.Name
	domainRestaurant.Address = dm.toDomainAddress(restaurantDto.Address)
	domainRestaurant.Menu = dm.toDomainMenu(restaurantDto.Menu)

	return &domainRestaurant
}

func (dm DefaultMapper) toDomainRestaurants(restaurants []*Restaurant) []*domain.Restaurant {
	if restaurants == nil {
		return nil
	}
	domainRestaurants := []*domain.Restaurant{}
	for _, restaurant := range restaurants {
		domainRestaurants = append(domainRestaurants, dm.toDomainRestaurant(restaurant))
	}

	return domainRestaurants
}

func (DefaultMapper) toDomainAddress(address *Address) *domain.Address {
	if address == nil {
		return nil
	}
	return domain.NewAddress(
		address.Street,
		address.City,
		address.State,
		address.Zip)
}

func (DefaultMapper) toDomainMenu(menuItems []*MenuItem) *domain.Menu {
	if menuItems == nil {
		return nil
	}
	domainItems := []*domain.MenuItem{}
	for _, item := range menuItems {
		f := new(big.Float)
		f.SetString(item.Price)
		domainItems = append(domainItems, domain.NewMenuItem(int16(item.Id), item.Name, f))
	}

	return domain.NewMenu(domainItems)
}
