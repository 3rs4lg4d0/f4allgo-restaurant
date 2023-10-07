package storage

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToAndFromDomainRestaurant(t *testing.T) {
	mapper := DefaultMapper{}
	type args struct {
		restaurants []*Restaurant
	}
	testcases := []struct {
		name string
		args args
	}{
		{
			name: "map a slice of storage restaurants",
			args: args{
				restaurants: []*Restaurant{newStorageRestaurant(), newStorageRestaurant()},
			},
		},
		{
			name: "map a nil slice",
			args: args{
				restaurants: nil,
			},
		},
		{
			name: "map a slice of nils",
			args: args{
				restaurants: []*Restaurant{nil, nil},
			},
		},
		{
			name: "map a slice of restaurants without address",
			args: args{
				restaurants: []*Restaurant{newStorageRestaurantWithoutAddress(), newStorageRestaurantWithoutAddress()},
			},
		},
		{
			name: "map a slice of restaurants without menu",
			args: args{
				restaurants: []*Restaurant{newStorageRestaurantWithoutMenu(), newStorageRestaurantWithoutMenu()},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			domainRestaurants := mapper.toDomainRestaurants(tc.args.restaurants)
			backToStorage := mapper.fromDomainRestaurants(domainRestaurants)

			assert.True(t, reflect.DeepEqual(tc.args.restaurants, backToStorage))
		})
	}
}

// --------------------------------------------------------------------------------
// Utility functions to create restaurants and menus.
// --------------------------------------------------------------------------------

func newStorageRestaurant() *Restaurant {
	return &Restaurant{Name: "aName", Address: newStorageAddress(), Menu: newStorageMenu()}
}

func newStorageRestaurantWithoutAddress() *Restaurant {
	return &Restaurant{Name: "aName", Menu: newStorageMenu()}
}

func newStorageRestaurantWithoutMenu() *Restaurant {
	return &Restaurant{Name: "aName", Address: newStorageAddress()}
}

func newStorageAddress() *Address {
	return &Address{Street: "aStreet", City: "aCity", State: "aState", Zip: "aZip"}
}

func newStorageMenu() []*MenuItem {
	item1 := MenuItem{Id: 1, Name: "Name1", Price: "12.99"}
	item2 := MenuItem{Id: 2, Name: "Name2", Price: "4.55"}
	return []*MenuItem{&item1, &item2}
}
