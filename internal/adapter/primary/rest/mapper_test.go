package rest

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mapper Mapper = DefaultMapper{}

func TestToAndFromDomainRestaurant(t *testing.T) {
	restRestaurant := newRestaurant()
	domainRestaurant := mapper.toDomainRestaurant(restRestaurant)
	backToRest := mapper.fromDomainRestaurant(domainRestaurant)

	assert.True(t, reflect.DeepEqual(restRestaurant, backToRest))
}

func TestToAndFromDomainRestaurants(t *testing.T) {
	restRestaurants := []*Restaurant{newRestaurant(), newRestaurant()}
	domainRestaurants := mapper.toDomainRestaurants(restRestaurants)
	backToRest := mapper.fromDomainRestaurants(domainRestaurants)

	assert.True(t, reflect.DeepEqual(restRestaurants, backToRest))
}

// --------------------------------------------------------------------------------
// Utility functions to create restaurants and menus.
// --------------------------------------------------------------------------------

func newRestaurant() *Restaurant {
	return &Restaurant{Name: "aName", Address: newAddress(), Menu: newMenu()}
}

func newAddress() *Address {
	return &Address{Street: "aStreet", City: "aCity", State: "aState", Zip: "aZip"}
}

func newMenu() *Menu {
	item1 := MenuItem{Id: 1, Name: "Name1", Price: "12.99"}
	item2 := MenuItem{Id: 2, Name: "Name2", Price: "4.55"}
	return &Menu{Items: []MenuItem{item1, item2}}
}
