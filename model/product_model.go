package model

import "time"

type ProductModel struct {
	Id          string
	ProductName string
	Thumbnail   string
	Price       float64
	Status      StatusEnum
	CategoryId  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Categories  []*CategoryModel
}

type InputProduct struct {
	ProductName string  `json:"productName" validate:"required,max=100"`
	Thumbnail   string  `json:"thumbnail" validate:"required,max=100"`
	Price       float64 `json:"price" validate:"required"`
	CategoryId  string  `json:"categoryId" validate:"required"`
}

type UpdateProduct struct {
	ProductName string `json:"productName" validate:"max=100"`
	Thumbnail   string
	Price       float64    `json:"price"`
	Status      StatusEnum `json:"status"`
	CategoryId  string     `json:"categoryId"`
	UpdatedAt   time.Time
}

type ProductImage struct {
	Id          string
	ProductId   string
	ProductFile string `json:"productName"`
}
