package models

import "time"

type Account struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CompanyID uint      `json:"company_id"`
	Company   Company   `json:"company" gorm:"foreignKey:CompanyID"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	IsAdmin   bool      `gorm:"default:false" json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
