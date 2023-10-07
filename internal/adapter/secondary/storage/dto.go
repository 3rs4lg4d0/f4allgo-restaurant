package storage

// Restaurant is a Gorm DTO that carries the information of domain restaurants.
type Restaurant struct {
	ID      int64
	Name    string
	Address *Address    `gorm:"embedded"`
	Menu    []*MenuItem `gorm:"foreignKey:RestaurantID"`
}

func (Restaurant) TableName() string {
	return "restaurant"
}

// Address is a Gorm DTO that carries the information of domain addresses.
type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

// MenuItem is a Gorm DTO that carries the information of domain menu items
type MenuItem struct {
	RestaurantID int64 `gorm:"primaryKey"`
	Id           int32 `gorm:"primaryKey"`
	Name         string
	Price        string
}

func (MenuItem) TableName() string {
	return "menu_item"
}
