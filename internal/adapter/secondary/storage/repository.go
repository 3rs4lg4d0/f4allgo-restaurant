package storage

import (
	"context"
	"f4allgo-restaurant/internal/core/domain"
	"f4allgo-restaurant/internal/core/port"
	"fmt"

	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"
	tally "github.com/uber-go/tally/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Postgres implementation of the secondary port RestaurantRepository. It uses
// GORM to persist and retrieve restaurants from Postgres.
type RestaurantPostgresRepository struct {
	mapper    Mapper
	db        *gorm.DB
	ctxGetter *trmgorm.CtxGetter
	timer     tally.Timer
}

// Interface compliance verification.
var _ port.RestaurantRepository = (*RestaurantPostgresRepository)(nil)

func NewRestaurantPostgresRepository(db *gorm.DB, ctxGetter *trmgorm.CtxGetter, timer tally.Timer) *RestaurantPostgresRepository {
	return &RestaurantPostgresRepository{mapper: DefaultMapper{}, db: db, ctxGetter: ctxGetter, timer: timer}
}

// FindAll restrieves all the registered restaurants.
func (r *RestaurantPostgresRepository) FindAll(ctx context.Context, offset uint32, limit uint8) ([]*domain.Restaurant, error) {
	var restaurants []*Restaurant
	if err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).Preload("Menu").Offset(int(offset)).Limit(int(limit)).Find(&restaurants).Error; err != nil {
		return nil, err
	}

	return r.mapper.toDomainRestaurants(restaurants), nil
}

// FindById retrieves a particular restaurant by its identifier. If the restaurant
// doesn't exist it doesn't fail, but return a nil restaurant.
func (r *RestaurantPostgresRepository) FindById(ctx context.Context, restaurantId uint64) (*domain.Restaurant, error) {
	var restaurants [1]*Restaurant
	if err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).Where("id = ?", restaurantId).Find(&restaurants).Error; err != nil {
		return nil, err
	}

	// If the restaurant doesn't exist, return nil.
	if restaurants[0] == nil {
		return nil, nil
	}

	return r.mapper.toDomainRestaurant(restaurants[0]), nil
}

// Save persists a restaurant in the database.
func (r *RestaurantPostgresRepository) Save(ctx context.Context, restaurant *domain.Restaurant) error {
	restaurantDto := r.mapper.fromDomainRestaurant(restaurant)
	if err := r.executeWithTimer(func() error {
		return r.ctxGetter.DefaultTrOrDB(ctx, r.db).Clauses(clause.OnConflict{DoNothing: true}).Create(restaurantDto).Error
	}); err != nil {
		return fmt.Errorf("when trying to create the restaurant: %w", err)
	}
	restaurant.Id = restaurantDto.ID
	return nil
}

// Update updates a restaurant and its relations.
func (r *RestaurantPostgresRepository) Update(ctx context.Context, restaurant *domain.Restaurant) error {
	// Delete previous relation in the database.
	err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).Where("restaurant_id = ?", restaurant.Id).Delete(&MenuItem{}).Error
	if err != nil {
		return fmt.Errorf("an error ocurrend when trying to delete the menu items related to the restaurant: %w", err)
	}

	// Persist the aggregate again.
	return r.ctxGetter.DefaultTrOrDB(ctx, r.db).Save(r.mapper.fromDomainRestaurant(restaurant)).Error
}

// Delete deletes a restaurant and its relations from the database.
func (r *RestaurantPostgresRepository) Delete(ctx context.Context, restaurantId uint64) error {
	return r.ctxGetter.DefaultTrOrDB(ctx, r.db).Select(clause.Associations).Delete(&Restaurant{ID: restaurantId}).Error
}

// executeWithTimer executes a function using a tally timer if present.
func (r *RestaurantPostgresRepository) executeWithTimer(fn func() error) error {
	if r.timer != nil {
		tsw := r.timer.Start()
		defer tsw.Stop()
	}
	return fn()
}
