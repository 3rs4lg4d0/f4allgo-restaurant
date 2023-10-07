package port

import (
	"context"
	"f4allgo-restaurant/internal/core/domain"
)

// RestaurantRepository manages persistent operations on restaurants dealing
// with an external storage system.
type RestaurantRepository interface {

	// FindAll restrieves all the registered restaurants with their menus.
	FindAll(ctx context.Context, offset int, limit int) ([]*domain.Restaurant, int64, error)

	// FindById retrieves a particular restaurant by its identifier. The method
	// allows the client to decide if the restaurant's menu should be fetched or not.
	FindById(ctx context.Context, restaurantId int64, fetchMenu bool) (*domain.Restaurant, error)

	// Save persists a restaurant in the external storage and sets its generated identifier
	// to the restaurant instance input argument.
	Save(ctx context.Context, restaurant *domain.Restaurant) error

	// Update updates a restaurant and its relations and return the number of rows affected
	// (this value is important because all the operations that assume the existence of the
	// restaurant in the first place can be wrong and we want to be able to detect it and
	// inform it properly from the service layer) Call this operation when you want to update
	// a restaurant's menu.
	Update(ctx context.Context, restaurant *domain.Restaurant) (int64, error)

	// Delete deletes a restaurant from the external storage and return the number of rows
	// affected.
	Delete(ctx context.Context, restaurantId int64) (int64, error)
}

// DomainEventPublisher publishes domain events to the outside world dealing with a
// message broker.
type DomainEventPublisher interface {

	// Publish publishes a domain event to the outside world.
	Publish(ctx context.Context, e domain.DomainEvent) error
}
