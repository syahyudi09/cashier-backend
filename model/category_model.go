package model

import "time"

type CategoryModel struct {
	Id           string
	CategoryName string
	Status       StatusEnum
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CategoryInput struct {
	CategoryName string `json:"category_name" validate:"required,max=100"`
}

type CategoryUpdate struct {
	CategoryName string     `json:"category_name" validate:"max=100"`
	Status       StatusEnum `json:"status"`
}
