package formatter

import "cashier/model"

type ProductFormatter struct {
	ProductId    string
	ProductName  string
	Price        float64
	Thumbnail    string
	Status       string
	CategoryId   string
	CategoryName string
}

func FormatterProduct(product []*model.ProductModel) []*ProductFormatter {
	var productFormatter []*ProductFormatter
	for _, p := range product {
		categoryName := ""
		// Memeriksa apakah produk memiliki kategori
		if len(p.Categories) > 0 {
			// Mengambil nama kategori dari kategori pertama (asumsikan satu produk hanya memiliki satu kategori)
			categoryName = p.Categories[0].CategoryName
		}
		productFormatter = append(productFormatter, &ProductFormatter{
			ProductId:    p.Id,
			ProductName:  p.ProductName,
			Price:        p.Price,
			Thumbnail:    p.Thumbnail,
			Status:       string(p.Status),
			CategoryId:   p.CategoryId,
			CategoryName: categoryName,
		})
	}

	return productFormatter
}
