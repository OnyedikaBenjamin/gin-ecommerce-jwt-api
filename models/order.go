package models

import "gorm.io/gorm"

type Order struct {
    gorm.Model
    UserID uint
    Status string `gorm:"default:pending"`
    Items  []OrderItem
}

type OrderItem struct {
    gorm.Model
    OrderID   uint
    ProductID uint
    Quantity  int
}
