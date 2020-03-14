package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Orders struct {
	gorm.Model
	UserID uint              `json:"user_id"`
	HousesID uint            `json:"houses_id"`
	Begin_date time.Time    `json:"begin_date" gorm:"not null"`
	End_date time.Time       `json:"end_date" gorm:"not null"`
	Days int                 `json:"days" gorm:"not null"`
	House_price int          `json:"house_price" gorm:"not null"`           //单价
	Amount int               `json:"amount" gorm:"not null"`               //总价
	Status string            `json:"status" gorm:"index;default:'WAIT_ACCEPT'"`
	Comment string           `json:"comment"`
	Trade_no string          `json:"trade_no"`
}

func AddOrders(data map[string]interface{}) bool {
	order := Orders{
		UserID:      data["user_id"].(uint),
		HousesID:    data["house_id"].(uint),
		Begin_date:  data["start_date"].(time.Time),
		End_date:    data["end_date"].(time.Time) ,
		Days:        data["days"].(int),
		House_price: data["house_price"].(int),
		Amount:      data["amount"].(int),
	}
	db.Create(&order)
	return true
}

func GetConflictOrder (data map[string]interface{}) (c_order Orders) {
	db.Where("houses_id = ? AND begin_date < ? AND end_date > ? AND status in (?)", data["house_id"],data["end_date"],data["start_date"],[]string{"WAIT_ACCEPT", "WAIT_PAYMENT", "PAID"}).First(&c_order)
	return
}

func GetOrders(houses_id []uint) (orders []Orders) {
	db.Where("houses_id in (?)",houses_id).Find(&orders)
	return
}

func GetUserOrders (user_id uint) (orders []Orders) {
	db.Where("user_id = ?", user_id).Find(&orders)
	return
}

func GetOrder(order_id uint,status string) (order Orders) {
	db.Where("id = ? AND status = ?", order_id,status).First(&order)
	return
}

func UpdateOrderStatus(status,reason string, order_id uint) bool {
	var order Orders
	db.Model(&order).Where("id = ?", order_id).Update(map[string]interface{}{"status": status, "comment":reason})
	return true
}

func GetOrderC(order_id uint,status string,user_id uint) (order Orders) {
	db.Where("id = ? AND status = ? AND user_id = ?", order_id,status,user_id).First(&order)
	return
}

