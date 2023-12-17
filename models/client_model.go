// Package models defines object-relational mapping for fluxy database
package models

import "gorm.io/gorm"

// Client is the sql model for fluxy rate-limiter clients
type Client struct{
  gorm.Model

  Name     string  `gorm:"unique;not null"`
  Limit    uint64  `gorm:"not null"`
  Duration uint64  `gorm:"not null"`
}
