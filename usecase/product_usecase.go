package usecase

import (
	"cashier/model"
	"cashier/repository"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ProductUsecase interface {
	Create(input model.InputProduct) error                                                          // 1
	ProductNameExits(name string) (bool, error)                                                      // 2
	FindById(id string) (*model.ProductModel, error)                                                 // 3
	Update(id string, updatedProduct model.UpdateProduct) error                                      //4
	GetAll(c context.Context, page, pageSize int, search string) ([]*model.ProductModel, int, error) // 5
	CheckNameForUpdate(name, id string) (bool, error)      
	Delete(ID string) error                                          // 6
	UploadProduct(productImage model.ProductImage) error
}

type productUsecase struct {
	productRepository repository.ProductRepository
}

func (pu *productUsecase) Create(input model.InputProduct) error {
	p := model.ProductModel{}
	p.ProductName = input.ProductName
	p.Thumbnail = input.Thumbnail
	p.Price = input.Price
	p.Status = "Active"
	p.CategoryId = input.CategoryId
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	
	err := pu.productRepository.Create(p)
	if err != nil {
		return fmt.Errorf("failed to create product: %v", err)
	}

	return nil
}

func (pu *productUsecase) GetAll(c context.Context, page, pageSize int, search string) ([]*model.ProductModel, int, error) {
	product, totalDocs, err := pu.productRepository.GetAll(c, page, pageSize, search)
	if err != nil {
		return nil, 0, fmt.Errorf("error on  *productUsecase.GetAll: %v", err)
	}

	return product, totalDocs, nil
}

func (pu *productUsecase) ProductNameExits(name string) (bool, error) {
	return pu.productRepository.ProductCheckName(name)
}

func (pu *productUsecase) Update(id string, updatedProduct model.UpdateProduct) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidUUIDFormat
	}

	product, err := pu.productRepository.FindById(id)
	if err != nil {
		return fmt.Errorf("error fetching product: %w", err)
	}

	if updatedProduct.ProductName != "" {
		product.ProductName = updatedProduct.ProductName
	}

	if updatedProduct.Price != 0 {
		product.Price = updatedProduct.Price
	}

	if updatedProduct.Status != ""{
		product.Status = updatedProduct.Status
	}

	if updatedProduct.CategoryId != ""{
		product.CategoryId = updatedProduct.CategoryId
	}

	if updatedProduct.Thumbnail != "" {
		product.ProductName = updatedProduct.Thumbnail
	}

	err = pu.productRepository.Update(id, *product)
	if err != nil {
		return fmt.Errorf("error fetching product: %w", err)
	}
	return nil
}

func (pu *productUsecase) FindById(id string) (*model.ProductModel, error) {
	return pu.productRepository.FindById(id)
}

func (pu *productUsecase) CheckNameForUpdate(name, id string) (bool, error) {
	exists, err := pu.productRepository.CheckNameForUpdate(name, id)
	if err != nil {
		return false, fmt.Errorf("error checking name for update: %w", err)
	}
	return exists, nil
}

func (pu *productUsecase) Delete(userID string) error {
	if _, err := uuid.Parse(userID); err != nil {
		return ErrInvalidUUIDFormat
	}
	return pu.productRepository.Delete(userID)
}

func(pu *productUsecase) UploadProduct(productImage model.ProductImage) error {
	return pu.productRepository.UploadProduct(productImage)
}


func NewProductUsecase(pr repository.ProductRepository) ProductUsecase {
	return &productUsecase{
		productRepository: pr,
	}
}
