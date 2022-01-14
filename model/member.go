package model

import (
	"time"
)

type Member struct {
	ID                int `gorm:"primaryKey"`
	Email             string
	Password          string
	AccessToken       string
	RefreshToken      string
	EmailAuth         bool
	ThirdAuthPassword string
	ActivatedAt       time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
