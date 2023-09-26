package rest

import (
	"reflect"
	"testing"
)

var mapper Mapper = DefaultMapper{}

const notEqualsErrorMsg string = "Not equals"

func TestToAndFromDomainAddress(t *testing.T) {
	restAddress := newAddress()
	domainAddress := mapper.toDomainAddress(restAddress)
	backToRest := mapper.fromDomainAddress(domainAddress)
	if !reflect.DeepEqual(restAddress, backToRest) {
		t.Fatalf(notEqualsErrorMsg)
	}
}

func TestToAndFromDomainMenu(t *testing.T) {
	restMenu := newMenu()
	domainMenu := mapper.toDomainMenu(restMenu)
	backToRest := mapper.fromDomainMenu(domainMenu)
	if !reflect.DeepEqual(restMenu, backToRest) {
		t.Fatalf(notEqualsErrorMsg)
	}
}

func TestToAndFromDomainRestaurant(t *testing.T) {
	restRestaurant := newRestaurant()
	domainRestaurant := mapper.toDomainRestaurant(restRestaurant)
	backToRest := mapper.fromDomainRestaurant(domainRestaurant)
	if !reflect.DeepEqual(restRestaurant, backToRest) {
		t.Fatalf(notEqualsErrorMsg)
	}
}

func TestFromDomainRestaurants(t *testing.T) {
	restRestaurants := []*Restaurant{newRestaurant(), newRestaurant()}
	domainRestaurants := mapper.toDomainRestaurants(restRestaurants)
	backToRest := mapper.fromDomainRestaurants(domainRestaurants)
	if !reflect.DeepEqual(restRestaurants, backToRest) {
		t.Fatalf(notEqualsErrorMsg)
	}
}

func newAddress() *Address {
	return &Address{Street: "aStreet", City: "aCity", State: "aState", Zip: "aZip"}
}

func newMenu() *Menu {
	item1 := MenuItem{Id: 1, Name: "Name1", Price: "12.99"}
	item2 := MenuItem{Id: 2, Name: "Name2", Price: "4.55"}
	return &Menu{Items: []MenuItem{item1, item2}}
}

func newRestaurant() *Restaurant {
	return &Restaurant{Name: "aName", Address: newAddress(), Menu: newMenu()}
}
