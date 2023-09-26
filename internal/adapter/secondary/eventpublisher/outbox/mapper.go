package outbox

import (
	"f4allgo-restaurant/internal/adapter/secondary/eventpublisher/avro"
	"f4allgo-restaurant/internal/core/domain"
)

type Mapper interface {

	// fromRestaurantCreated maps a domain.RestaurantCreated event into an outbox row.
	fromRestaurantCreated(restaurant *domain.RestaurantCreated) *avro.RestaurantCreatedAvro

	// fromRestaurantDeleted maps a domain.RestaurantDeleted event into an outbox row.
	fromRestaurantDeleted(address *domain.RestaurantDeleted) *avro.RestaurantDeletedAvro

	// fromRestaurantMenuUpdated maps a domain.RestaurantMenuUpdated into an outbox row.
	fromRestaurantMenuUpdated(menu *domain.RestaurantMenuUpdated) *avro.RestaurantMenuUpdatedAvro
}

// DefaultMapper is the default implementation of Mapper.
// It provides a plain manual mapping without using any third party library.
type DefaultMapper struct{}

// Interface compliance verification.
var _ Mapper = (*DefaultMapper)(nil)

func (dm DefaultMapper) fromRestaurantCreated(event *domain.RestaurantCreated) *avro.RestaurantCreatedAvro {
	return &avro.RestaurantCreatedAvro{
		Id:      int64(event.Restaurant.Id),
		Name:    event.Restaurant.Name,
		Address: dm.fromDomainAddress(event.Restaurant.Address),
		Menu:    dm.fromDomainMenu(event.Restaurant.Menu),
	}
}

func (dm DefaultMapper) fromRestaurantDeleted(event *domain.RestaurantDeleted) *avro.RestaurantDeletedAvro {
	return &avro.RestaurantDeletedAvro{
		RestaurantId: int64(event.RestaurantId),
	}
}

func (dm DefaultMapper) fromRestaurantMenuUpdated(event *domain.RestaurantMenuUpdated) *avro.RestaurantMenuUpdatedAvro {
	return &avro.RestaurantMenuUpdatedAvro{
		RestaurantId: int64(event.RestaurantId),
		Menu:         dm.fromDomainMenu(event.Menu),
	}
}

func (dm DefaultMapper) fromDomainAddress(address *domain.Address) avro.AdressAvro {
	return avro.AdressAvro{
		Street: address.Street(),
		City:   address.City(),
		State:  address.State(),
		Zip:    address.Zip(),
	}
}
func (dm DefaultMapper) fromDomainMenu(menu *domain.Menu) avro.MenuAvro {
	avroItems := []avro.MenuItemAvro{}
	for _, item := range menu.GetItems() {
		avroItems = append(avroItems, avro.MenuItemAvro{Id: int64(item.GetId()), Name: item.GetName(), Price: []byte(item.GetPrice().Text('f', 2))})
	}
	return avro.MenuAvro{
		Items: avroItems,
	}
}
