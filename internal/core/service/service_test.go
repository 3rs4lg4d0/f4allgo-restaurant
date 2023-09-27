package service

import (
	"context"
	"errors"
	"f4allgo-restaurant/internal/core/domain"
	coreerrors "f4allgo-restaurant/internal/core/service/errors"
	"f4allgo-restaurant/internal/core/service/mocks"
	"math/big"
	"reflect"
	"testing"

	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/stretchr/testify/assert"
)

func TestDefaultRestaurantService_FindAll(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name         string
		expectations func(repository *mocks.MockRestaurantRepository)
		args         args
		want         []*domain.Restaurant
		wantErr      bool
	}{
		{
			name: "When RestaurantRepository returns an empty slice",
			expectations: func(repository *mocks.MockRestaurantRepository) {
				repository.EXPECT().FindAll(context.Background(), uint32(0), uint8(100)).Return(make([]*domain.Restaurant, 0), nil).Once()
			},
			args: args{
				ctx: context.Background(),
			},
			want:    make([]*domain.Restaurant, 0),
			wantErr: false,
		},
		{
			name: "When RestaurantRepository returns an error",
			expectations: func(repository *mocks.MockRestaurantRepository) {
				repository.EXPECT().FindAll(context.Background(), uint32(0), uint8(100)).Return(nil, errors.New("error")).Once()
			},
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository := mocks.NewMockRestaurantRepository(t)
			tt.expectations(mockRepository)

			rs := NewDefaultRestaurantService(mockRepository, nil, nil)
			got, err := rs.FindAll(tt.args.ctx, 0, 100)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultRestaurantService.FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultRestaurantService.FindAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultRestaurantService_Create(t *testing.T) {
	type mockFields struct {
		restaurantRepository *mocks.MockRestaurantRepository
		domainEventPublisher *mocks.MockDomainEventPublisher
	}
	type args struct {
		ctx        context.Context
		restaurant *domain.Restaurant
	}
	tests := []struct {
		name         string
		expectations func(args args, mocks mockFields)
		args         args
		wantErr      bool
		assertions   func(mocks mockFields)
	}{
		{
			name: "When RestaurantRepository and DomainEventPublisher succeed",
			expectations: func(args args, mocks mockFields) {
				mocks.restaurantRepository.EXPECT().Save(args.ctx, args.restaurant).Return(nil).Once()
				mocks.domainEventPublisher.EXPECT().Publish(args.ctx, domain.NewRestaurantCreated(args.restaurant)).Return(nil).Once()
			},
			args: args{
				ctx:        context.Background(),
				restaurant: &domain.Restaurant{},
			},
			wantErr: false,
		},
		{
			name: "When DomainEventPublisher fails",
			expectations: func(args args, mocks mockFields) {
				mocks.restaurantRepository.EXPECT().Save(args.ctx, args.restaurant).Return(nil).Once()
				mocks.domainEventPublisher.EXPECT().Publish(args.ctx, domain.NewRestaurantCreated(args.restaurant)).Return(errors.New("error")).Once()
			},
			args: args{
				ctx:        context.Background(),
				restaurant: &domain.Restaurant{},
			},
			wantErr: true,
		},
		{
			name: "When RestaurantRepository fails",
			expectations: func(args args, mocks mockFields) {
				mocks.restaurantRepository.EXPECT().Save(args.ctx, args.restaurant).Return(errors.New("error")).Once()
			},
			args: args{
				ctx:        context.Background(),
				restaurant: &domain.Restaurant{},
			},
			wantErr: true,
			assertions: func(mocks mockFields) {
				mocks.domainEventPublisher.AssertNotCalled(t, "Publish")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mockFields{
				restaurantRepository: mocks.NewMockRestaurantRepository(t),
				domainEventPublisher: mocks.NewMockDomainEventPublisher(t),
			}
			tt.expectations(tt.args, m)
			rs := &DefaultRestaurantService{
				restaurantRepository: m.restaurantRepository,
				domainEventPublisher: m.domainEventPublisher,
				trManager:            &mockTrManager{},
			}
			if err := rs.Create(tt.args.ctx, tt.args.restaurant); (err != nil) != tt.wantErr {
				t.Errorf("DefaultRestaurantService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assertions != nil {
				tt.assertions(m)
			}
		})
	}
}

func TestDefaultRestaurantService_UpdateMenu(t *testing.T) {
	type mockFields struct {
		restaurantRepository *mocks.MockRestaurantRepository
		domainEventPublisher *mocks.MockDomainEventPublisher
	}
	type args struct {
		ctx          context.Context
		restaurantId uint64
		menu         *domain.Menu
	}
	tests := []struct {
		name         string
		expectations func(args args, mocks mockFields)
		args         args
		wantErr      bool
		assertions   func(mocks mockFields)
	}{
		{
			name: "When RestaurantRepository and DomainEventPublisher succeed",
			expectations: func(args args, mocks mockFields) {
				restaurant := &domain.Restaurant{}
				mocks.restaurantRepository.EXPECT().FindById(args.ctx, args.restaurantId).Return(restaurant, nil).Once()
				mocks.restaurantRepository.EXPECT().Update(args.ctx, restaurant).Return(nil).Once()
				mocks.domainEventPublisher.EXPECT().Publish(args.ctx, domain.NewRestaurantMenuUpdated(args.restaurantId, args.menu)).Return(nil).Once()
			},
			args: args{
				ctx:          context.Background(),
				restaurantId: 1,
				menu:         domain.NewMenu([]*domain.MenuItem{domain.NewMenuItem(1, "one", big.NewFloat(1.0))}),
			},
			wantErr: false,
		},
		{
			name: "When DomainEventPublisher fails",
			expectations: func(args args, mocks mockFields) {
				restaurant := &domain.Restaurant{}
				mocks.restaurantRepository.EXPECT().FindById(args.ctx, args.restaurantId).Return(restaurant, nil).Once()
				mocks.restaurantRepository.EXPECT().Update(args.ctx, restaurant).Return(nil).Once()
				mocks.domainEventPublisher.EXPECT().Publish(args.ctx, domain.NewRestaurantMenuUpdated(args.restaurantId, args.menu)).Return(errors.New("error")).Once()
			},
			args: args{
				ctx:          context.Background(),
				restaurantId: 1,
				menu:         domain.NewMenu([]*domain.MenuItem{domain.NewMenuItem(1, "one", big.NewFloat(1.0))}),
			},
			wantErr: true,
		},
		{
			name: "When RestaurantRepository fails updating the restaurant",
			expectations: func(args args, mocks mockFields) {
				restaurant := &domain.Restaurant{}
				mocks.restaurantRepository.EXPECT().FindById(args.ctx, args.restaurantId).Return(restaurant, nil).Once()
				mocks.restaurantRepository.EXPECT().Update(args.ctx, restaurant).Return(errors.New("error")).Once()
			},
			args: args{
				ctx:          context.Background(),
				restaurantId: 1,
				menu:         domain.NewMenu([]*domain.MenuItem{domain.NewMenuItem(1, "one", big.NewFloat(1.0))}),
			},
			wantErr: true,
			assertions: func(mocks mockFields) {
				mocks.domainEventPublisher.AssertNotCalled(t, "Publish")
			},
		},
		{
			name: "When invalid menu",
			expectations: func(args args, mocks mockFields) {
				restaurant := &domain.Restaurant{}
				mocks.restaurantRepository.EXPECT().FindById(args.ctx, args.restaurantId).Return(restaurant, nil).Once()
			},
			args: args{
				ctx:          context.Background(),
				restaurantId: 1,
				menu:         &domain.Menu{},
			},
			wantErr: true,
			assertions: func(mocks mockFields) {
				mocks.restaurantRepository.AssertNotCalled(t, "Update")
				mocks.domainEventPublisher.AssertNotCalled(t, "Publish")
			},
		},
		{
			name: "When RestaurantRepository fails finding the restaurant",
			expectations: func(args args, mocks mockFields) {
				mocks.restaurantRepository.EXPECT().FindById(args.ctx, args.restaurantId).Return(nil, errors.New("error")).Once()
			},
			args: args{
				ctx:          context.Background(),
				restaurantId: 1,
				menu:         domain.NewMenu([]*domain.MenuItem{domain.NewMenuItem(1, "one", big.NewFloat(1.0))}),
			},
			wantErr: true,
			assertions: func(mocks mockFields) {
				mocks.restaurantRepository.AssertNotCalled(t, "Update")
				mocks.domainEventPublisher.AssertNotCalled(t, "Publish")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mockFields{
				restaurantRepository: mocks.NewMockRestaurantRepository(t),
				domainEventPublisher: mocks.NewMockDomainEventPublisher(t),
			}
			tt.expectations(tt.args, m)
			rs := &DefaultRestaurantService{
				restaurantRepository: m.restaurantRepository,
				domainEventPublisher: m.domainEventPublisher,
				trManager:            &mockTrManager{},
			}
			if err := rs.UpdateMenu(tt.args.ctx, tt.args.restaurantId, tt.args.menu); (err != nil) != tt.wantErr {
				t.Errorf("DefaultRestaurantService.UpdateMenu() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assertions != nil {
				tt.assertions(m)
			}
		})
	}
}

func TestDefaultRestaurantService_DeleteMenu(t *testing.T) {
	type mockFields struct {
		restaurantRepository *mocks.MockRestaurantRepository
		domainEventPublisher *mocks.MockDomainEventPublisher
	}
	type args struct {
		ctx          context.Context
		restaurantId uint64
	}
	tests := []struct {
		name         string
		expectations func(args args, mocks mockFields)
		args         args
		wantErr      bool
		assertions   func(mocks mockFields)
	}{
		{
			name: "When RestaurantRepository and DomainEventPublisher succeed",
			expectations: func(args args, mocks mockFields) {
				restaurant := &domain.Restaurant{}
				mocks.restaurantRepository.EXPECT().FindById(args.ctx, args.restaurantId).Return(restaurant, nil).Once()
				mocks.restaurantRepository.EXPECT().Delete(args.ctx, args.restaurantId).Return(nil).Once()
				mocks.domainEventPublisher.EXPECT().Publish(args.ctx, domain.NewRestaurantDeleted(args.restaurantId)).Return(nil).Once()
			},
			args: args{
				ctx:          context.Background(),
				restaurantId: 1,
			},
			wantErr: false,
		},
		{
			name: "When DomainEventPublisher fails",
			expectations: func(args args, mocks mockFields) {
				restaurant := &domain.Restaurant{}
				mocks.restaurantRepository.EXPECT().FindById(args.ctx, args.restaurantId).Return(restaurant, nil).Once()
				mocks.restaurantRepository.EXPECT().Delete(args.ctx, args.restaurantId).Return(nil).Once()
				mocks.domainEventPublisher.EXPECT().Publish(args.ctx, domain.NewRestaurantDeleted(args.restaurantId)).Return(errors.New("error")).Once()
			},
			args: args{
				ctx:          context.Background(),
				restaurantId: 1,
			},
			wantErr: true,
		},
		{
			name: "When RestaurantRepository fails deleting",
			expectations: func(args args, mocks mockFields) {
				restaurant := &domain.Restaurant{}
				mocks.restaurantRepository.EXPECT().FindById(args.ctx, args.restaurantId).Return(restaurant, nil).Once()
				mocks.restaurantRepository.EXPECT().Delete(args.ctx, args.restaurantId).Return(errors.New("error")).Once()
			},
			args: args{
				ctx:          context.Background(),
				restaurantId: 1,
			},
			wantErr: true,
			assertions: func(mocks mockFields) {
				mocks.domainEventPublisher.AssertNotCalled(t, "Publish")
			},
		},
		{
			name: "When RestaurantRepository fails finding the restaurant",
			expectations: func(args args, mocks mockFields) {
				mocks.restaurantRepository.EXPECT().FindById(args.ctx, args.restaurantId).Return(nil, errors.New("error")).Once()
			},
			args: args{
				ctx:          context.Background(),
				restaurantId: 1,
			},
			wantErr: true,
			assertions: func(mocks mockFields) {
				mocks.restaurantRepository.AssertNotCalled(t, "Delete")
				mocks.domainEventPublisher.AssertNotCalled(t, "Publish")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mockFields{
				restaurantRepository: mocks.NewMockRestaurantRepository(t),
				domainEventPublisher: mocks.NewMockDomainEventPublisher(t),
			}
			tt.expectations(tt.args, m)
			rs := &DefaultRestaurantService{
				restaurantRepository: m.restaurantRepository,
				domainEventPublisher: m.domainEventPublisher,
				trManager:            &mockTrManager{},
			}
			if err := rs.Delete(tt.args.ctx, tt.args.restaurantId); (err != nil) != tt.wantErr {
				t.Errorf("DefaultRestaurantService.UpdateMenu() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assertions != nil {
				tt.assertions(m)
			}
		})
	}
}

// mockTrManager mocks the transaction manager. It just calls the closure passed as
// argument and return the result.
type mockTrManager struct{}

func (mtr *mockTrManager) Do(ctx context.Context, f func(ctx context.Context) error) error {
	return f(ctx)
}

func (mtr *mockTrManager) DoWithSettings(ctx context.Context, _ trm.Settings, f func(ctx context.Context) error) error {
	return f(ctx)
}

func TestDefaultRestaurantService_findById(t *testing.T) {
	type mockFields struct {
		restaurantRepository *mocks.MockRestaurantRepository
	}
	type args struct {
		ctx          context.Context
		restaurantId uint64
	}
	tests := []struct {
		name             string
		expectations     func(args args, mocks mockFields)
		args             args
		expectedError    error
		expectedErrorMsg string
	}{
		{
			name: "When RestaurantRepository fails finding the restaurant",
			expectations: func(args args, mocks mockFields) {
				mocks.restaurantRepository.EXPECT().FindById(args.ctx, args.restaurantId).Return(nil, errors.New("customError")).Once()
			},
			args: args{
				ctx:          context.Background(),
				restaurantId: 1,
			},
			expectedError:    &coreerrors.RepositoryError{},
			expectedErrorMsg: "customError",
		},
		{
			name: "When RestaurantRepository returns no restaurant",
			expectations: func(args args, mocks mockFields) {
				mocks.restaurantRepository.EXPECT().FindById(args.ctx, args.restaurantId).Return(nil, nil).Once()
			},
			args: args{
				ctx:          context.Background(),
				restaurantId: 1,
			},
			expectedError:    &coreerrors.RestaurantNotFoundError{},
			expectedErrorMsg: "restaurant not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mockFields{
				restaurantRepository: mocks.NewMockRestaurantRepository(t),
			}
			tt.expectations(tt.args, m)
			rs := &DefaultRestaurantService{
				restaurantRepository: m.restaurantRepository,
			}
			_, err := rs.findById(tt.args.ctx, tt.args.restaurantId)
			assert.Equal(t, true, reflect.TypeOf(err) == reflect.TypeOf(tt.expectedError))
			assert.Equal(t, tt.expectedErrorMsg, err.Error())
		})
	}
}
