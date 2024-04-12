package formatter

import "cashier/model"

type CategoryFormatter struct {
	CategoryId   string
	CategoryName string
	Status       string
}

func FormatterCategory(category []*model.CategoryModel) []*CategoryFormatter {
	var categoryFormatter []*CategoryFormatter
	for _, c := range category {
		categoryFormatter = append(categoryFormatter, &CategoryFormatter{
			CategoryId:   c.Id,
			CategoryName: c.CategoryName,
			Status:       string(c.Status),
		})
	}

	return categoryFormatter
}
