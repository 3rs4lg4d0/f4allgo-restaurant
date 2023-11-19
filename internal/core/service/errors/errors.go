package errors

// RestaurantNotFoundError is returned when searching for a particular restaurant
// in the database and no results are found matching the criteria.
type RestaurantNotFoundError struct{}

func NewRestaurantNotFoundError() *RestaurantNotFoundError {
	return &RestaurantNotFoundError{}
}

func (r *RestaurantNotFoundError) Error() string {
	return ""
}

// RepositoryError is returned when an unexpected error occurr using
// repositories.
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
// event publishers.
type EventPublisherError struct {
	err error
}

func NewEventPublisherError(err error) *EventPublisherError {
	return &EventPublisherError{err: err}
}

func (e *EventPublisherError) Error() string {
	return e.err.Error()
}

// CoreError is returned when an unexpected error occur executing a core
// domain operation.
type CoreError struct {
	err error
}

func NewCoreError(err error) *CoreError {
	return &CoreError{err: err}
}

func (e *CoreError) Error() string {
	return e.err.Error()
}
