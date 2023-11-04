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
	timers    map[timerEnum]tally.Timer
}

// Interface compliance verification.
var _ port.RestaurantRepository = (*RestaurantPostgresRepository)(nil)

// Custom type for timer types.
type timerEnum int

// Enumeration of timers for repository operations.
const (
	findAll timerEnum = iota
	findById
	save
	update
	delete
)

func NewRestaurantPostgresRepository(db *gorm.DB, ctxGetter *trmgorm.CtxGetter, scope tally.Scope) *RestaurantPostgresRepository {
	var timers map[timerEnum]tally.Timer
	if scope != nil {
		FindAll := scope.Tagged(map[string]string{"repository": "restaurant", "operation": "FindAll"}).Timer("repository_latencies")
		FindById := scope.Tagged(map[string]string{"repository": "restaurant", "operation": "FindById"}).Timer("repository_latencies")
		Save := scope.Tagged(map[string]string{"repository": "restaurant", "operation": "Save"}).Timer("repository_latencies")
		Update := scope.Tagged(map[string]string{"repository": "restaurant", "operation": "Update"}).Timer("repository_latencies")
		Delete := scope.Tagged(map[string]string{"repository": "restaurant", "operation": "Delete"}).Timer("repository_latencies")

		timers = make(map[timerEnum]tally.Timer)
		timers[findAll] = FindAll
		timers[findById] = FindById
		timers[save] = Save
		timers[update] = Update
		timers[delete] = Delete
	}
	return &RestaurantPostgresRepository{mapper: DefaultMapper{}, db: db, ctxGetter: ctxGetter, timers: timers}
}

// FindAll restrieves all the registered restaurants.
func (r *RestaurantPostgresRepository) FindAll(ctx context.Context, offset int, limit int) ([]*domain.Restaurant, int64, error) {
	var restaurants []*Restaurant
	var total int64

	if err := r.executeWithTimer(findAll, func() error {
		if err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).Preload("Menu").Order("id ASC").Offset(offset).Limit(limit).Find(&restaurants).Error; err != nil {
			return err
		}
		if err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).Model(&Restaurant{}).Count(&total).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, 0, err
	}

	return r.mapper.toDomainRestaurants(restaurants), total, nil
}

// FindById retrieves a particular restaurant by its identifier. The method allows the client to decide if
// the restaurant's menu should be fetched or not.
func (r *RestaurantPostgresRepository) FindById(ctx context.Context, restaurantId int64, fetchMenu bool) (*domain.Restaurant, error) {
	var restaurant *Restaurant
	if err := r.executeWithTimer(findById, func() error {
		if fetchMenu {
			return r.ctxGetter.DefaultTrOrDB(ctx, r.db).Preload("Menu").First(&restaurant, restaurantId).Error
		} else {
			return r.ctxGetter.DefaultTrOrDB(ctx, r.db).First(&restaurant, restaurantId).Error
		}
	}); err != nil {
		return nil, err
	}

	return r.mapper.toDomainRestaurant(restaurant), nil
}

// Save persists a restaurant in the database.
func (r *RestaurantPostgresRepository) Save(ctx context.Context, restaurant *domain.Restaurant) error {
	restaurantDto := r.mapper.fromDomainRestaurant(restaurant)
	if err := r.executeWithTimer(save, func() error {
		return r.ctxGetter.DefaultTrOrDB(ctx, r.db).Clauses(clause.OnConflict{DoNothing: true}).Create(restaurantDto).Error
	}); err != nil {
		return err
	}
	restaurant.Id = restaurantDto.ID
	return nil
}

// Update updates a restaurant and its relations.
func (r *RestaurantPostgresRepository) Update(ctx context.Context, restaurant *domain.Restaurant) (int64, error) {
	var result *gorm.DB
	if err := r.executeWithTimer(update, func() error {
		// Delete previous menu items.
		err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).Where("restaurant_id = ?", restaurant.Id).Delete(&MenuItem{}).Error
		if err != nil {
			return err
		}

		// Persist the aggregate again. If the aggregate doesn't exist in the database in the first place
		// en error will happen (inserting the menu items or the restaurant, it doesn't matter) and the
		// affected rows will be zero (that's exactly what we want).
		result = r.ctxGetter.DefaultTrOrDB(ctx, r.db).Save(r.mapper.fromDomainRestaurant(restaurant))
		return result.Error
	}); err != nil {
		return 0, err
	}

	return result.RowsAffected, nil
}

// Delete deletes a restaurant and its relations from the database.
func (r *RestaurantPostgresRepository) Delete(ctx context.Context, restaurantId int64) (int64, error) {
	var result *gorm.DB
	if err := r.executeWithTimer(delete, func() error {
		result = r.ctxGetter.DefaultTrOrDB(ctx, r.db).Select(clause.Associations).Delete(&Restaurant{ID: restaurantId})
		return result.Error
	}); err != nil {
		return 0, err
	}

	return result.RowsAffected, nil
}

// executeWithTimer executes a function using a tally timer if present.
func (r *RestaurantPostgresRepository) executeWithTimer(t timerEnum, fn func() error) error {
	if r.timers[t] != nil {
		tsw := r.timers[t].Start()
		defer tsw.Stop()
	}
	return fn()
}
