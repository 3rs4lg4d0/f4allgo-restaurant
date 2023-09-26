package cli

// --------------------------------------------------------------------------------
// OpenAPI components :: model
// --------------------------------------------------------------------------------

type Restaurant struct {
	Id      int64    `json:"id"`
	Name    string   `json:"name" binding:"required,max=255"`
	Address *Address `json:"address" binding:"required"`
	Menu    *Menu    `json:"menu" binding:"required"`
}

type Address struct {
	Street string `json:"street" binding:"required,max=255"`
	City   string `json:"city" binding:"required,max=255"`
	State  string `json:"state" binding:"required,max=255"`
	Zip    string `json:"zip" binding:"required,max=255"`
}

type Menu struct {
	Items []MenuItem `json:"items" binding:"required,min=2,max=1000,dive"`
}

type MenuItem struct {
	Id    int32  `json:"id" binding:"required"`
	Name  string `json:"name" binding:"required,max=255"`
	Price string `json:"price" binding:"required,max=10"`
}

// --------------------------------------------------------------------------------
// OpenAPI components :: requests/responses
// --------------------------------------------------------------------------------

type CreateRestaurantRequest struct {
	Restaurant *Restaurant `json:"restaurant" binding:"required"`
}

type UpdateMenuRequest struct {
	Menu *Menu `json:"menu" binding:"required"`
}

type GetRestaurantsResponse struct {
	Restaurants []*Restaurant `json:"restaurants"`
}
