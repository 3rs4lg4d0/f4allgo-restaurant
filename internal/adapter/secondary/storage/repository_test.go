package storage

import (
	"context"
	"errors"
	"f4allgo-restaurant/internal/adapter/secondary/eventpublisher/outbox"
	"f4allgo-restaurant/internal/boot"
	"f4allgo-restaurant/internal/core/domain"
	"f4allgo-restaurant/internal/core/port"
	"f4allgo-restaurant/test"
	"fmt"
	"os"
	"testing"

	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/stretchr/testify/assert"
	pgcontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	database             *pgcontainer.PostgresContainer
	db                   *gorm.DB
	trManager            *manager.Manager
	restaurantRepository port.RestaurantRepository
)

var mapper Mapper = DefaultMapper{}

func TestMain(m *testing.M) {
	var err error
	var dsn string
	ctx := context.Background()

	database, err = test.InitApplicationDatabase(ctx)
	if err != nil {
		fmt.Printf("A problem occurred initializing the database: %v", err)
		os.Exit(1)
	}

	dsn, err = database.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Printf("A problem occurred getting the connection string: %v", err)
		os.Exit(1)
	}

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect to database")
	}

	trManager = boot.GetTransactionManager(db)
	restaurantRepository = NewRestaurantPostgresRepository(db, trmgorm.DefaultCtxGetter, nil)

	code := m.Run()

	test.TerminateApplicationDatabase(ctx, database)
	os.Exit(code)
}

func TestRestaurantPostgresRepository_FindAll(t *testing.T) {
	ctx := context.Background()
	restaurants, err := restaurantRepository.FindAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, restaurants, 3)
	assert.Equal(t, uint64(1000), restaurants[0].Id)
}

func TestRestaurantPostgresRepository_FindById(t *testing.T) {
	ctx := context.Background()
	restaurant, err := restaurantRepository.FindById(ctx, uint64(1000))
	assert.NoError(t, err)
	assert.Equal(t, "restaurant1", restaurant.Name)
	assert.Equal(t, "city1", restaurant.Address.City())
	assert.Equal(t, "state1", restaurant.Address.State())
	assert.Equal(t, "street1", restaurant.Address.Street())
	assert.Equal(t, "zip1", restaurant.Address.Zip())
	assert.Equal(t, 0, len(restaurant.Menu.GetItems()))

	restaurant, err = restaurantRepository.FindById(ctx, uint64(5555))
	assert.NoError(t, err)
	assert.Nil(t, restaurant)
}

func TestRestaurantPostgresRepository_Save(t *testing.T) {
	ctx := context.Background()
	trManager.Do(ctx, func(ctx context.Context) error {
		// Create some fake restaurants.
		restaurant1 := newRestaurant()
		restaurant2 := newRestaurant()
		err := restaurantRepository.Save(ctx, restaurant1)
		assert.NoError(t, err)
		err = restaurantRepository.Save(ctx, restaurant2)
		assert.NoError(t, err)

		// Find restaurants.
		restaurants, err := restaurantRepository.FindAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, restaurants, 5)

		return errors.New("rollback")
	})
}

func TestRestaurantPostgresRepository_Update(t *testing.T) {
	ctx := context.Background()
	trManager.Do(ctx, func(ctx context.Context) error {
		restaurant, err := restaurantRepository.FindById(ctx, uint64(1000))
		assert.NoError(t, err)

		newMenu := newMenu()
		err = restaurant.UpdateMenu(mapper.toDomainMenu(newMenu))
		assert.NoError(t, err)

		err = restaurantRepository.Update(ctx, restaurant)
		assert.NoError(t, err)

		restaurants, err := restaurantRepository.FindAll(ctx)
		assert.NoError(t, err)
		for _, r := range restaurants {
			if r.Id == uint64(1000) {
				assert.Equal(t, "restaurant1", restaurant.Name)
				assert.Equal(t, "city1", restaurant.Address.City())
				assert.Equal(t, "state1", restaurant.Address.State())
				assert.Equal(t, "street1", restaurant.Address.Street())
				assert.Equal(t, "zip1", restaurant.Address.Zip())
				assert.Equal(t, 2, len(restaurant.Menu.GetItems()))
			}
		}

		return errors.New("rollback")
	})
}

func TestRestaurantPostgresRepository_Deletex(t *testing.T) {
	ctx := context.Background()
	trManager.Do(ctx, func(ctx context.Context) error {
		restaurants, err := restaurantRepository.FindAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(restaurants))

		err = restaurantRepository.Delete(ctx, uint64(1000))
		restaurants, _ = restaurantRepository.FindAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(restaurants))

		// Test that Delete is silent when trying to delete an unexisting row.
		err = restaurantRepository.Delete(ctx, uint64(1001))
		assert.NoError(t, err)

		return errors.New("rollback")
	})
}

func TestGormExamples(t *testing.T) {
	var outbox []*outbox.Outbox
	err := db.Find(&outbox).Error
	assert.NoError(t, err)
	assert.Equal(t, 5, len(outbox))

	err = db.Delete(outbox).Error
	assert.NoError(t, err)

	db.Find(&outbox)
	assert.Equal(t, 0, len(outbox))
}

func newAddress() *Address {
	return &Address{Street: "aStreet", City: "aCity", State: "aState", Zip: "aZip"}
}

func newMenu() []*MenuItem {
	item1 := MenuItem{Id: 1, Name: "aName1", Price: "6.99"}
	item2 := MenuItem{Id: 2, Name: "aName2", Price: "7.99"}
	return []*MenuItem{&item1, &item2}
}

func newRestaurant() *domain.Restaurant {
	return mapper.toDomainRestaurant(&Restaurant{Name: "aName", Address: newAddress(), Menu: newMenu()})
}
