package model

import "github.com/jinzhu/gorm"

type Area struct {
	gorm.Model
	Name string           `json:"name" gorm:"not null"`
	Houses []Houses       `gorm:"foreignkey:AreaID"`
}

type Areas struct {
	ID uint        `json:"id"`
	Name string    `json:"name"`
}

func GetArea(id uint) (area Area, err error) {
	db.Where("id = ?", id).First(&area)
	return area,nil
}

func GetAllArea() (areas []Areas) {
	db.Table("area").Select("id,name").Scan(&areas)
	return
}

func AddArea(data map[string]interface{}) bool {
	area := Area{
		Name: data["name"].(string),
	}
	db.Create(&area)
	return true
}

func GetAllAreaHouses(area Area) (houses []Houses) {
	houses = make([]Houses,512)
	db.Model(&area).Related(&houses)
	return houses
}