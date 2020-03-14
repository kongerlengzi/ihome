package model

import "github.com/jinzhu/gorm"

type Houses struct {
	gorm.Model
	UserID uint            `json:"user_id" gorm:"not null"`
	AreaID uint            `json:"area_id" gorm:"not null"`
	Title string           `json:"title" gorm:"not null"`
	Price int              `json:"price"`
	Address string         `json:"address"`
	Room_count int         `json:"room_count"`
	Acreage int            `json:"acreage"`
	Unit string            `json:"unit"`
	Capacity int           `json:"capacity"`
	Beds string            `json:"beds"`
	Deposit int            `json:"deposit"`
	Min_days int           `json:"min_days" gorm:"default:1"`
	Max_days int           `json:"max_dayst" gorm:"default:0"`
	Order_count int        `json:"order_count" gorm:"default:0"`
	Index_image_url string `json:"index_image_url" gorm:"default:''"`
	Facilities []Facility  `gorm:"many2many:houses_facility;"`
	Images []HouseImage
	Orders []Orders
}





type HouseImage struct {
	gorm.Model
	HousesID uint         `json:"houses_id"`
	Url string            `json:"url"`
}

type Facility struct {
	gorm.Model
	Name string           `json:"name" gorm:"not null"`
	Houses []Houses       `gorm:"many2many:houses_facility;"`
}

type HousesFacility struct {
	gorm.Model
	HousesID uint
	FacilityID uint
}

func AddHouse(data map[string]interface{}) bool {
	house := Houses{
		UserID: data["user_id"].(uint),
		AreaID: data["area_id"].(uint),
		Title: data["title"].(string),
		Price: data["price"].(int),
		Address :data["address"].(string),
	}
	db.Create(&house)
	return true
}

func GetUserHouse(userid uint) (house []Houses) {
	db.Where("user_id = ?",userid).Find(&house)
	return
}

func GetHouse(id uint) (house Houses) {
	db.Where("id = ?",id).First(&house)
	return
}



func IndexHouse() (houses []Houses) {
	db.Limit(10).Order("created_at desc").Find(&houses)
	return
}

func UpdateHouseOderById(id uint) bool {
	var house Houses
	house.Order_count ++
	db.Model(&house).Where("id = ?", id).Update("order_count",house.Order_count)
	return true
}



