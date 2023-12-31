package rest

import (
	"net/http"
	"strconv"

	"f4allgo-restaurant/internal/core/port"

	coreerrors "f4allgo-restaurant/internal/core/service/errors"

	"github.com/gin-gonic/gin"
)

type RestaurantHandler struct {
	mapper            Mapper
	restaurantService port.RestaurantService
}

// NewRestaurantHandler builds a new RestaurantHandler struct.
func NewRestaurantHandler(s port.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{mapper: DefaultMapper{}, restaurantService: s}
}

// CreateRestaurant creates a restaurant.
func (rh *RestaurantHandler) CreateRestaurant(ctx *gin.Context) {
	var request CreateRestaurantRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restaurantId, err := rh.restaurantService.Create(ctx, rh.mapper.toDomainRestaurant(request.Restaurant))
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
	ctx.JSON(http.StatusCreated, CreateRestaurantResponse{RestaurantId: restaurantId})
}

// DeleteRestaurant deletes a restaurant.
func (rh *RestaurantHandler) DeleteRestaurant(ctx *gin.Context) {
	restaurantId, err := strconv.ParseInt(ctx.Param("restaurantId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	err = rh.restaurantService.Delete(ctx, restaurantId)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

// GetRestaurants gets the list of all restaurants.
func (rh *RestaurantHandler) GetRestaurants(ctx *gin.Context) {
	offset, limit := getOffsetAndLimit(ctx)
	domainRestaurants, total, err := rh.restaurantService.FindAll(ctx, offset, limit)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	dtoRestaurants := rh.mapper.fromDomainRestaurants(domainRestaurants)
	ctx.JSON(http.StatusOK, GetRestaurantsResponse{Restaurants: dtoRestaurants, Total: total})
}

// GetRestaurant gets a restaurant by its ID.
func (rh *RestaurantHandler) GetRestaurant(ctx *gin.Context) {
	restaurantId, err := strconv.ParseInt(ctx.Param("restaurantId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	domainRestaurant, err := rh.restaurantService.FindById(ctx, restaurantId)
	if err != nil {
		handleError(ctx, err)
		return
	}

	dtoRestaurant := rh.mapper.fromDomainRestaurant(domainRestaurant)
	ctx.JSON(http.StatusOK, GetRestaurantResponse{Restaurant: dtoRestaurant})
}

// UpdateMenu updates the menu of a restaurant.
func (rh *RestaurantHandler) UpdateMenu(ctx *gin.Context) {
	restaurantId, err := strconv.ParseInt(ctx.Param("restaurantId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	var request UpdateMenuRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = rh.restaurantService.UpdateMenu(ctx, restaurantId, rh.mapper.toDomainMenu(request.Menu))
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}

// handleError translates core errors into the proper http error codes.
func handleError(ctx *gin.Context, err error) {
	switch e := err.(type) {

	case *coreerrors.RestaurantNotFoundError:
		ctx.AbortWithStatus(http.StatusNotFound)

	case *coreerrors.RepositoryError:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": e.Error()})

	case *coreerrors.EventPublisherError:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": e.Error()})

	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
	}
}

func getOffsetAndLimit(ctx *gin.Context) (int, int) {
	offsetStr := ctx.DefaultQuery("offset", "0")
	limitStr := ctx.DefaultQuery("limit", "10")

	offset, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		offset = 0
	}

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		limit = 10
	}

	// Set a hardcoded maximum limit of 100 elements per page.
	if limit > 100 {
		limit = 100
	}

	return int(offset), int(limit)
}
