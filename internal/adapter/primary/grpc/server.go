package grpc

import (
	context "context"
	"f4allgo-restaurant/internal/core/port"
)

// Implementation of the generated interface RestaurantServiceServer.
type restaurantServiceServer struct {
	UnimplementedRestaurantServiceServer
	mapper            Mapper
	restaurantService port.RestaurantService
}

// Interface compliance verification.
var _ RestaurantServiceServer = (*restaurantServiceServer)(nil)

func NewRestaurantServiceServer(s port.RestaurantService) *restaurantServiceServer {
	return &restaurantServiceServer{mapper: DefaultMapper{}, restaurantService: s}
}

func (rs *restaurantServiceServer) GetRestaurants(ctx context.Context, req *GetRestaurantsRequest) (*GetRestaurantsResponse, error) {
	offset, limit := getOffsetAndLimit(req)
	domainRestaurants, total, err := rs.restaurantService.FindAll(ctx, int(offset), int(limit))
	if err != nil {
		return nil, err
	}
	dtoRestaurants := rs.mapper.fromDomainRestaurants(domainRestaurants)
	return &GetRestaurantsResponse{Restaurants: dtoRestaurants, Total: int32(total)}, nil
}

func (rs *restaurantServiceServer) CreateRestaurant(ctx context.Context, req *CreateRestaurantRequest) (*CreateRestaurantResponse, error) {
	restaurantId, err := rs.restaurantService.Create(ctx, rs.mapper.toDomainRestaurant(req.Restaurant))

	return &CreateRestaurantResponse{RestaurantId: restaurantId}, err
}

func (rs *restaurantServiceServer) UpdateMenu(ctx context.Context, req *UpdateMenuRequest) (*UpdateMenuResponse, error) {
	err := rs.restaurantService.UpdateMenu(ctx, req.RestaurantId, rs.mapper.toDomainMenu(req.Menu))
	return &UpdateMenuResponse{}, err
}

func (rs *restaurantServiceServer) GetRestaurantById(ctx context.Context, req *GetRestaurantByIdRequest) (*GetRestaurantResponse, error) {
	domainRestaurant, err := rs.restaurantService.FindById(ctx, req.RestaurantId)
	if err != nil {
		return nil, err
	}

	dtoRestaurant := rs.mapper.fromDomainRestaurant(domainRestaurant)
	return &GetRestaurantResponse{Restaurant: dtoRestaurant}, nil
}

func (rs *restaurantServiceServer) DeleteRestaurant(ctx context.Context, req *DeleteRestaurantRequest) (*DeleteRestaurantResponse, error) {
	err := rs.restaurantService.Delete(ctx, req.RestaurantId)
	return &DeleteRestaurantResponse{}, err
}

func getOffsetAndLimit(req *GetRestaurantsRequest) (int, int) {
	offset := req.Offset
	limit := req.Limit

	if limit == 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	return int(offset), int(limit)
}
