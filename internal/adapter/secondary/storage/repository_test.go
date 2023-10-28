package storage

import (
	"context"
	"errors"
	"f4allgo-restaurant/internal/boot"
	"f4allgo-restaurant/internal/core/domain"
	"f4allgo-restaurant/internal/core/port"
	"f4allgo-restaurant/test"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"
	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/stretchr/testify/assert"
	pgcontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const ROLLBACK_PLEASE string = "rollback to keep database clean from testcase to testcase"

var (
	database             *pgcontainer.PostgresContainer
	db                   *gorm.DB
	trManager            trm.Manager
	restaurantRepository port.RestaurantRepository
)

var mapper Mapper = DefaultMapper{}

// TestMain prepares the database setup needed to run these tests. As you can see
// the database layer is tested against a real Postgres containerized instance, but
// for some specific cases (mostly to simulate errors) a sqlmock instance is used.
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
	restaurantRepository = NewRestaurantPostgresRepository(db, trmgorm.DefaultCtxGetter, test.NewNopTallyTimer())

	code := m.Run()

	err = test.TerminateApplicationDatabase(ctx, database)
	if err != nil {
		fmt.Printf("an error ocurred terminating the database container: %w", err)
	}
	os.Exit(code)
}

func TestFindAll(t *testing.T) {
	type args struct {
		offset int
		limit  int
	}
	testcases := []struct {
		name               string
		args               args
		mockExpectations   func(sqlmock.Sqlmock)
		wantLen            int
		wantTotal          int64
		wantFirstElementId int64
		wantErr            bool
		wantErrMsg         string
	}{
		{
			name: "offset 0 and limit 100",
			args: args{
				offset: 0,
				limit:  100,
			},
			wantLen:            3,
			wantTotal:          3,
			wantFirstElementId: 1000,
			wantErr:            false,
		},
		{
			name: "offset 0 and limit 2",
			args: args{
				offset: 0,
				limit:  2,
			},
			wantLen:            2,
			wantTotal:          3,
			wantFirstElementId: 1000,
			wantErr:            false,
		},
		{
			name: "offset 1 and limit 2",
			args: args{
				offset: 1,
				limit:  2,
			},
			wantLen:            2,
			wantTotal:          3,
			wantFirstElementId: 2000,
			wantErr:            false,
		},
		{
			name: "simulate error when fetching restaurants with mocked DB",
			args: args{
				offset: 0,
				limit:  100,
			},
			mockExpectations: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT .+ FROM "restaurant" .+`).WillReturnError(errors.New("error#1"))
			},
			wantErr:    true,
			wantErrMsg: "error#1",
		},
		{
			name: "simulate error when counting total restaurants with mocked DB",
			args: args{
				offset: 0,
				limit:  100,
			},
			mockExpectations: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT .+ FROM "restaurant" .+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "city", "state", "street", "zip"}).AddRow(1000, "restaurant1", "city1", "state1", "street1", "zip1"))
				mock.ExpectQuery(`SELECT .+ FROM "menu_item" .+`).WillReturnRows(sqlmock.NewRows(make([]string, 0)))
				mock.ExpectQuery(`SELECT count.+ FROM "restaurant`).WillReturnError(errors.New("error#2"))
			},
			wantErr:    true,
			wantErrMsg: "error#2",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var repository = restaurantRepository
			if tc.mockExpectations != nil {
				var mock sqlmock.Sqlmock
				repository, _, mock = createMockRepository()
				tc.mockExpectations(mock)
			}
			restaurants, total, err := repository.FindAll(context.Background(), tc.args.offset, tc.args.limit)
			if !tc.wantErr {
				assert.NoError(t, err)
				assert.Len(t, restaurants, tc.wantLen)
				assert.Equal(t, tc.wantTotal, total)
				assert.Equal(t, tc.wantFirstElementId, restaurants[0].Id)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tc.wantErrMsg, err.Error())
			}
		})
	}
}

func TestFindById(t *testing.T) {
	type args struct {
		restaurantId int64
		fetchMenu    bool
	}
	testcases := []struct {
		name           string
		args           args
		wantRestaurant *Restaurant
		wantErr        bool
		wantErrMsg     string
	}{
		{
			name: "find a restaurant that exists fetching its menu",
			args: args{
				restaurantId: 1000,
				fetchMenu:    true,
			},
			wantRestaurant: newTestRestaurant(),
			wantErr:        false,
		},
		{
			name: "find a restaurant that doesn't exist fetching its menu",
			args: args{
				restaurantId: 1001,
				fetchMenu:    true,
			},
			wantRestaurant: nil,
			wantErr:        true,
			wantErrMsg:     "record not found",
		},
		{
			name: "find a restaurant that exists not fetching its menu",
			args: args{
				restaurantId: 1000,
				fetchMenu:    false,
			},
			wantRestaurant: newTestRestaurantWithoutMenu(),
			wantErr:        false,
		},
		{
			name: "find a restaurant that doesn't exist not fetching its menu",
			args: args{
				restaurantId: 1001,
				fetchMenu:    false,
			},
			wantRestaurant: nil,
			wantErr:        true,
			wantErrMsg:     "record not found",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			restaurant, err := restaurantRepository.FindById(context.Background(), tc.args.restaurantId, tc.args.fetchMenu)
			assert.Equal(t, tc.wantRestaurant, mapper.fromDomainRestaurant(restaurant))
			if !tc.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tc.wantErrMsg, err.Error())
			}
		})
	}
}

func TestSave(t *testing.T) {
	type args struct {
		restaurant *domain.Restaurant
	}
	testcases := []struct {
		name             string
		args             args
		mockExpectations func(sqlmock.Sqlmock)
		wantErr          bool
		wantErrMsg       string
	}{
		{
			name: "save a restaurant successfully",
			args: args{
				restaurant: mapper.toDomainRestaurant(newTestRestaurantWithoutId()),
			},
			wantErr: false,
		},
		{
			name: "simulate error with mocked DB",
			args: args{
				restaurant: mapper.toDomainRestaurant(newTestRestaurantWithoutId()),
			},
			mockExpectations: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO .+").WillReturnError(errors.New("error#3"))
				mock.ExpectRollback()
			},
			wantErr:    true,
			wantErrMsg: "error#3",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			var repository = restaurantRepository
			var trm = trManager
			if tc.mockExpectations != nil {
				var mock sqlmock.Sqlmock
				repository, trm, mock = createMockRepository()
				tc.mockExpectations(mock)
			}
			err := trm.Do(ctx, func(ctx context.Context) error {
				assert.Equal(t, int64(0), tc.args.restaurant.Id)
				err := repository.Save(ctx, tc.args.restaurant)
				if !tc.wantErr {
					assert.NoError(t, err)
					_, total, _ := repository.FindAll(ctx, 0, 100)
					assert.Equal(t, int64(4), total)
					actualRestaurant, _ := repository.FindById(ctx, tc.args.restaurant.Id, true)
					assert.True(t, reflect.DeepEqual(tc.args.restaurant, actualRestaurant))
				} else {
					assert.Error(t, err)
					assert.Equal(t, tc.wantErrMsg, err.Error())
				}

				// Returning an error instructs the transaction manager to
				// rollback the transaction so we keep clean the database for
				// the other tests.
				return errors.New(ROLLBACK_PLEASE)
			})

			assert.Error(t, err)
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		restaurant *domain.Restaurant
	}
	testcases := []struct {
		name             string
		args             args
		mockExpectations func(sqlmock.Sqlmock)
		wantRowsAffected int64
		wantMenuItems    int
		wantErr          bool
		wantErrMsg       string
	}{
		{
			name: "update restaurant name",
			args: args{
				restaurant: func() *domain.Restaurant {
					r := mapper.toDomainRestaurant(newTestRestaurant())
					r.Name = "A different name"
					return r
				}(),
			},
			wantRowsAffected: 1,
			wantMenuItems:    3,
			wantErr:          false,
		},
		{
			name: "update restaurant menu",
			args: args{
				restaurant: func() *domain.Restaurant {
					r := mapper.toDomainRestaurant(newTestRestaurant())
					r.Menu = mapper.toDomainMenu(newTestMenu()[0:2])
					return r
				}(),
			},
			wantRowsAffected: 1,
			wantMenuItems:    2,
			wantErr:          false,
		},
		{
			name: "update a restaurant that doesn't exist",
			args: args{
				restaurant: func() *domain.Restaurant {
					r := mapper.toDomainRestaurant(newTestRestaurant())
					r.Id = 1001
					return r
				}(),
			},
			wantRowsAffected: 0,
			wantErr:          true,
		},
		{
			name: "simulate error when deleting menu with mocked DB",
			args: args{
				restaurant: mapper.toDomainRestaurant(newTestRestaurant()),
			},
			mockExpectations: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM .+").WillReturnError(errors.New("error#4"))
				mock.ExpectRollback()
			},
			wantErr:    true,
			wantErrMsg: "error#4",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			var repository = restaurantRepository
			var trm = trManager
			if tc.mockExpectations != nil {
				var mock sqlmock.Sqlmock
				repository, trm, mock = createMockRepository()
				tc.mockExpectations(mock)
			}
			err := trm.Do(ctx, func(ctx context.Context) error {
				ra, err := repository.Update(ctx, tc.args.restaurant)
				if !tc.wantErr {
					assert.NoError(t, err)
					assert.Equal(t, tc.wantRowsAffected, ra)
					if ra > 0 {
						actualRestaurant, _ := repository.FindById(ctx, tc.args.restaurant.Id, true)
						assert.True(t, reflect.DeepEqual(tc.args.restaurant, actualRestaurant))
						assert.Len(t, actualRestaurant.Menu.GetItems(), tc.wantMenuItems)
					}
				} else {
					assert.Error(t, err)
					if len(tc.wantErrMsg) > 0 {
						assert.Equal(t, tc.wantErrMsg, err.Error())
					}
				}

				return errors.New(ROLLBACK_PLEASE)
			})

			assert.Error(t, err)
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		restaurantId int64
	}
	testcases := []struct {
		name             string
		args             args
		wantRowsAffected int64
	}{
		{
			name: "delete restaurant",
			args: args{
				restaurantId: 1000,
			},
			wantRowsAffected: 1,
		},
		{
			name: "delete a restaurant that doesn't exist",
			args: args{
				restaurantId: 1001,
			},
			wantRowsAffected: 0,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			err := trManager.Do(ctx, func(ctx context.Context) error {
				ra, _ := restaurantRepository.Delete(ctx, tc.args.restaurantId)
				assert.Equal(t, tc.wantRowsAffected, ra)

				return errors.New(ROLLBACK_PLEASE)
			})

			assert.Error(t, err)
		})
	}
}

// createMockRepository creates a repository that uses gorm as usual but a sqlmock
// connection as the underlying database connection, returning also the sqlmock instance
// that we can use to define mockExpectations and a NOP transaction manager.
func createMockRepository() (*RestaurantPostgresRepository, trm.Manager, sqlmock.Sqlmock) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		DriverName:           "go-sqlmock",
		DSN:                  "go-sqlmock",
		PreferSimpleProtocol: true,
		WithoutReturning:     true,
		Conn:                 mockDb,
	})
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect to database")
	}
	return NewRestaurantPostgresRepository(db, trmgorm.DefaultCtxGetter, nil), test.NewNopTrManager(), mock
}

// --------------------------------------------------------------------------------
// Utility functions to create restaurants and menus.
// --------------------------------------------------------------------------------

func newTestRestaurant() *Restaurant {
	return &Restaurant{ID: 1000, Name: "restaurant1", Address: newTestAddress(), Menu: newTestMenu()}
}

func newTestRestaurantWithoutMenu() *Restaurant {
	return &Restaurant{ID: 1000, Name: "restaurant1", Address: newTestAddress(), Menu: nil}
}

func newTestRestaurantWithoutId() *Restaurant {
	return &Restaurant{Name: "restaurant1", Address: newTestAddress(), Menu: newTestMenu()}
}

func newTestAddress() *Address {
	return &Address{Street: "street1", City: "city1", State: "state1", Zip: "zip1"}
}

func newTestMenu() []*MenuItem {
	item1 := MenuItem{Id: 1, Name: "item1.1", Price: "13.14"}
	item2 := MenuItem{Id: 2, Name: "item1.2", Price: "14.15"}
	item3 := MenuItem{Id: 3, Name: "item1.3", Price: "15.16"}
	return []*MenuItem{&item1, &item2, &item3}
}
