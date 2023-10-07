package storage

import (
	"context"
	"f4allgo-restaurant/internal/core/domain"
	"f4allgo-restaurant/internal/core/port"

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
func (r *RestaurantPostgresRepository) FindAll(ctx context.Context, offset int, limit int) ([]*domain.Restaurant, int64, error) {
	var restaurants []*Restaurant
	var total int64
	if err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).Preload("Menu").Offset(offset).Limit(limit).Find(&restaurants).Error; err != nil {
		return nil, 0, err
	}
	if err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).Model(&Restaurant{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return r.mapper.toDomainRestaurants(restaurants), total, nil
}

// FindById retrieves a particular restaurant by its identifier. The method allows the client to decide if
// the restaurant's menu should be fetched or not.
func (r *RestaurantPostgresRepository) FindById(ctx context.Context, restaurantId int64, fetchMenu bool) (*domain.Restaurant, error) {
	var restaurant *Restaurant
	var err error

	if fetchMenu {
		err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).Preload("Menu").First(&restaurant, restaurantId).Error
	} else {
		err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).First(&restaurant, restaurantId).Error
	}

	if err != nil {
		return nil, err
	}

	return r.mapper.toDomainRestaurant(restaurant), nil
}

// Save persists a restaurant in the database.
func (r *RestaurantPostgresRepository) Save(ctx context.Context, restaurant *domain.Restaurant) error {
	restaurantDto := r.mapper.fromDomainRestaurant(restaurant)
	if err := r.executeWithTimer(func() error {
		return r.ctxGetter.DefaultTrOrDB(ctx, r.db).Clauses(clause.OnConflict{DoNothing: true}).Create(restaurantDto).Error
	}); err != nil {
		return err
	}
	restaurant.Id = restaurantDto.ID
	return nil
}

// Update updates a restaurant and its relations.
func (r *RestaurantPostgresRepository) Update(ctx context.Context, restaurant *domain.Restaurant) (int64, error) {
	// Delete previous menu items.
	err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).Where("restaurant_id = ?", restaurant.Id).Delete(&MenuItem{}).Error
	if err != nil {
		return 0, err
	}

	// Persist the aggregate again. If the aggregate doesn't exist in the database in the first place
	// en error will happen (inserting the menu items or the restaurant, it doesn't matter) and the
	// affected rows will be zero (that's exactly what we want).
	result := r.ctxGetter.DefaultTrOrDB(ctx, r.db).Save(r.mapper.fromDomainRestaurant(restaurant))

	return result.RowsAffected, result.Error
}

// Delete deletes a restaurant and its relations from the database.
func (r *RestaurantPostgresRepository) Delete(ctx context.Context, restaurantId int64) (int64, error) {
	result := r.ctxGetter.DefaultTrOrDB(ctx, r.db).Select(clause.Associations).Delete(&Restaurant{ID: restaurantId})
	return result.RowsAffected, result.Error
}

// executeWithTimer executes a function using a tally timer if present.
func (r *RestaurantPostgresRepository) executeWithTimer(fn func() error) error {
	if r.timer != nil {
		tsw := r.timer.Start()
		defer tsw.Stop()
	}
	return fn()
}
