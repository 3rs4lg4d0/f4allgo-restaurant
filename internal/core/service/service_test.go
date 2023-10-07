package service

import (
	"context"
	"errors"
	"f4allgo-restaurant/internal/core/domain"
	coreerrors "f4allgo-restaurant/internal/core/service/errors"
	"f4allgo-restaurant/internal/core/service/mocks"
	"f4allgo-restaurant/test"
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindAll(t *testing.T) {
	type args struct {
		ctx    context.Context
		offset int
		limit  int
	}
	testcases := []struct {
		name             string
		args             args
		mockExpectations func(args, *mocks.MockRestaurantRepository)
		wantRestaurants  []*domain.Restaurant
		wantTotal        int64
		wantErr          bool
		wantErrType      error
	}{
		{
			name: "mock a successful execution",
			args: args{
				ctx:    context.Background(),
				offset: 0,
				limit:  100,
			},
			mockExpectations: func(args args, repository *mocks.MockRestaurantRepository) {
				repository.EXPECT().FindAll(args.ctx, args.offset, args.limit).Return([]*domain.Restaurant{newTestRestaurant()}, int64(10), nil).Once()
			},
			wantRestaurants: []*domain.Restaurant{newTestRestaurant()},
			wantTotal:       10,
			wantErr:         false,
		},
		{
			name: "mock a failure execution",
			args: args{
				ctx:    context.Background(),
				offset: 0,
				limit:  100,
			},
			mockExpectations: func(args args, repository *mocks.MockRestaurantRepository) {
				repository.EXPECT().FindAll(args.ctx, args.offset, args.limit).Return(nil, 0, errors.New("error")).Once()
			},
			wantRestaurants: nil,
			wantTotal:       0,
			wantErr:         true,
			wantErrType:     &coreerrors.RepositoryError{},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepository := mocks.NewMockRestaurantRepository(t)
			tc.mockExpectations(tc.args, mockRepository)
			rs := NewDefaultRestaurantService(mockRepository, nil, nil)
			actualRestaurants, actualTotal, err := rs.FindAll(tc.args.ctx, 0, 100)
			if !tc.wantErr {
				assert.NoError(t, err)
				assert.True(t, reflect.DeepEqual(tc.wantRestaurants, actualRestaurants))
				assert.Equal(t, tc.wantTotal, actualTotal)
			} else {
				assert.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
			}
		})
	}
}

func TestFindById(t *testing.T) {
	type args struct {
		ctx          context.Context
		restaurantId int64
	}
	testcases := []struct {
		name             string
		args             args
		mockExpectations func(args, *mocks.MockRestaurantRepository)
		wantRestaurant   *domain.Restaurant
		wantErr          bool
		wantErrType      error
	}{
		{
			name: "mock a successful execution",
			args: args{
				ctx:          context.Background(),
				restaurantId: 1000,
			},
			mockExpectations: func(args args, repository *mocks.MockRestaurantRepository) {
				repository.EXPECT().FindById(args.ctx, args.restaurantId, true).Return(newTestRestaurant(), nil).Once()
			},
			wantRestaurant: newTestRestaurant(),
			wantErr:        false,
		},
		{
			name: "mock record not found",
			args: args{
				ctx: context.Background(),
			},
			mockExpectations: func(args args, repository *mocks.MockRestaurantRepository) {
				repository.EXPECT().FindById(args.ctx, args.restaurantId, true).Return(nil, errors.New("record not found")).Once()
			},
			wantRestaurant: nil,
			wantErr:        true,
			wantErrType:    &coreerrors.RestaurantNotFoundError{},
		},
		{
			name: "mock generic repository error",
			args: args{
				ctx: context.Background(),
			},
			mockExpectations: func(args args, repository *mocks.MockRestaurantRepository) {
				repository.EXPECT().FindById(args.ctx, args.restaurantId, true).Return(nil, errors.New("error#1")).Once()
			},
			wantRestaurant: nil,
			wantErr:        true,
			wantErrType:    &coreerrors.RepositoryError{},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepository := mocks.NewMockRestaurantRepository(t)
			tc.mockExpectations(tc.args, mockRepository)
			rs := NewDefaultRestaurantService(mockRepository, nil, nil)
			actualRestaurant, err := rs.FindById(tc.args.ctx, tc.args.restaurantId)
			if !tc.wantErr {
				assert.NoError(t, err)
				assert.True(t, reflect.DeepEqual(tc.wantRestaurant, actualRestaurant))
			} else {
				assert.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	type args struct {
		ctx        context.Context
		restaurant *domain.Restaurant
	}
	testcases := []struct {
		name                 string
		args                 args
		mockExpectations     func(args, *mocks.MockRestaurantRepository, *mocks.MockDomainEventPublisher)
		wantErr              bool
		wantErrType          error
		additionalAssertions func(args, *mocks.MockRestaurantRepository, *mocks.MockDomainEventPublisher)
	}{
		{
			name: "mock a successful execution",
			args: args{
				ctx:        context.Background(),
				restaurant: newTestRestaurant(),
			},
			mockExpectations: func(args args, mr *mocks.MockRestaurantRepository, mp *mocks.MockDomainEventPublisher) {
				mr.EXPECT().Save(args.ctx, args.restaurant).Return(nil).Once()
				mp.EXPECT().Publish(args.ctx, domain.NewRestaurantCreated(args.restaurant)).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name: "mock a DomainEventPublisher failure",
			args: args{
				ctx:        context.Background(),
				restaurant: newTestRestaurant(),
			},
			mockExpectations: func(args args, mr *mocks.MockRestaurantRepository, mp *mocks.MockDomainEventPublisher) {
				mr.EXPECT().Save(args.ctx, args.restaurant).Return(nil).Once()
				mp.EXPECT().Publish(args.ctx, domain.NewRestaurantCreated(args.restaurant)).Return(errors.New("error")).Once()
			},
			wantErr:     true,
			wantErrType: &coreerrors.EventPublisherError{},
		},
		{
			name: "mock a RestaurantRepository failure",
			args: args{
				ctx:        context.Background(),
				restaurant: &domain.Restaurant{},
			},
			mockExpectations: func(args args, mr *mocks.MockRestaurantRepository, mp *mocks.MockDomainEventPublisher) {
				mr.EXPECT().Save(args.ctx, args.restaurant).Return(errors.New("error")).Once()
			},
			wantErr:     true,
			wantErrType: &coreerrors.RepositoryError{},
			additionalAssertions: func(args args, _ *mocks.MockRestaurantRepository, mp *mocks.MockDomainEventPublisher) {
				mp.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mr := mocks.NewMockRestaurantRepository(t)
			mp := mocks.NewMockDomainEventPublisher(t)
			tc.mockExpectations(tc.args, mr, mp)
			rs := NewDefaultRestaurantService(mr, mp, test.NewNopTrManager())
			err := rs.Create(tc.args.ctx, tc.args.restaurant)
			if !tc.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
			}
			if tc.additionalAssertions != nil {
				tc.additionalAssertions(tc.args, mr, mp)
			}
		})
	}
}

func TestUpdateMenu(t *testing.T) {
	type args struct {
		ctx          context.Context
		restaurantId int64
		menu         *domain.Menu
	}
	testcases := []struct {
		name                 string
		args                 args
		mockExpectations     func(args, *mocks.MockRestaurantRepository, *mocks.MockDomainEventPublisher)
		wantErr              bool
		wantErrType          error
		additionalAssertions func(args, *mocks.MockRestaurantRepository, *mocks.MockDomainEventPublisher)
	}{
		{
			name: "mock a successful execution",
			args: args{
				ctx:          context.Background(),
				restaurantId: 1000,
				menu:         domain.NewMenu([]*domain.MenuItem{domain.NewMenuItem(1, "one", big.NewFloat(1.0))}),
			},
			mockExpectations: func(args args, mr *mocks.MockRestaurantRepository, mp *mocks.MockDomainEventPublisher) {
				r := newTestRestaurant()
				mr.EXPECT().FindById(args.ctx, args.restaurantId, false).Return(r, nil).Once()
				rr := &domain.Restaurant{Id: r.Id, Name: r.Name, Address: r.Address, Menu: args.menu}
				mr.EXPECT().Update(args.ctx, rr).Return(1, nil).Once()
				mp.EXPECT().Publish(args.ctx, domain.NewRestaurantMenuUpdated(args.restaurantId, args.menu)).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name: "mock a DomainEventPublisher failure",
			args: args{
				ctx:          context.Background(),
				restaurantId: 1000,
				menu:         domain.NewMenu([]*domain.MenuItem{domain.NewMenuItem(1, "one", big.NewFloat(1.0))}),
			},
			mockExpectations: func(args args, mr *mocks.MockRestaurantRepository, mp *mocks.MockDomainEventPublisher) {
				r := newTestRestaurant()
				mr.EXPECT().FindById(args.ctx, args.restaurantId, false).Return(r, nil).Once()
				rr := &domain.Restaurant{Id: r.Id, Name: r.Name, Address: r.Address, Menu: args.menu}
				mr.EXPECT().Update(args.ctx, rr).Return(1, nil).Once()
				mp.EXPECT().Publish(args.ctx, domain.NewRestaurantMenuUpdated(args.restaurantId, args.menu)).Return(errors.New("error")).Once()
			},
			wantErr:     true,
			wantErrType: &coreerrors.EventPublisherError{},
		},
		{
			name: "mock a RestaurantRepository failure when updating",
			args: args{
				ctx:          context.Background(),
				restaurantId: 1000,
				menu:         domain.NewMenu([]*domain.MenuItem{domain.NewMenuItem(1, "one", big.NewFloat(1.0))}),
			},
			mockExpectations: func(args args, mr *mocks.MockRestaurantRepository, mp *mocks.MockDomainEventPublisher) {
				r := newTestRestaurant()
				mr.EXPECT().FindById(args.ctx, args.restaurantId, false).Return(r, nil).Once()
				rr := &domain.Restaurant{Id: r.Id, Name: r.Name, Address: r.Address, Menu: args.menu}
				mr.EXPECT().Update(args.ctx, rr).Return(1, errors.New("error")).Once()
			},
			wantErr:     true,
			wantErrType: &coreerrors.RepositoryError{},
			additionalAssertions: func(args args, mr *mocks.MockRestaurantRepository, mp *mocks.MockDomainEventPublisher) {
				mp.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything)
			},
		},
		{
			name: "provide an invalid menu not compliant with invariants",
			args: args{
				ctx:          context.Background(),
				restaurantId: 1000,
				menu:         &domain.Menu{},
			},
			mockExpectations: func(args args, mr *mocks.MockRestaurantRepository, _ *mocks.MockDomainEventPublisher) {
				restaurant := newTestRestaurant()
				mr.EXPECT().FindById(args.ctx, args.restaurantId, false).Return(restaurant, nil).Once()
			},
			wantErr:     true,
			wantErrType: &coreerrors.CoreError{},
			additionalAssertions: func(args args, mr *mocks.MockRestaurantRepository, mp *mocks.MockDomainEventPublisher) {
				restaurant := newTestRestaurant()
				restaurant.Menu = args.menu
				mr.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
				mp.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything)
			},
		},
		{
			name: "mock a RestaurantRepository failure when fetching",
			args: args{
				ctx:          context.Background(),
				restaurantId: 1000,
				menu:         domain.NewMenu([]*domain.MenuItem{domain.NewMenuItem(1, "one", big.NewFloat(1.0))}),
			},
			mockExpectations: func(args args, mr *mocks.MockRestaurantRepository, _ *mocks.MockDomainEventPublisher) {
				mr.EXPECT().FindById(args.ctx, args.restaurantId, false).Return(nil, errors.New("error")).Once()
			},
			wantErr:     true,
			wantErrType: &coreerrors.RepositoryError{},
			additionalAssertions: func(args args, mr *mocks.MockRestaurantRepository, mp *mocks.MockDomainEventPublisher) {
				mr.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
				mp.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mr := mocks.NewMockRestaurantRepository(t)
			mp := mocks.NewMockDomainEventPublisher(t)
			tc.mockExpectations(tc.args, mr, mp)
			rs := &DefaultRestaurantService{
				restaurantRepository: mr,
				domainEventPublisher: mp,
				trManager:            test.NewNopTrManager(),
			}
			err := rs.UpdateMenu(tc.args.ctx, tc.args.restaurantId, tc.args.menu)
			if !tc.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
			}
			if tc.additionalAssertions != nil {
				tc.additionalAssertions(tc.args, mr, mp)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		ctx          context.Context
		restaurantId int64
	}
	testcases := []struct {
		name                 string
		args                 args
		mockExpectations     func(args, *mocks.MockRestaurantRepository, *mocks.MockDomainEventPublisher)
		wantErr              bool
		wantErrType          error
		additionalAssertions func(args, *mocks.MockRestaurantRepository, *mocks.MockDomainEventPublisher)
	}{
		{
			name: "mock a successful execution",
			args: args{
				ctx:          context.Background(),
				restaurantId: 1,
			},
			mockExpectations: func(args args, mr *mocks.MockRestaurantRepository, mp *mocks.MockDomainEventPublisher) {
				mr.EXPECT().Delete(args.ctx, args.restaurantId).Return(1, nil).Once()
				mp.EXPECT().Publish(args.ctx, domain.NewRestaurantDeleted(args.restaurantId)).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name: "mock a DomainEventPublisher failure",
			args: args{
				ctx:          context.Background(),
				restaurantId: 1,
			},
			mockExpectations: func(args args, mr *mocks.MockRestaurantRepository, mp *mocks.MockDomainEventPublisher) {
				mr.EXPECT().Delete(args.ctx, args.restaurantId).Return(1, nil).Once()
				mp.EXPECT().Publish(args.ctx, domain.NewRestaurantDeleted(args.restaurantId)).Return(errors.New("error")).Once()
			},
			wantErr:     true,
			wantErrType: &coreerrors.EventPublisherError{},
		},
		{
			name: "mock a RestaurantRepository failure",
			args: args{
				ctx:          context.Background(),
				restaurantId: 1,
			},
			mockExpectations: func(args args, mr *mocks.MockRestaurantRepository, mp *mocks.MockDomainEventPublisher) {
				mr.EXPECT().Delete(args.ctx, args.restaurantId).Return(1, errors.New("error")).Once()
			},
			wantErr:     true,
			wantErrType: &coreerrors.RepositoryError{},
			additionalAssertions: func(args args, mr *mocks.MockRestaurantRepository, mp *mocks.MockDomainEventPublisher) {
				mp.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything)
			},
		},
		{
			name: "mock a RestaurantNotFound failure",
			args: args{
				ctx:          context.Background(),
				restaurantId: 1,
			},
			mockExpectations: func(args args, mr *mocks.MockRestaurantRepository, mp *mocks.MockDomainEventPublisher) {
				mr.EXPECT().Delete(args.ctx, args.restaurantId).Return(0, nil).Once()
			},
			wantErr:     true,
			wantErrType: &coreerrors.RestaurantNotFoundError{},
			additionalAssertions: func(args args, mr *mocks.MockRestaurantRepository, mp *mocks.MockDomainEventPublisher) {
				mp.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mr := mocks.NewMockRestaurantRepository(t)
			mp := mocks.NewMockDomainEventPublisher(t)
			tc.mockExpectations(tc.args, mr, mp)
			rs := &DefaultRestaurantService{
				restaurantRepository: mr,
				domainEventPublisher: mp,
				trManager:            test.NewNopTrManager(),
			}
			err := rs.Delete(tc.args.ctx, tc.args.restaurantId)
			if !tc.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
			}
			if tc.additionalAssertions != nil {
				tc.additionalAssertions(tc.args, mr, mp)
			}
		})
	}
}

// --------------------------------------------------------------------------------
// Utility functions to create restaurants and menus.
// --------------------------------------------------------------------------------

func newTestRestaurant() *domain.Restaurant {
	return &domain.Restaurant{Id: 1000, Name: "restaurant1", Address: newTestAddress(), Menu: newTestMenu()}
}

func newTestAddress() *domain.Address {
	return domain.NewAddress("street1", "city1", "state1", "zip1")
}

func newTestMenu() *domain.Menu {
	f := new(big.Float)
	f.SetString("13.14")
	item1 := domain.NewMenuItem(1, "item1.1", f)
	f = new(big.Float)
	f.SetString("14.15")
	item2 := domain.NewMenuItem(1, "item1.2", f)
	f = new(big.Float)
	f.SetString("15.16")
	item3 := domain.NewMenuItem(1, "item1.3", f)
	return domain.NewMenu([]*domain.MenuItem{item1, item2, item3})
}
