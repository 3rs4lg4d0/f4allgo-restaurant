package domain

type DomainEvent interface {
	GetType() string
}

// --------------------------------------------------------------------------------
// Event :: RestaurantCreated
// --------------------------------------------------------------------------------

// RestaurantCreated event is raised every time a new restaurant is created.
type RestaurantCreated struct {
	Restaurant *Restaurant
}

// Interface compliance verification.
var _ DomainEvent = (*RestaurantCreated)(nil)

func NewRestaurantCreated(restaurant *Restaurant) *RestaurantCreated {
	return &RestaurantCreated{Restaurant: restaurant}
}

func (r *RestaurantCreated) GetType() string {
	return "RestaurantCreated"
}

// --------------------------------------------------------------------------------
// Event :: RestaurantDeleted
// --------------------------------------------------------------------------------

// RestaurantDeleted event is raised every time a restaurant is deleted.
type RestaurantDeleted struct {
	RestaurantId int64
}

// Interface compliance verification.
var _ DomainEvent = (*RestaurantDeleted)(nil)

func NewRestaurantDeleted(restaurantId int64) *RestaurantDeleted {
	return &RestaurantDeleted{RestaurantId: restaurantId}
}

func (r *RestaurantDeleted) GetType() string {
	return "RestaurantDeleted"
}

// --------------------------------------------------------------------------------
// Event :: RestaurantMenuUpdated
// --------------------------------------------------------------------------------

// RestaurantMenuUpdated event is raised every time a restaurant updates its menu.
type RestaurantMenuUpdated struct {
	RestaurantId int64
	Menu         *Menu
}

// Interface compliance verification.
var _ DomainEvent = (*RestaurantMenuUpdated)(nil)

func NewRestaurantMenuUpdated(restaurantId int64, menu *Menu) *RestaurantMenuUpdated {
	return &RestaurantMenuUpdated{RestaurantId: restaurantId, Menu: menu}
}

func (r *RestaurantMenuUpdated) GetType() string {
	return "RestaurantMenuUpdated"
}
