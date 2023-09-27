package cli

import (
	"context"
	"encoding/json"
	"f4allgo-restaurant/internal/core/port"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

type RestaurantCli struct {
	mapper            Mapper
	restaurantService port.RestaurantService
	ctx               context.Context
	validate          *validator.Validate
}

// NewRestaurantHandler builds a new RestaurantHandler struct.
func NewRestaurantCli(s port.RestaurantService) *RestaurantCli {
	validator := validator.New()
	validator.SetTagName("binding")
	return &RestaurantCli{mapper: DefaultMapper{}, restaurantService: s, ctx: context.Background(), validate: validator}
}

func (rc *RestaurantCli) Execute() error {
	var rootCmd = &cobra.Command{Use: "f4allgorestaurant-cli"}

	// Level 1 subcomands.
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get resources",
		Long:  "Get various types of resources",
	}
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create resources",
		Long:  "Create various types of resources",
	}
	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update resources",
		Long:  "Update various types of resources",
	}
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete resources",
		Long:  "Delete various types of resources",
	}

	// Level 2 subcommands.
	var getRestaurantsCmd = &cobra.Command{
		Use:   "restaurants",
		Short: "Get restaurants",
		RunE: func(cmd *cobra.Command, args []string) error {
			offset, limit := getOffsetAndLimit(cmd)
			return rc.getRestaurants(offset, limit)
		},
	}
	var createRestaurantCmd = &cobra.Command{
		Use:   "restaurant",
		Short: "Create restaurant",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonData, err := cmd.Flags().GetString("json")
			if err != nil {
				return err
			}
			var request CreateRestaurantRequest
			if err := json.Unmarshal([]byte(jsonData), &request); err != nil {
				return err
			}
			return rc.createRestaurant(request)
		},
	}
	createRestaurantCmd.PersistentFlags().String("json", "", "the JSON payload")

	var updateRestaurantMenuCmd = &cobra.Command{
		Use:   "menu",
		Short: "Update restaurant menu",
		RunE: func(cmd *cobra.Command, args []string) error {
			restaurantId, err := cmd.Flags().GetUint64("restaurantId")
			if err != nil {
				return err
			}
			jsonData, err := cmd.Flags().GetString("json")
			if err != nil {
				return err
			}
			var request UpdateMenuRequest
			if err := json.Unmarshal([]byte(jsonData), &request); err != nil {
				return err
			}
			return rc.updateMenu(restaurantId, request)
		},
	}
	updateRestaurantMenuCmd.PersistentFlags().Uint64("restaurantId", 0, "the restaurant id")
	updateRestaurantMenuCmd.PersistentFlags().String("json", "", "the JSON payload")

	var deleteRestaurantCmd = &cobra.Command{
		Use:   "restaurant",
		Short: "Delete restaurant",
		RunE: func(cmd *cobra.Command, args []string) error {
			restaurantId, err := cmd.Flags().GetUint64("restaurantId")
			if err != nil {
				return err
			}
			return rc.deleteRestaurant(restaurantId)
		},
	}
	deleteRestaurantCmd.PersistentFlags().Uint64("restaurantId", 0, "the restaurant id")

	// Subcommands for 'get'.
	getCmd.AddCommand(getRestaurantsCmd)

	// Subcommands for 'create'.
	createCmd.AddCommand(createRestaurantCmd)

	// Subcommands for 'update'.
	updateCmd.AddCommand(updateRestaurantMenuCmd)

	// Subcommands for 'delete'.
	deleteCmd.AddCommand(deleteRestaurantCmd)

	// Register all the subcommands.
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(deleteCmd)

	return rootCmd.Execute()
}

// CreateRestaurant creates a restaurant.
func (rc *RestaurantCli) createRestaurant(request CreateRestaurantRequest) error {
	if err := rc.validate.Struct(request); err != nil {
		return err
	}
	return rc.restaurantService.Create(rc.ctx, rc.mapper.toDomainRestaurant(request.Restaurant))
}

// DeleteRestaurant deletes a restaurant.
func (rc *RestaurantCli) deleteRestaurant(restaurantId uint64) error {
	return rc.restaurantService.Delete(rc.ctx, restaurantId)
}

// GetRestaurants gets the list of all restaurants.
func (rc *RestaurantCli) getRestaurants(offset uint32, limit uint8) error {
	domainRestaurants, err := rc.restaurantService.FindAll(rc.ctx, offset, limit)
	if err != nil {
		return err
	}
	dtoRestaurants := rc.mapper.fromDomainRestaurants(domainRestaurants)
	return printJSON(GetRestaurantsResponse{Restaurants: dtoRestaurants})
}

// UpdateMenu updates the menu of a restaurant.
func (rc *RestaurantCli) updateMenu(restaurantId uint64, request UpdateMenuRequest) error {
	if err := rc.validate.Struct(request); err != nil {
		return err
	}
	return rc.restaurantService.UpdateMenu(rc.ctx, restaurantId, rc.mapper.toDomainMenu(request.Menu))
}

func printJSON(jsonStruct any) error {
	jsonData, err := json.MarshalIndent(jsonStruct, "", "  ")
	if err != nil {
		return err
	}
	jsonStr := string(jsonData)
	fmt.Println(jsonStr)
	fmt.Println()
	return nil
}

func getOffsetAndLimit(cmd *cobra.Command) (int32, int8) {
	offsetStr, _ := cmd.Flags().GetString("offset")
	limitStr, _ := cmd.Flags().GetString("limit")

	offset, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		offset = 0
	}

	limit, err := strconv.ParseInt(limitStr, 10, 8)
	if err != nil {
		limit = 10
	}

	// Set a hardcoded maximum limit of 100 elements per page.
	if limit > 100 {
		limit = 100
	}

	return int32(offset), int8(limit)
}
