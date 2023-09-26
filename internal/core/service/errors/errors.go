package errors

// RestaurantNotFoundError is returned when searching for a particular restaurant
// in the database and no results are found matching the criteria.
type RestaurantNotFoundError struct {
	err error
}

func NewRestaurantNotFoundError(err error) *RestaurantNotFoundError {
	return &RestaurantNotFoundError{err: err}
}

func (r *RestaurantNotFoundError) Error() string {
	return r.err.Error()
}

// RepositoryError is returned when an unexpected error occurr using
// repositories. Use it as the default option if no specific error applies.
type RepositoryError struct {
	err error
}

func NewRepositoryError(err error) *RepositoryError {
	return &RepositoryError{err: err}
}

func (r *RepositoryError) Error() string {
	return r.err.Error()
}

// EventPublisherError is returned when an unexpected error occur using
// event publishers. Use it as the default option if no specific error applies.
type EventPublisherError struct {
	err error
}

func NewEventPublisherError(err error) *EventPublisherError {
	return &EventPublisherError{err: err}
}

func (e *EventPublisherError) Error() string {
	return e.err.Error()
}
